package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{
		Success: false,
		Error:   http.StatusText(status),
		Message: message,
	})
}

func WriteSuccess(w http.ResponseWriter, data interface{}, message ...string) {
	response := SuccessResponse{Success: true, Data: data}
	if len(message) > 0 {
		response.Message = message[0]
	}
	WriteJSON(w, http.StatusOK, response)
}