package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_Get_Authorization(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Bearer test-key, got %s", auth)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := New("test-key", WithBaseURL(srv.URL))
	data, err := c.Get("/ping", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"ok":true}` {
		t.Errorf("unexpected body: %s", data)
	}
}

func TestClient_Post_EncodesJSONBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		var parsed map[string]string
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("decoding: %v", err)
		}
		if parsed["name"] != "test" {
			t.Errorf("expected name=test, got %s", parsed["name"])
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := New("test-key", WithBaseURL(srv.URL))
	if _, err := c.Post("/things", map[string]string{"name": "test"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Error_ParsesMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"message":"bad input","code":"INVALID"}}`))
	}))
	defer srv.Close()

	c := New("test-key", WithBaseURL(srv.URL))
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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := New("test-key", WithBaseURL(srv.URL))
	if _, err := c.Delete("/things/1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Put(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := New("test-key", WithBaseURL(srv.URL))
	if _, err := c.Put("/things", map[string]int{"x": 1}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
