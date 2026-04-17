package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// jsonResponse builds an in-memory *http.Response for a RoundTripper fake.
func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestClient_Get_Authorization(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
			t.Errorf("expected Bearer test-key, got %s", auth)
		}
		return jsonResponse(http.StatusOK, `{"ok":true}`), nil
	})
	c := New("test-key",
		WithBaseURL("http://api.test"),
		WithHTTPClient(&http.Client{Transport: transport}),
	)

	data, err := c.Get("/ping", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"ok":true}` {
		t.Errorf("unexpected body: %s", data)
	}
}

func TestClient_Post_EncodesJSONBody(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", got)
		}
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, r.Body); err != nil {
			t.Fatalf("reading body: %v", err)
		}
		var parsed map[string]string
		if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
			t.Fatalf("decoding: %v", err)
		}
		if parsed["name"] != "test" {
			t.Errorf("expected name=test, got %s", parsed["name"])
		}
		return jsonResponse(http.StatusOK, `{}`), nil
	})
	c := New("test-key",
		WithBaseURL("http://api.test"),
		WithHTTPClient(&http.Client{Transport: transport}),
	)

	if _, err := c.Post("/things", map[string]string{"name": "test"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Error_ParsesMessage(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusBadRequest, `{"error":{"message":"bad input","code":"INVALID"}}`), nil
	})
	c := New("test-key",
		WithBaseURL("http://api.test"),
		WithHTTPClient(&http.Client{Transport: transport}),
	)

	_, err := c.Get("/bad", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 || apiErr.Message != "bad input" || apiErr.Code != "INVALID" {
		t.Errorf("unexpected error fields: %+v", apiErr)
	}
	if !strings.Contains(apiErr.Error(), "bad input") {
		t.Errorf("expected error text to contain 'bad input', got %s", apiErr.Error())
	}
}

func TestClient_Delete(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		return jsonResponse(http.StatusOK, `{}`), nil
	})
	c := New("test-key",
		WithBaseURL("http://api.test"),
		WithHTTPClient(&http.Client{Transport: transport}),
	)

	if _, err := c.Delete("/things/1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Put(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		return jsonResponse(http.StatusOK, `{}`), nil
	})
	c := New("test-key",
		WithBaseURL("http://api.test"),
		WithHTTPClient(&http.Client{Transport: transport}),
	)

	if _, err := c.Put("/things", map[string]int{"x": 1}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
