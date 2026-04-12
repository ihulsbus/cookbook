package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/sony/gobreaker"
)

// RetryConfig defines retry behavior
type RetryConfig struct {
	MaxRetries     int           // Maximum number of retry attempts (default: 3)
	InitialBackoff time.Duration // Initial backoff duration (default: 100ms)
	MaxBackoff     time.Duration // Maximum backoff duration (default: 5s)
	BackoffFactor  float64       // Backoff multiplier (default: 2.0 for exponential)
}

// CircuitBreakerConfig defines circuit breaker behavior
type CircuitBreakerConfig struct {
	MaxRequests uint32        // Max requests allowed in half-open state (default: 3)
	Interval    time.Duration // Interval to clear internal counters (default: 60s)
	Timeout     time.Duration // Time in open state before attempting recovery (default: 30s)
	// ReadyToTrip is called when a request fails. If it returns true, circuit opens
	ReadyToTrip func(counts gobreaker.Counts) bool
}

// EnhancedHTTPClient provides HTTP client with circuit breaker and retry logic
type EnhancedHTTPClient struct {
	Client         *http.Client
	AuthToken      string
	TokenMutex     sync.RWMutex
	KeycloakURL    string
	KeycloakRealm  string
	ClientID       string
	ClientSecret   string
	Logger         LoggerInterface
	CircuitBreaker *gobreaker.CircuitBreaker
	RetryConfig    RetryConfig
	requestCounter uint64
	requestMutex   sync.Mutex
}

// DefaultRetryConfig returns sensible retry defaults
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     5 * time.Second,
		BackoffFactor:  2.0,
	}
}

// DefaultCircuitBreakerConfig returns sensible circuit breaker defaults
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Open circuit if:
			// - At least 10 requests have been made
			// - More than 50% of recent requests failed
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 10 && failureRatio >= 0.5
		},
	}
}

// NewEnhancedHTTPClient creates an HTTP client with circuit breaker and retry logic
func NewEnhancedHTTPClient(
	timeout time.Duration,
	keycloakURL, keycloakRealm, clientID, clientSecret string,
	logger LoggerInterface,
	retryConfig *RetryConfig,
	cbConfig *CircuitBreakerConfig,
) *EnhancedHTTPClient {
	if timeout < time.Second || timeout > 30*time.Second {
		panic("timeout must be between 1 and 30 seconds")
	}

	// Use defaults if not provided
	if retryConfig == nil {
		defaultRetry := DefaultRetryConfig()
		retryConfig = &defaultRetry
	}
	if cbConfig == nil {
		defaultCB := DefaultCircuitBreakerConfig()
		cbConfig = &defaultCB
	}

	// Create circuit breaker
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "HTTPClient",
		MaxRequests: cbConfig.MaxRequests,
		Interval:    cbConfig.Interval,
		Timeout:     cbConfig.Timeout,
		ReadyToTrip: cbConfig.ReadyToTrip,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Warnf("Circuit breaker state changed: %s -> %s", from.String(), to.String())
		},
	})

	return &EnhancedHTTPClient{
		Client: &http.Client{
			Timeout: timeout,
		},
		KeycloakURL:    keycloakURL,
		KeycloakRealm:  keycloakRealm,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		Logger:         logger,
		CircuitBreaker: cb,
		RetryConfig:    *retryConfig,
	}
}

// Do sends an HTTP request with circuit breaker and retry logic
func (h *EnhancedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Ensure auth token is present
	if err := h.ensureAuthToken(); err != nil {
		return nil, fmt.Errorf("failed to ensure auth token: %w", err)
	}

	// Execute with circuit breaker
	result, err := h.CircuitBreaker.Execute(func() (interface{}, error) {
		return h.doWithRetry(req)
	})

	if err != nil {
		return nil, err
	}

	return result.(*http.Response), nil
}

// doWithRetry performs the request with exponential backoff retry
func (h *EnhancedHTTPClient) doWithRetry(req *http.Request) (*http.Response, error) {
	var lastErr error
	backoff := h.RetryConfig.InitialBackoff

	for attempt := 0; attempt <= h.RetryConfig.MaxRetries; attempt++ {
		// Clone the request for retry (body might be consumed)
		clonedReq, err := h.cloneRequest(req)
		if err != nil {
			return nil, fmt.Errorf("failed to clone request: %w", err)
		}

		// Set authorization header
		h.TokenMutex.RLock()
		clonedReq.Header.Set("Authorization", "Bearer "+h.AuthToken)
		h.TokenMutex.RUnlock()

		// Execute request
		start := time.Now()
		resp, err := h.Client.Do(clonedReq)
		elapsed := time.Since(start)

		// Log request
		statusCode := 0
		if resp != nil {
			statusCode = resp.StatusCode
		}
		h.Logger.Debugf("HTTP request: method=%s, url=%s, attempt=%d/%d, duration=%s, status=%d, error=%v",
			clonedReq.Method, clonedReq.URL.String(), attempt+1, h.RetryConfig.MaxRetries+1, elapsed, statusCode, err)

		// Success cases
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		// Store error for potential return
		if err != nil {
			lastErr = err
		} else {
			// Server error (5xx)
			lastErr = fmt.Errorf("server returned %d status code", resp.StatusCode)
			resp.Body.Close()
		}

		// Don't retry on last attempt
		if attempt == h.RetryConfig.MaxRetries {
			break
		}

		// Exponential backoff
		h.Logger.Warnf("Request failed (attempt %d/%d): %v. Retrying in %s...",
			attempt+1, h.RetryConfig.MaxRetries+1, lastErr, backoff)

		time.Sleep(backoff)

		// Calculate next backoff with jitter
		backoff = time.Duration(float64(backoff) * h.RetryConfig.BackoffFactor)
		if backoff > h.RetryConfig.MaxBackoff {
			backoff = h.RetryConfig.MaxBackoff
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", h.RetryConfig.MaxRetries+1, lastErr)
}

// ensureAuthToken ensures we have a valid auth token
func (h *EnhancedHTTPClient) ensureAuthToken() error {
	h.TokenMutex.RLock()
	hasToken := h.AuthToken != ""
	h.TokenMutex.RUnlock()

	if hasToken {
		return nil
	}

	h.TokenMutex.Lock()
	defer h.TokenMutex.Unlock()

	// Double-check after acquiring write lock
	if h.AuthToken != "" {
		return nil
	}

	return h.refreshAuthToken()
}

// refreshAuthToken obtains a new authentication token
func (h *EnhancedHTTPClient) refreshAuthToken() error {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", h.KeycloakURL, h.KeycloakRealm)
	data := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", h.ClientID, h.ClientSecret)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := h.Client.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected token response status: %d, body: %s", resp.StatusCode, body)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return fmt.Errorf("access_token not found in response")
	}

	h.AuthToken = token
	h.Logger.Debugf("Successfully refreshed authentication token")
	return nil
}

// cloneRequest creates a copy of the request for retry purposes
func (h *EnhancedHTTPClient) cloneRequest(req *http.Request) (*http.Request, error) {
	// Clone the request
	cloned := req.Clone(context.Background())

	// If there's a body, we need to restore it
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		// Restore original body
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Set cloned body
		cloned.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return cloned, nil
}

// GetCircuitBreakerState returns the current state of the circuit breaker
func (h *EnhancedHTTPClient) GetCircuitBreakerState() gobreaker.State {
	return h.CircuitBreaker.State()
}

// GetCircuitBreakerCounts returns current circuit breaker statistics
func (h *EnhancedHTTPClient) GetCircuitBreakerCounts() gobreaker.Counts {
	return h.CircuitBreaker.Counts()
}

// calculateBackoff calculates exponential backoff with jitter
func calculateBackoff(attempt int, config RetryConfig) time.Duration {
	backoff := float64(config.InitialBackoff) * math.Pow(config.BackoffFactor, float64(attempt))
	if backoff > float64(config.MaxBackoff) {
		backoff = float64(config.MaxBackoff)
	}
	return time.Duration(backoff)
}
