package responses

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJsonResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       interface{}
	}{
		{
			name:       "success response",
			statusCode: http.StatusOK,
			data:       Response{Message: "success"},
		},
		{
			name:       "created response",
			statusCode: http.StatusCreated,
			data:       Response{Message: "created"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			JsonResponse(w, tt.statusCode, tt.data)

			if w.Code != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, w.Code)
			}

			if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("expected content type application/json, got %s", contentType)
			}

			var response Response
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			expectedResponse, ok := tt.data.(Response)
			if !ok {
				t.Fatalf("expected data to be of type Response")
			}

			if response.Message != expectedResponse.Message {
				t.Errorf("expected message %q, got %q", expectedResponse.Message, response.Message)
			}
		})
	}
}

func TestJsonErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		err        error
	}{
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			err:        errors.New("invalid input"),
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("something went wrong"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			JsonErrorResponse(w, tt.statusCode, tt.err)

			if w.Code != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, w.Code)
			}

			if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("expected content type application/json, got %s", contentType)
			}
		})
	}
}
