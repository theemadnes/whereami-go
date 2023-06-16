package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWhereami(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	getRoot(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if res.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected content-type to be application/json, got %v", res.Header.Get("Content-Type"))
	}
}
