package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sony/gobreaker"
)

type MockLogger struct{}

func (m *MockLogger) Debugf(format string, args ...interface{}) {}
func (m *MockLogger) Warnf(format string, args ...interface{})  {}

func TestEnhancedHTTPClient_SuccessfulRequest(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	client := NewEnhancedHTTPClient(
		5*time.Second,
		"http://keycloak.example.com",
		"test-realm",
		"test-client",
		"test-secret",
		&MockLogger{},
		nil, // Use defaults
		nil, // Use defaults
	)

	// Set token manually to avoid Keycloak call
	client.AuthToken = "test-token"

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", resp.StatusCode)
	}
}

func TestEnhancedHTTPClient_RetryOnFailure(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success after retries"))
	}))
	defer server.Close()

	retryConfig := RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	client := NewEnhancedHTTPClient(
		5*time.Second,
		"http://keycloak.example.com",
		"test-realm",
		"test-client",
		"test-secret",
		&MockLogger{},
		&retryConfig,
		nil,
	)
	client.AuthToken = "test-token"

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("Expected no error after retries, got: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", resp.StatusCode)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got: %d", attempts)
	}
}

func TestEnhancedHTTPClient_MaxRetriesExceeded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	retryConfig := RetryConfig{
		MaxRetries:     2,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	client := NewEnhancedHTTPClient(
		5*time.Second,
		"http://keycloak.example.com",
		"test-realm",
		"test-client",
		"test-secret",
		&MockLogger{},
		&retryConfig,
		nil,
	)
	client.AuthToken = "test-token"

	req, _ := http.NewRequest("GET", server.URL, nil)
	_, err := client.Do(req)

	if err == nil {
		t.Fatal("Expected error after max retries, got nil")
	}
}

func TestEnhancedHTTPClient_CircuitBreakerOpens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	retryConfig := RetryConfig{
		MaxRetries:     0, // No retries to make test faster
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	cbConfig := CircuitBreakerConfig{
		MaxRequests: 1,
		Interval:    1 * time.Second,
		Timeout:     1 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Open after 5 failures
			return counts.ConsecutiveFailures >= 5
		},
	}

	client := NewEnhancedHTTPClient(
		5*time.Second,
		"http://keycloak.example.com",
		"test-realm",
		"test-client",
		"test-secret",
		&MockLogger{},
		&retryConfig,
		&cbConfig,
	)
	client.AuthToken = "test-token"

	// Make requests until circuit opens
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequest("GET", server.URL, nil)
		client.Do(req)
	}

	// Circuit should be open now
	state := client.GetCircuitBreakerState()
	if state != gobreaker.StateOpen {
		t.Errorf("Expected circuit breaker to be open, got: %s", state.String())
	}

	// Next request should fail immediately with circuit breaker error
	req, _ := http.NewRequest("GET", server.URL, nil)
	_, err := client.Do(req)

	if err == nil {
		t.Fatal("Expected circuit breaker error, got nil")
	}

	if err != gobreaker.ErrOpenState {
		t.Errorf("Expected ErrOpenState, got: %v", err)
	}
}

func TestEnhancedHTTPClient_CircuitBreakerRecovery(t *testing.T) {
	failureCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failureCount++
		if failureCount <= 5 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		// Start succeeding after 5 failures
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	retryConfig := RetryConfig{
		MaxRetries:     0,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	cbConfig := CircuitBreakerConfig{
		MaxRequests: 2,
		Interval:    500 * time.Millisecond,
		Timeout:     200 * time.Millisecond, // Short timeout for testing
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
	}

	client := NewEnhancedHTTPClient(
		5*time.Second,
		"http://keycloak.example.com",
		"test-realm",
		"test-client",
		"test-secret",
		&MockLogger{},
		&retryConfig,
		&cbConfig,
	)
	client.AuthToken = "test-token"

	// Trigger circuit opening
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", server.URL, nil)
		client.Do(req)
	}

	// Verify circuit is open
	if client.GetCircuitBreakerState() != gobreaker.StateOpen {
		t.Fatal("Circuit should be open after failures")
	}

	// Wait for circuit to enter half-open state
	time.Sleep(300 * time.Millisecond)

	// Make successful request to close circuit
	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)

	if err != nil {
		t.Logf("Request error (might be expected in half-open): %v", err)
	}

	// Wait a bit and make another successful request
	time.Sleep(100 * time.Millisecond)
	req, _ = http.NewRequest("GET", server.URL, nil)
	resp, err = client.Do(req)

	if err != nil {
		t.Fatalf("Expected success after recovery, got: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 after recovery, got: %d", resp.StatusCode)
	}
}

func TestEnhancedHTTPClient_NoRetryOn4xx(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest) // 4xx should not retry
	}))
	defer server.Close()

	retryConfig := RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	client := NewEnhancedHTTPClient(
		5*time.Second,
		"http://keycloak.example.com",
		"test-realm",
		"test-client",
		"test-secret",
		&MockLogger{},
		&retryConfig,
		nil,
	)
	client.AuthToken = "test-token"

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("Expected no error on 4xx, got: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got: %d", resp.StatusCode)
	}

	// Should only attempt once (no retries for 4xx)
	if attempts != 1 {
		t.Errorf("Expected 1 attempt for 4xx error, got: %d", attempts)
	}
}

func TestDefaultConfigs(t *testing.T) {
	retryConfig := DefaultRetryConfig()
	if retryConfig.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries=3, got: %d", retryConfig.MaxRetries)
	}

	cbConfig := DefaultCircuitBreakerConfig()
	if cbConfig.MaxRequests != 3 {
		t.Errorf("Expected MaxRequests=3, got: %d", cbConfig.MaxRequests)
	}
}
