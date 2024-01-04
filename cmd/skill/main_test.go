package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhook(t *testing.T) {
	handler := http.HandlerFunc(webhook)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	successBody := `{
        "response": {
            "text": "Извините, я пока ничего не умею"
        },
        "version": "1.0"
    }`

	tests := []struct {
		method       string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodGet, expectedCode: http.StatusMethodNotAllowed, expectedBody: ""},
		{method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed, expectedBody: ""},
		{method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed, expectedBody: ""},
		{method: http.MethodPost, expectedCode: http.StatusOK, expectedBody: successBody},
	}

	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = test.method
			req.URL = srv.URL

			resp, err := req.Send()
			assert.NoErrorf(t, err, "Error making HTTP request")

			assert.Equal(t, test.expectedCode, resp.StatusCode(), "Response code didn't match expected")
			if test.expectedBody != "" {
				assert.JSONEq(t, test.expectedBody, string(resp.Body()), "Response body didn't match expected")
			}
		})
	}
}
