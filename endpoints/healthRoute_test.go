package endpoints

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewHealthHanderContext(t *testing.T) {
	ctx := context.Background()
	authorization := true
	health := NewHealthHandlerContext(ctx, authorization)
	if health == nil {
		t.Errorf("Expected a new HealthRoute, but got nil")
	}

	if health != nil {

		if health.ctx != nil {
			t.Errorf("Expected context to be %v, but got %v", ctx, health.ctx)
		}

		if health.authorization != authorization {
			t.Errorf("Expected authorization to be %v, but got %v", authorization, health.authorization)
		}

	}

}

func TestHealthHTTPHandler(t *testing.T) {

	health := NewHealthHandlerContext(context.Background(), false)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(health.HealthHTTPHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"status":"OK"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func TestHealthHTTPHandlerInvalidMethod(t *testing.T) {

	health := NewHealthHandlerContext(context.Background(), false)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(health.HealthHTTPHandler)
	handler.ServeHTTP(rr, req)
	msg := strings.TrimSpace(rr.Body.String())

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `Invalid request method.`
	if msg != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			msg, expected)
	}

}

func TestHealthHTTPHandlerInvalidToken(t *testing.T) {

	health := NewHealthHandlerContext(context.Background(), true)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(health.HealthHTTPHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `Invalid credentials`
	msg := strings.TrimSpace(rr.Body.String())
	if msg != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func TestHealthHTTPHandlerValidToken(t *testing.T) {

	health := NewHealthHandlerContext(context.Background(), true)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header = http.Header{
		"Authorization": []string{"Bearer 931A3B02-8DCC-543F-A1B2-69423D1A0B94"},
	}

	handler := http.HandlerFunc(health.HealthHTTPHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"status":"OK"}`
	msg := strings.TrimSpace(rr.Body.String())
	if msg != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
