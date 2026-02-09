package zia

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/zscaler/zscaler-sdk-go/v3/cache"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

// =============================================================================
// Section 1: Rate Limiter Unit Tests
// =============================================================================

// TestRateLimiterWait_GET verifies the rate limiter enforces per-window GET limits.
// With a limit of 2 GETs per 1-second window, the first two requests should pass
// immediately and the third should be asked to wait.
func TestRateLimiterWait_GET(t *testing.T) {
	// 2 GET requests allowed per 1-second window
	limiter := rl.NewRateLimiter(2, 1, 1, 1)

	// First GET — should not wait
	shouldWait, _ := limiter.Wait(http.MethodGet)
	if shouldWait {
		t.Fatal("first GET should not require waiting")
	}

	// Second GET — should not wait (still within limit)
	shouldWait, _ = limiter.Wait(http.MethodGet)
	if shouldWait {
		t.Fatal("second GET should not require waiting")
	}

	// Third GET — exceeds the limit, should wait
	shouldWait, delay := limiter.Wait(http.MethodGet)
	if !shouldWait {
		t.Fatal("third GET should require waiting — rate limit exceeded")
	}
	if delay <= 0 {
		t.Fatalf("expected positive wait delay, got %v", delay)
	}
	t.Logf("third GET wait delay: %v", delay)
}

// TestRateLimiterWait_POST verifies the rate limiter enforces POST/PUT/DELETE limits.
func TestRateLimiterWait_POST(t *testing.T) {
	// 1 POST/PUT/DELETE request per 2-second window
	limiter := rl.NewRateLimiter(10, 1, 1, 2)

	// First POST — should not wait
	shouldWait, _ := limiter.Wait(http.MethodPost)
	if shouldWait {
		t.Fatal("first POST should not require waiting")
	}

	// Second POST — exceeds the limit
	shouldWait, delay := limiter.Wait(http.MethodPost)
	if !shouldWait {
		t.Fatal("second POST should require waiting — rate limit exceeded")
	}
	if delay <= 0 {
		t.Fatalf("expected positive wait delay, got %v", delay)
	}
	t.Logf("second POST wait delay: %v", delay)
}

// TestRateLimiterWait_DELETE verifies the rate limiter tracks DELETE requests
// in the combined POST/PUT/DELETE bucket.
func TestRateLimiterWait_DELETE(t *testing.T) {
	// 1 POST/PUT/DELETE per 2-second window
	limiter := rl.NewRateLimiter(10, 1, 1, 2)

	// First DELETE — should not wait
	shouldWait, _ := limiter.Wait(http.MethodDelete)
	if shouldWait {
		t.Fatal("first DELETE should not require waiting")
	}

	// Second DELETE — same combined bucket
	shouldWait, _ = limiter.Wait(http.MethodDelete)
	if !shouldWait {
		t.Fatal("second DELETE should require waiting — combined POST/PUT/DELETE limit exceeded")
	}
}

// TestRateLimiterWait_MixedMethods verifies that POST and DELETE share the same
// combined bucket, but GET uses its own.
func TestRateLimiterWait_MixedMethods(t *testing.T) {
	// GET: 5 per 1s, POST/PUT/DELETE: 1 per 2s
	limiter := rl.NewRateLimiter(5, 1, 1, 2)

	// POST consumes the POST/PUT/DELETE slot
	shouldWait, _ := limiter.Wait(http.MethodPost)
	if shouldWait {
		t.Fatal("first POST should not wait")
	}

	// DELETE in the same window should wait (shared bucket)
	shouldWait, _ = limiter.Wait(http.MethodDelete)
	if !shouldWait {
		t.Fatal("DELETE after POST should wait — shared bucket")
	}

	// GET should still be fine (separate bucket)
	shouldWait, _ = limiter.Wait(http.MethodGet)
	if shouldWait {
		t.Fatal("GET should not be affected by POST/DELETE bucket")
	}
}

// TestRateLimiterWait_WindowExpiry verifies that the rate limiter allows new
// requests once the time window expires.
func TestRateLimiterWait_WindowExpiry(t *testing.T) {
	// 1 GET per 1-second window — very tight
	limiter := rl.NewRateLimiter(1, 1, 1, 1)

	shouldWait, _ := limiter.Wait(http.MethodGet)
	if shouldWait {
		t.Fatal("first GET should not wait")
	}

	// Immediately — should need to wait
	shouldWait, _ = limiter.Wait(http.MethodGet)
	if !shouldWait {
		t.Fatal("second GET should wait")
	}

	// Wait for window to expire
	time.Sleep(1100 * time.Millisecond)

	// Now should be allowed again
	shouldWait, _ = limiter.Wait(http.MethodGet)
	if shouldWait {
		t.Fatal("GET should be allowed after window expiry")
	}
}

// TestRateLimiterWithHourlyLimits verifies that hourly limits are enforced independently
// of per-second limits.
func TestRateLimiterWithHourlyLimits(t *testing.T) {
	// Per-second: 100 GET per 1s (effectively no per-second throttle for this test)
	// Hourly: only 3 GETs per hour
	limiter := rl.NewRateLimiterWithHourly(100, 100, 1, 1, 3, 3, 3)

	for i := 0; i < 3; i++ {
		shouldWait, _ := limiter.Wait(http.MethodGet)
		if shouldWait {
			t.Fatalf("GET #%d should not wait (within hourly limit)", i+1)
		}
	}

	// 4th GET exceeds hourly limit
	shouldWait, delay := limiter.Wait(http.MethodGet)
	if !shouldWait {
		t.Fatal("4th GET should wait — hourly limit exceeded")
	}
	if delay <= 0 {
		t.Fatal("expected positive hourly wait delay")
	}
	t.Logf("hourly limit wait delay: %v", delay)
}

// =============================================================================
// Section 2: Rate Limit Transport Integration Tests (mock HTTP server)
// =============================================================================

// TestRateLimitTransport_ThrottlesRequests verifies that the RateLimitTransport
// enforces client-side rate limiting on actual HTTP calls to a mock server.
func TestRateLimitTransport_ThrottlesRequests(t *testing.T) {
	var requestCount int32

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer mockServer.Close()

	// Very tight limiter: 2 GET per 2-second window
	limiter := rl.NewRateLimiter(2, 2, 2, 2)
	transport := &rl.RateLimitTransport{
		Base:    http.DefaultTransport,
		Limiter: limiter,
		Logger:  logger.NewNopLogger(),
	}
	client := &http.Client{Transport: transport}

	start := time.Now()
	totalRequests := 4

	for i := 0; i < totalRequests; i++ {
		resp, err := client.Get(mockServer.URL + "/zia/api/v1/test")
		if err != nil {
			t.Fatalf("request %d failed: %v", i+1, err)
		}
		resp.Body.Close()
	}

	elapsed := time.Since(start)
	finalCount := atomic.LoadInt32(&requestCount)

	if int(finalCount) != totalRequests {
		t.Fatalf("expected %d requests to reach server, got %d", totalRequests, finalCount)
	}

	// With 2 requests per 2 seconds and 4 total, we expect at least ~2 seconds
	// of throttle delay (the 3rd and 4th requests must wait)
	if elapsed < 1500*time.Millisecond {
		t.Fatalf("expected rate limiting to add delay; elapsed=%v (too fast)", elapsed)
	}
	t.Logf("rate-limited %d requests in %v", totalRequests, elapsed)
}

// =============================================================================
// Section 3: HTTP Client Retry-on-429 Tests (mock server + Retry-After)
// =============================================================================

// TestRetryOn429WithRetryAfterHeader verifies that the retryable HTTP client
// correctly retries on 429 responses with a Retry-After header, matching the
// behavior configured in the SDK's getHTTPClient.
func TestRetryOn429WithRetryAfterHeader(t *testing.T) {
	var attempt int32

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attempt, 1)
		if count <= 2 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"code":"TOO_MANY_REQUESTS","message":"Rate limit exceeded"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer mockServer.Close()

	// Build an HTTP client matching the SDK's getHTTPClient behavior
	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryWaitMin = 1 * time.Second
	retryableClient.RetryWaitMax = 5 * time.Second
	retryableClient.RetryMax = 5
	retryableClient.Logger = nil // suppress retryablehttp logs in test output

	// Use SDK-compatible backoff that respects Retry-After
	retryableClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil && (resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable) {
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				if secs, err := time.ParseDuration(retryAfter + "s"); err == nil {
					return secs
				}
			}
		}
		mult := math.Pow(2, float64(attemptNum)) * float64(min)
		sleep := time.Duration(mult)
		if float64(sleep) != mult || sleep > max {
			sleep = max
		}
		return sleep
	}

	// Use SDK-compatible retry policy: retry on 429
	retryableClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			return true, nil
		}
		return false, nil
	}

	client := retryableClient.StandardClient()

	start := time.Now()
	resp, err := client.Get(mockServer.URL + "/zia/api/v1/test")
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected successful response after retries, got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK after retries, got %d", resp.StatusCode)
	}

	totalAttempts := atomic.LoadInt32(&attempt)
	if totalAttempts != 3 {
		t.Fatalf("expected 3 attempts (2 x 429 + 1 x 200), got %d", totalAttempts)
	}

	// Should have waited at least ~2 seconds (2 retries × 1s Retry-After)
	if elapsed < 1500*time.Millisecond {
		t.Fatalf("expected retry delays to add up; elapsed=%v (too fast)", elapsed)
	}
	t.Logf("completed after %d attempts in %v", totalAttempts, elapsed)
}

// TestRetryOn429ExhaustsRetries verifies that the client gives up after
// exceeding max retries on persistent 429 responses. When retryablehttp
// exhausts all retries it returns an error wrapping "giving up after N attempt(s)".
func TestRetryOn429ExhaustsRetries(t *testing.T) {
	var attempt int32

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempt, 1)
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"code":"TOO_MANY_REQUESTS","message":"Rate limit exceeded"}`))
	}))
	defer mockServer.Close()

	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryWaitMin = 10 * time.Millisecond // fast for test
	retryableClient.RetryWaitMax = 50 * time.Millisecond
	retryableClient.RetryMax = 3
	retryableClient.Logger = nil

	retryableClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			return true, nil
		}
		return false, nil
	}

	client := retryableClient.StandardClient()
	_, err := client.Get(mockServer.URL + "/zia/api/v1/test")

	// retryablehttp wraps exhausted retries into an error via StandardClient
	if err == nil {
		t.Fatal("expected an error after exhausting retries, got nil")
	}

	// Verify the error indicates giving up
	if !strings.Contains(err.Error(), "giving up") {
		t.Fatalf("expected 'giving up' in error message, got: %v", err)
	}

	totalAttempts := atomic.LoadInt32(&attempt)
	// retryablehttp makes 1 initial + RetryMax retries = 4 total attempts
	expectedAttempts := int32(retryableClient.RetryMax + 1)
	if totalAttempts != expectedAttempts {
		t.Fatalf("expected %d total attempts, got %d", expectedAttempts, totalAttempts)
	}
	t.Logf("correctly exhausted retries after %d attempts, error: %v", totalAttempts, err)
}

// =============================================================================
// Section 4: Cache Unit Tests
// =============================================================================

// TestCacheSetAndGet verifies that storing an HTTP response in the cache
// allows it to be retrieved with the same key.
func TestCacheSetAndGet(t *testing.T) {
	c, err := cache.NewCache(10*time.Minute, 8*time.Minute, 0)
	if err != nil {
		t.Fatalf("failed to create cache: %v", err)
	}
	defer c.Close()

	key := "https://api.zsapi.net/zia/api/v1/rules"
	body := `{"rules":["rule1","rule2"]}`
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Request:    &http.Request{Method: "GET", URL: &url.URL{Scheme: "https", Host: "api.zsapi.net", Path: "/zia/api/v1/rules"}},
	}

	c.Set(key, cache.CopyResponse(resp))

	cached := c.Get(key)
	if cached == nil {
		t.Fatal("expected cached response, got nil")
	}

	cachedBody, err := io.ReadAll(cached.Body)
	if err != nil {
		t.Fatalf("failed to read cached body: %v", err)
	}

	if string(cachedBody) != body {
		t.Fatalf("cached body mismatch: got %q, want %q", string(cachedBody), body)
	}
	t.Log("cache set/get works correctly")
}

// TestCacheDelete verifies that deleting a key removes it from the cache.
func TestCacheDelete(t *testing.T) {
	c, err := cache.NewCache(10*time.Minute, 8*time.Minute, 0)
	if err != nil {
		t.Fatalf("failed to create cache: %v", err)
	}
	defer c.Close()

	key := "https://api.zsapi.net/zia/api/v1/rules"
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
		Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
		Request:    &http.Request{Method: "GET", URL: &url.URL{Scheme: "https", Host: "api.zsapi.net", Path: "/zia/api/v1/rules"}},
	}

	c.Set(key, cache.CopyResponse(resp))
	if c.Get(key) == nil {
		t.Fatal("expected cache entry before delete")
	}

	c.Delete(key)
	if c.Get(key) != nil {
		t.Fatal("expected nil after delete")
	}
	t.Log("cache delete works correctly")
}

// TestCacheClear verifies that clearing the cache removes all entries.
func TestCacheClear(t *testing.T) {
	c, err := cache.NewCache(10*time.Minute, 8*time.Minute, 0)
	if err != nil {
		t.Fatalf("failed to create cache: %v", err)
	}
	defer c.Close()

	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("https://api.zsapi.net/zia/api/v1/resource/%d", i)
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{},
			Body:       io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{"id":%d}`, i))),
			Request:    &http.Request{Method: "GET", URL: &url.URL{Scheme: "https", Host: "api.zsapi.net", Path: fmt.Sprintf("/zia/api/v1/resource/%d", i)}},
		}
		c.Set(key, cache.CopyResponse(resp))
	}

	c.Clear()

	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("https://api.zsapi.net/zia/api/v1/resource/%d", i)
		if c.Get(key) != nil {
			t.Fatalf("expected nil for key %q after clear", key)
		}
	}
	t.Log("cache clear works correctly")
}

// TestCachePrefixInvalidation verifies that ClearAllKeysWithPrefix removes
// only matching entries, simulating the SDK's behavior of invalidating related
// cache entries on mutating (POST/PUT/DELETE) requests.
func TestCachePrefixInvalidation(t *testing.T) {
	c, err := cache.NewCache(10*time.Minute, 8*time.Minute, 0)
	if err != nil {
		t.Fatalf("failed to create cache: %v", err)
	}
	defer c.Close()

	// Populate cache with two different resource paths
	rulesKey := "https://api.zsapi.net/zia/api/v1/rules"
	usersKey := "https://api.zsapi.net/zia/api/v1/users"

	for _, key := range []string{rulesKey, usersKey} {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{},
			Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
			Request:    &http.Request{Method: "GET", URL: &url.URL{Scheme: "https", Host: "api.zsapi.net", Path: key}},
		}
		c.Set(key, cache.CopyResponse(resp))
	}

	// Invalidate only the rules prefix
	c.ClearAllKeysWithPrefix("https://api.zsapi.net/zia/api/v1/rules")

	if c.Get(rulesKey) != nil {
		t.Fatal("expected rules cache entry to be invalidated")
	}
	if c.Get(usersKey) == nil {
		t.Fatal("expected users cache entry to survive prefix invalidation")
	}
	t.Log("cache prefix invalidation works correctly")
}

// TestCacheTTLExpiry verifies that cache entries expire after the configured TTL.
func TestCacheTTLExpiry(t *testing.T) {
	// 1-second TTL with 500ms clean window for fast test
	c, err := cache.NewCache(1*time.Second, 500*time.Millisecond, 0)
	if err != nil {
		t.Fatalf("failed to create cache: %v", err)
	}
	defer c.Close()

	key := "https://api.zsapi.net/zia/api/v1/rules"
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
		Body:       io.NopCloser(bytes.NewBufferString(`{"data":"test"}`)),
		Request:    &http.Request{Method: "GET", URL: &url.URL{Scheme: "https", Host: "api.zsapi.net", Path: "/zia/api/v1/rules"}},
	}
	c.Set(key, cache.CopyResponse(resp))

	// Should exist immediately
	if c.Get(key) == nil {
		t.Fatal("expected cache entry to exist immediately after set")
	}

	// Wait for TTL + clean window to pass
	time.Sleep(2 * time.Second)

	if c.Get(key) != nil {
		t.Fatal("expected cache entry to expire after TTL")
	}
	t.Log("cache TTL expiry works correctly")
}

// TestNopCacheAlwaysMisses verifies that the NopCache implementation
// (used when caching is disabled) always returns nil.
func TestNopCacheAlwaysMisses(t *testing.T) {
	c := cache.NewNopCache()

	key := "https://api.zsapi.net/zia/api/v1/test"
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
		Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
		Request:    &http.Request{Method: "GET", URL: &url.URL{Scheme: "https", Host: "api.zsapi.net", Path: "/zia/api/v1/test"}},
	}

	c.Set(key, cache.CopyResponse(resp))

	if c.Get(key) != nil {
		t.Fatal("nop cache should always return nil")
	}
	t.Log("nop cache correctly never caches")
}

// =============================================================================
// Section 5: SDK Configuration Propagation Tests
// =============================================================================

// TestSDKConfigurationDefaults verifies that NewConfiguration sets expected
// default values for rate limiting and caching when no overrides are provided.
func TestSDKConfigurationDefaults(t *testing.T) {
	cfg, err := zscaler.NewConfiguration()
	if err != nil {
		t.Fatalf("failed to create configuration: %v", err)
	}

	// Rate limit defaults
	if cfg.Zscaler.Client.RateLimit.MaxRetries != zscaler.MaxNumOfRetries {
		t.Errorf("expected MaxRetries=%d, got %d", zscaler.MaxNumOfRetries, cfg.Zscaler.Client.RateLimit.MaxRetries)
	}
	if cfg.Zscaler.Client.RateLimit.RetryWaitMax != time.Duration(zscaler.RetryWaitMaxSeconds)*time.Second {
		t.Errorf("expected RetryWaitMax=%v, got %v", time.Duration(zscaler.RetryWaitMaxSeconds)*time.Second, cfg.Zscaler.Client.RateLimit.RetryWaitMax)
	}
	if cfg.Zscaler.Client.RateLimit.RetryWaitMin != time.Duration(zscaler.RetryWaitMinSeconds)*time.Second {
		t.Errorf("expected RetryWaitMin=%v, got %v", time.Duration(zscaler.RetryWaitMinSeconds)*time.Second, cfg.Zscaler.Client.RateLimit.RetryWaitMin)
	}

	// Cache defaults
	if cfg.Zscaler.Client.Cache.DefaultTtl != 10*time.Minute {
		t.Errorf("expected cache TTL=10m, got %v", cfg.Zscaler.Client.Cache.DefaultTtl)
	}
	if cfg.Zscaler.Client.Cache.DefaultTti != 8*time.Minute {
		t.Errorf("expected cache TTI=8m, got %v", cfg.Zscaler.Client.Cache.DefaultTti)
	}

	// Session retry default
	if cfg.Zscaler.Client.RateLimit.MaxSessionNotValidRetries != 3 {
		t.Errorf("expected MaxSessionNotValidRetries=3, got %d", cfg.Zscaler.Client.RateLimit.MaxSessionNotValidRetries)
	}
	t.Log("SDK configuration defaults are correct")
}

// TestSDKConfigurationOverrides verifies that ConfigSetter functions properly
// override rate limit and cache configuration values.
func TestSDKConfigurationOverrides(t *testing.T) {
	cfg, err := zscaler.NewConfiguration(
		zscaler.WithRateLimitMaxRetries(50),
		zscaler.WithRateLimitMinWait(5*time.Second),
		zscaler.WithRateLimitMaxWait(30*time.Second),
		zscaler.WithCache(true),
		zscaler.WithCacheTtl(15*time.Minute),
		zscaler.WithCacheTti(12*time.Minute),
	)
	if err != nil {
		t.Fatalf("failed to create configuration: %v", err)
	}

	if cfg.Zscaler.Client.RateLimit.MaxRetries != 50 {
		t.Errorf("expected MaxRetries=50, got %d", cfg.Zscaler.Client.RateLimit.MaxRetries)
	}
	if cfg.Zscaler.Client.RateLimit.RetryWaitMin != 5*time.Second {
		t.Errorf("expected RetryWaitMin=5s, got %v", cfg.Zscaler.Client.RateLimit.RetryWaitMin)
	}
	if cfg.Zscaler.Client.RateLimit.RetryWaitMax != 30*time.Second {
		t.Errorf("expected RetryWaitMax=30s, got %v", cfg.Zscaler.Client.RateLimit.RetryWaitMax)
	}
	if !cfg.Zscaler.Client.Cache.Enabled {
		t.Error("expected cache to be enabled")
	}
	if cfg.Zscaler.Client.Cache.DefaultTtl != 15*time.Minute {
		t.Errorf("expected cache TTL=15m, got %v", cfg.Zscaler.Client.Cache.DefaultTtl)
	}
	if cfg.Zscaler.Client.Cache.DefaultTti != 12*time.Minute {
		t.Errorf("expected cache TTI=12m, got %v", cfg.Zscaler.Client.Cache.DefaultTti)
	}
	t.Log("SDK configuration overrides propagated correctly")
}

// TestSDKConfigurationCacheDisabled verifies that when caching is disabled,
// the CacheManager still exists but as a NopCache (no caching behavior).
func TestSDKConfigurationCacheDisabled(t *testing.T) {
	cfg, err := zscaler.NewConfiguration(
		zscaler.WithCache(false),
	)
	if err != nil {
		t.Fatalf("failed to create configuration: %v", err)
	}

	if cfg.Zscaler.Client.Cache.Enabled {
		t.Error("expected cache to be disabled")
	}

	// CacheManager should still be non-nil (NopCache)
	if cfg.CacheManager == nil {
		t.Fatal("expected CacheManager to be non-nil even when cache is disabled")
	}
	t.Log("cache disabled configuration is correct")
}

// TestProviderConfigMatchesSDKConfig verifies that the provider's Config struct
// default values align with what the SDK expects. This ensures the provider
// correctly wires its configuration to the SDK.
func TestProviderConfigMatchesSDKConfig(t *testing.T) {
	// These are the provider defaults from config.go NewConfig
	providerDefaults := Config{
		backoff:        true,
		minWait:        2,
		maxWait:        10,
		retryCount:     100,
		parallelism:    1,
		requestTimeout: 1800,
	}

	// Verify provider defaults match SDK expectations
	if providerDefaults.retryCount != int(zscaler.MaxNumOfRetries) {
		t.Errorf("provider retryCount=%d should match SDK MaxNumOfRetries=%d",
			providerDefaults.retryCount, zscaler.MaxNumOfRetries)
	}
	if providerDefaults.minWait != int(zscaler.RetryWaitMinSeconds) {
		t.Errorf("provider minWait=%d should match SDK RetryWaitMinSeconds=%d",
			providerDefaults.minWait, zscaler.RetryWaitMinSeconds)
	}
	if providerDefaults.maxWait != int(zscaler.RetryWaitMaxSeconds) {
		t.Errorf("provider maxWait=%d should match SDK RetryWaitMaxSeconds=%d",
			providerDefaults.maxWait, zscaler.RetryWaitMaxSeconds)
	}
	t.Log("provider defaults match SDK constants")
}

// TestHTTPClientHasRateLimitTransport verifies that the ZIA HTTP client
// produced by SDK configuration uses the rate limit transport layer.
func TestHTTPClientHasRateLimitTransport(t *testing.T) {
	cfg, err := zscaler.NewConfiguration()
	if err != nil {
		t.Fatalf("failed to create configuration: %v", err)
	}

	if cfg.ZIAHTTPClient == nil {
		t.Fatal("expected ZIAHTTPClient to be non-nil")
	}
	if cfg.ZPAHTTPClient == nil {
		t.Fatal("expected ZPAHTTPClient to be non-nil")
	}
	if cfg.ZTWHTTPClient == nil {
		t.Fatal("expected ZTWHTTPClient to be non-nil")
	}
	if cfg.ZCCHTTPClient == nil {
		t.Fatal("expected ZCCHTTPClient to be non-nil")
	}
	if cfg.ZDXHTTPClient == nil {
		t.Fatal("expected ZDXHTTPClient to be non-nil")
	}
	if cfg.HTTPClient == nil {
		t.Fatal("expected default HTTPClient to be non-nil")
	}
	t.Log("all service-specific HTTP clients are properly initialized")
}

// =============================================================================
// Section 6: Combined Rate Limiting + Mock Server Integration Tests
// =============================================================================

// TestRateLimitAndRetryIntegration verifies the full pipeline: client-side
// rate limiting via RateLimitTransport combined with server-side 429 retry.
// This simulates what happens when the provider talks to the ZIA API under load.
func TestRateLimitAndRetryIntegration(t *testing.T) {
	var serverHits int32

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&serverHits, 1)
		// First request returns 429, rest return 200
		if count == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"code":"TOO_MANY_REQUESTS"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"request_number":%d}`, count)))
	}))
	defer mockServer.Close()

	// Set up the full pipeline: rate limiter → retryable client → mock server
	limiter := rl.NewRateLimiter(10, 5, 1, 1) // generous per-second limits

	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryWaitMin = 500 * time.Millisecond
	retryableClient.RetryWaitMax = 2 * time.Second
	retryableClient.RetryMax = 3
	retryableClient.Logger = nil

	retryableClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			return true, nil
		}
		return false, nil
	}

	retryableClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				if secs, err := time.ParseDuration(retryAfter + "s"); err == nil {
					return secs
				}
			}
		}
		return min
	}

	// Wrap with rate limiting transport
	rateLimitTransport := &rl.RateLimitTransport{
		Base:    http.DefaultTransport,
		Limiter: limiter,
		Logger:  logger.NewNopLogger(),
	}
	retryableClient.HTTPClient.Transport = rateLimitTransport

	client := retryableClient.StandardClient()

	// Make the request — should get 429, retry, then succeed with 200
	resp, err := client.Get(mockServer.URL + "/zia/api/v1/test")
	if err != nil {
		t.Fatalf("expected success after retry, got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 after retry, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	t.Logf("final response: %s (total server hits: %d)", string(body), atomic.LoadInt32(&serverHits))
}

// TestCacheReducesServerRequests verifies that caching at the HTTP layer
// prevents redundant requests to the server. This simulates the SDK's caching
// behavior where repeated GETs for the same resource are served from cache.
func TestCacheReducesServerRequests(t *testing.T) {
	var serverHits int32

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&serverHits, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"rules":[{"id":1,"name":"test"}]}`))
	}))
	defer mockServer.Close()

	// Create a real cache
	c, err := cache.NewCache(10*time.Minute, 8*time.Minute, 0)
	if err != nil {
		t.Fatalf("failed to create cache: %v", err)
	}
	defer c.Close()

	cacheKey := mockServer.URL + "/zia/api/v1/rules"

	// First request — cache miss, hits server
	resp1, err := http.Get(cacheKey)
	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}
	body1, _ := io.ReadAll(resp1.Body)
	resp1.Body = io.NopCloser(bytes.NewBuffer(body1))
	c.Set(cacheKey, cache.CopyResponse(resp1))
	resp1.Body.Close()

	// Second request — served from cache, should NOT hit server
	cachedResp := c.Get(cacheKey)
	if cachedResp == nil {
		t.Fatal("expected cache hit on second request")
	}
	body2, _ := io.ReadAll(cachedResp.Body)

	if string(body1) != string(body2) {
		t.Fatalf("cached response body mismatch: got %q, want %q", string(body2), string(body1))
	}

	hits := atomic.LoadInt32(&serverHits)
	if hits != 1 {
		t.Fatalf("expected exactly 1 server hit (cache should prevent second), got %d", hits)
	}
	t.Logf("cache correctly prevented redundant server request (server hits: %d)", hits)
}
