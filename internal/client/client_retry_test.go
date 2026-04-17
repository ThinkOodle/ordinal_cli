package client

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func okResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func rateLimitResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
			"Retry-After":  []string{"0"},
		},
		Body: io.NopCloser(strings.NewReader(`{"error":{"message":"rate limited"}}`)),
	}
}

// withInstantBackoff neutralizes the exponential sleep so retry tests run fast.
func withInstantBackoff(t *testing.T) {
	t.Helper()
	orig := initialBackoffValue
	initialBackoffValue = 0
	t.Cleanup(func() { initialBackoffValue = orig })
}

func TestClient_DoesNotRetryPOSTOnTransportError(t *testing.T) {
	withInstantBackoff(t)

	var attempts int32
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		atomic.AddInt32(&attempts, 1)
		return nil, errors.New("connection reset")
	})
	c := New("k",
		WithBaseURL("http://api.test"),
		WithHTTPClient(&http.Client{Transport: transport}),
	)

	if _, err := c.Post("/things", map[string]string{"name": "x"}); err == nil {
		t.Fatal("expected error")
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Errorf("POST should not retry on transport error; got %d attempts", got)
	}
}

func TestClient_DoesNotRetryPATCHOnTransportError(t *testing.T) {
	withInstantBackoff(t)

	var attempts int32
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		atomic.AddInt32(&attempts, 1)
		return nil, errors.New("connection reset")
	})
	c := New("k",
		WithBaseURL("http://api.test"),
		WithHTTPClient(&http.Client{Transport: transport}),
	)

	if _, err := c.Patch("/things/1", map[string]string{"name": "x"}); err == nil {
		t.Fatal("expected error")
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Errorf("PATCH should not retry on transport error; got %d attempts", got)
	}
}

func TestClient_RetriesGETOnTransportError(t *testing.T) {
	withInstantBackoff(t)

	var attempts int32
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		n := atomic.AddInt32(&attempts, 1)
		if n < 3 {
			return nil, errors.New("connection reset")
		}
		return okResponse(`{"ok":true}`), nil
	})
	c := New("k",
		WithBaseURL("http://api.test"),
		WithHTTPClient(&http.Client{Transport: transport}),
	)

	data, err := c.Get("/things", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"ok":true}` {
		t.Errorf("unexpected body: %s", data)
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Errorf("expected 3 GET attempts, got %d", got)
	}
}

func TestClient_RetriesPOSTOn429(t *testing.T) {
	withInstantBackoff(t)

	var attempts int32
	var bodies []string
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		n := atomic.AddInt32(&attempts, 1)
		if r.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = io.Copy(buf, r.Body)
			bodies = append(bodies, strings.TrimSpace(buf.String()))
		}
		if n < 2 {
			return rateLimitResponse(), nil
		}
		return okResponse(`{"ok":true}`), nil
	})
	c := New("k",
		WithBaseURL("http://api.test"),
		WithHTTPClient(&http.Client{Transport: transport}),
	)

	if _, err := c.Post("/things", map[string]string{"name": "x"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Errorf("expected 2 POST attempts (429 then OK), got %d", got)
	}
	if len(bodies) != 2 || bodies[0] != bodies[1] {
		t.Errorf("expected identical POST bodies on retry, got %v", bodies)
	}
}

func TestParseRetryAfter(t *testing.T) {
	now := time.Date(2026, 4, 17, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		header string
		want   time.Duration
		ok     bool
	}{
		{"seconds", "5", 5 * time.Second, true},
		{"zero seconds ignored", "0", 0, false},
		{"negative seconds ignored", "-3", 0, false},
		{"http-date future", now.Add(10 * time.Second).UTC().Format(http.TimeFormat), 10 * time.Second, true},
		{"http-date past ignored", now.Add(-1 * time.Second).UTC().Format(http.TimeFormat), 0, false},
		{"garbage ignored", "soon", 0, false},
		{"empty ignored", "", 0, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := parseRetryAfter(tc.header, now)
			if ok != tc.ok {
				t.Fatalf("ok: got %v, want %v", ok, tc.ok)
			}
			if ok && got != tc.want {
				t.Errorf("duration: got %v, want %v", got, tc.want)
			}
		})
	}
}

// TestClient_RateLimitBackoffSurvivesZeroRetryAfter locks in that a server
// returning "Retry-After: 0" does NOT collapse the retry loop into a
// zero-delay hammer. We set a small-but-nonzero initial backoff and assert
// that the elapsed time between retries reflects that backoff rather than
// racing through them instantly.
func TestClient_RateLimitBackoffSurvivesZeroRetryAfter(t *testing.T) {
	orig := initialBackoffValue
	initialBackoffValue = 25 * time.Millisecond
	t.Cleanup(func() { initialBackoffValue = orig })

	var attempts int32
	var timestamps []time.Time
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		timestamps = append(timestamps, time.Now())
		n := atomic.AddInt32(&attempts, 1)
		if n < 3 {
			return rateLimitResponse(), nil
		}
		return okResponse(`{"ok":true}`), nil
	})
	c := New("k",
		WithBaseURL("http://api.test"),
		WithHTTPClient(&http.Client{Transport: transport}),
	)

	if _, err := c.Get("/things", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Fatalf("expected 3 attempts, got %d", got)
	}
	// First retry should respect initialBackoffValue; second should double it.
	// Use generous lower bounds to avoid flakiness but still catch the bug,
	// which produced zero-gap retries.
	if gap := timestamps[1].Sub(timestamps[0]); gap < 20*time.Millisecond {
		t.Errorf("first retry gap %v collapsed below initial backoff", gap)
	}
	if gap := timestamps[2].Sub(timestamps[1]); gap < 40*time.Millisecond {
		t.Errorf("second retry gap %v did not double initial backoff", gap)
	}
}

// withMaxTotalDuration shrinks the overall retry budget for a single test so
// we can observe the budget-exceeded path without waiting seconds.
func withMaxTotalDuration(t *testing.T, d time.Duration) {
	t.Helper()
	orig := maxTotalDuration
	maxTotalDuration = d
	t.Cleanup(func() { maxTotalDuration = orig })
}

// TestClient_RetryLoopIsBoundedByTotalDuration guards the fix for the
// unbounded-retry-time bug: a long run of 429s combined with exponential
// backoff could make a single command blow well past its advertised 30s
// timeout. The whole do() call must now be capped by maxTotalDuration.
// We set the budget very small and use a non-trivial initial backoff so
// that after one or two retries the budget runs out; the request must
// fail well inside the budget window rather than plowing through all 5
// retries.
func TestClient_RetryLoopIsBoundedByTotalDuration(t *testing.T) {
	withMaxTotalDuration(t, 150*time.Millisecond)
	orig := initialBackoffValue
	initialBackoffValue = 80 * time.Millisecond
	t.Cleanup(func() { initialBackoffValue = orig })

	var attempts int32
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		atomic.AddInt32(&attempts, 1)
		return rateLimitResponse(), nil
	})
	c := New("k",
		WithBaseURL("http://api.test"),
		WithHTTPClient(&http.Client{Transport: transport}),
	)

	start := time.Now()
	_, err := c.Get("/things", nil)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected error when budget exhausted by repeated 429s")
	}
	// 80ms + 160ms > 150ms budget, so we should bail after at most 2
	// attempts. Give a generous cap to absorb scheduler noise while still
	// catching the regression (which would run all 6 attempts).
	if got := atomic.LoadInt32(&attempts); got > 2 {
		t.Errorf("expected retries to stop inside budget; got %d attempts", got)
	}
	// The total call must not drag on far beyond the budget. If the sleep
	// is no longer interruptible / pre-checked, this test will catch it.
	if elapsed > 400*time.Millisecond {
		t.Errorf("do() took %v; expected to stay near the 150ms budget", elapsed)
	}
}

// TestClient_OversizedRetryAfterDoesNotBlock guards against honoring a server
// Retry-After that would exceed the overall budget. Without the pre-sleep
// budget check, a Retry-After of "60" would make the CLI block for a full
// minute before ultimately timing out — this test asserts we bail promptly
// instead.
func TestClient_OversizedRetryAfterDoesNotBlock(t *testing.T) {
	withMaxTotalDuration(t, 100*time.Millisecond)
	orig := initialBackoffValue
	initialBackoffValue = 1 * time.Millisecond
	t.Cleanup(func() { initialBackoffValue = orig })

	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
				// 60 seconds — dwarfs the 100ms test budget.
				"Retry-After": []string{"60"},
			},
			Body: io.NopCloser(strings.NewReader(`{"error":{"message":"rate limited"}}`)),
		}, nil
	})
	c := New("k",
		WithBaseURL("http://api.test"),
		WithHTTPClient(&http.Client{Transport: transport}),
	)

	start := time.Now()
	_, err := c.Get("/things", nil)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected error from 429 storm")
	}
	// Must bail well before Retry-After would elapse. Anything longer than
	// a handful of budget windows signals we actually slept on Retry-After.
	if elapsed > 1*time.Second {
		t.Errorf("call blocked for %v on an oversized Retry-After; expected near-immediate bail", elapsed)
	}
}

func TestClient_IsIdempotent(t *testing.T) {
	for _, m := range []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete, http.MethodOptions} {
		if !isIdempotent(m) {
			t.Errorf("%s should be idempotent", m)
		}
	}
	for _, m := range []string{http.MethodPost, http.MethodPatch} {
		if isIdempotent(m) {
			t.Errorf("%s should not be idempotent", m)
		}
	}
}
