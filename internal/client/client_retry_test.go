package client

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
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
