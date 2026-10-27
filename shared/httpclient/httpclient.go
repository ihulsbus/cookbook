package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type LoggerInterface interface {
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

// HTTPClient provides a shared HTTP client with default settings.
type HTTPClient struct {
	Client        *http.Client
	AuthToken     string
	TokenMutex    sync.RWMutex
	KeycloakURL   string
	KeycloakRealm string
	ClientID      string
	ClientSecret  string
	Logger        LoggerInterface
}

// NewHTTPClient initializes and returns an HTTPClient instance.
func NewHTTPClient(timeout time.Duration, keycloakURL, keycloakRealm, clientID, clientSecret string, logger LoggerInterface) *HTTPClient {
	if timeout < time.Second || timeout > 10*time.Second {
		panic("timeout must be between 1 and 10 seconds")
	}

	return &HTTPClient{
		Client: &http.Client{
			Timeout: timeout,
		},
		KeycloakURL:   keycloakURL,
		KeycloakRealm: keycloakRealm,
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		Logger:        logger,
	}
}

// Do sends an HTTP request and returns the HTTP response.
func (h *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	h.TokenMutex.RLock()
	token := h.AuthToken
	h.TokenMutex.RUnlock()

	if token == "" {
		h.TokenMutex.Lock()
		if h.AuthToken == "" {
			if err := h.refreshAuthToken(); err != nil {
				h.TokenMutex.Unlock()
				return nil, fmt.Errorf("failed to refresh auth token: %w", err)
			}
		}
		token = h.AuthToken
		h.TokenMutex.Unlock()
	}

	req.Header.Set("Authorization", "Bearer "+h.AuthToken)
	start := time.Now()
	resp, err := h.Client.Do(req)
	elapsed := time.Since(start)

	h.Logger.Debugf("Http request: method=%s, url=%s, duration=%s, status=%d, error=%v", req.Method, req.URL.String(), elapsed, resp.StatusCode, err)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *HTTPClient) refreshAuthToken() error {
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
	return nil
}
