package shortener

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortner(t *testing.T) {
	t.Run("it short(ns on /shorturl", func(t *testing.T) {
		payload, _ := json.Marshal(map[string]string{
			"long_url": "http://www.google.com",
		})

		request, _ := http.NewRequest(http.MethodPost, "/s", bytes.NewBuffer(payload))
		request.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		GetShortUrl(res, request)
		if res.Code != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, res.Code)
		}
		var response map[string]string
		json.Unmarshal(res.Body.Bytes(), &response)
		if response["short_url"] != "http://localhost:8090/l?q=ed646a3334" {
			t.Errorf("Unexpected response body: %v", response)
		}
	})
}
