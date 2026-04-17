package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/client"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func newTestClient(fn roundTripFunc) *client.Client {
	return client.New(
		"test-key",
		client.WithBaseURL("http://api.test"),
		client.WithHTTPClient(&http.Client{Transport: fn}),
	)
}

func jsonResponse(t testing.TB, status int, body interface{}) *http.Response {
	t.Helper()

	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshaling response: %v", err)
	}

	return &http.Response{
		StatusCode: status,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(payload)),
	}
}

func rawResponse(t testing.TB, status int, payload []byte) *http.Response {
	t.Helper()
	return &http.Response{
		StatusCode: status,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(payload)),
	}
}
