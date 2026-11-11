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

type RetryConfig struct {
	MaxRetries     int           // default: 3
	InitialBackoff time.Duration // default: 100ms
	MaxBackoff     time.Duration // default: 5s
	BackoffFactor  float64       // default: 2.0 (exponential)
}

type CircuitBreakerConfig struct {
	MaxRequests uint32        // max requests allowed in half-open state (default: 3)
	Interval    time.Duration // interval to clear internal counters (default: 60s)
	Timeout     time.Duration // time in open state before attempting recovery (default: 30s)
	ReadyToTrip func(counts gobreaker.Counts) bool
}

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
	// SkipAuth skips the OAuth client-credentials flow entirely (no
	// Authorization header set), for calls that never cross the OIDC
	// trust boundary. See NewInternalHTTPClient.
	SkipAuth bool
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     5 * time.Second,
		BackoffFactor:  2.0,
	}
}

func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Open circuit once >=10 requests seen and >=50% failed.
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 10 && failureRatio >= 0.5
		},
	}
}

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

	if retryConfig == nil {
		defaultRetry := DefaultRetryConfig()
		retryConfig = &defaultRetry
	}
	if cbConfig == nil {
		defaultCB := DefaultCircuitBreakerConfig()
		cbConfig = &defaultCB
	}

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

// NewInternalHTTPClient creates an EnhancedHTTPClient (retry + circuit
// breaker) with SkipAuth set, for calls to a network-internal endpoint
// that isn't behind the OIDC trust boundary (e.g. one protected by a
// NetworkPolicy instead). It never attempts the OAuth client-credentials
// flow.
func NewInternalHTTPClient(
	timeout time.Duration,
	logger LoggerInterface,
	retryConfig *RetryConfig,
	cbConfig *CircuitBreakerConfig,
) *EnhancedHTTPClient {
	client := NewEnhancedHTTPClient(timeout, "", "", "", "", logger, retryConfig, cbConfig)
	client.SkipAuth = true
	return client
}

func (h *EnhancedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if !h.SkipAuth {
		if err := h.ensureAuthToken(); err != nil {
			return nil, fmt.Errorf("failed to ensure auth token: %w", err)
		}
	}

	result, err := h.CircuitBreaker.Execute(func() (interface{}, error) {
		return h.doWithRetry(req)
	})

	if err != nil {
		return nil, err
	}

	return result.(*http.Response), nil
}

func (h *EnhancedHTTPClient) doWithRetry(req *http.Request) (*http.Response, error) {
	var lastErr error
	backoff := h.RetryConfig.InitialBackoff

	for attempt := 0; attempt <= h.RetryConfig.MaxRetries; attempt++ {
		// Request body may already be consumed, so clone before each attempt.
		clonedReq, err := h.cloneRequest(req)
		if err != nil {
			return nil, fmt.Errorf("failed to clone request: %w", err)
		}

		if !h.SkipAuth {
			h.TokenMutex.RLock()
			clonedReq.Header.Set("Authorization", "Bearer "+h.AuthToken)
			h.TokenMutex.RUnlock()
		}

		start := time.Now()
		resp, err := h.Client.Do(clonedReq)
		elapsed := time.Since(start)

		statusCode := 0
		if resp != nil {
			statusCode = resp.StatusCode
		}
		h.Logger.Debugf("HTTP request: method=%s, url=%s, attempt=%d/%d, duration=%s, status=%d, error=%v",
			clonedReq.Method, clonedReq.URL.String(), attempt+1, h.RetryConfig.MaxRetries+1, elapsed, statusCode, err)

		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		if err != nil {
			lastErr = err
		} else {
			lastErr = fmt.Errorf("server returned %d status code", resp.StatusCode)
			resp.Body.Close()
		}

		if attempt == h.RetryConfig.MaxRetries {
			break
		}

		h.Logger.Warnf("Request failed (attempt %d/%d): %v. Retrying in %s...",
			attempt+1, h.RetryConfig.MaxRetries+1, lastErr, backoff)

		time.Sleep(backoff)

		backoff = time.Duration(float64(backoff) * h.RetryConfig.BackoffFactor)
		if backoff > h.RetryConfig.MaxBackoff {
			backoff = h.RetryConfig.MaxBackoff
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", h.RetryConfig.MaxRetries+1, lastErr)
}

func (h *EnhancedHTTPClient) ensureAuthToken() error {
	h.TokenMutex.RLock()
	hasToken := h.AuthToken != ""
	h.TokenMutex.RUnlock()

	if hasToken {
		return nil
	}

	h.TokenMutex.Lock()
	defer h.TokenMutex.Unlock()

	// Re-check: another goroutine may have refreshed while we waited for the lock.
	if h.AuthToken != "" {
		return nil
	}

	return h.refreshAuthToken()
}

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

func (h *EnhancedHTTPClient) cloneRequest(req *http.Request) (*http.Request, error) {
	cloned := req.Clone(context.Background())

	// Body is a stream that gets consumed; read it once and give both the
	// original and the clone their own copy.
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		cloned.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return cloned, nil
}

func (h *EnhancedHTTPClient) GetCircuitBreakerState() gobreaker.State {
	return h.CircuitBreaker.State()
}

func (h *EnhancedHTTPClient) GetCircuitBreakerCounts() gobreaker.Counts {
	return h.CircuitBreaker.Counts()
}

func calculateBackoff(attempt int, config RetryConfig) time.Duration {
	backoff := float64(config.InitialBackoff) * math.Pow(config.BackoffFactor, float64(attempt))
	if backoff > float64(config.MaxBackoff) {
		backoff = float64(config.MaxBackoff)
	}
	return time.Duration(backoff)
}
