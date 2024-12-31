package responses

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message interface{} `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func JsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if str, ok := data.(string); ok {
		data = Response{Message: str}
	}

	json.NewEncoder(w).Encode(data)
}

func JsonErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	JsonResponse(w, statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: err.Error(),
	})
}
