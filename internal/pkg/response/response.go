package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	JSON(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, status int, err error) {
	JSON(w, status, Response{
		Success: false,
		Error:   err.Error(),
	})
}

func Created(w http.ResponseWriter, data interface{}) {
	Success(w, http.StatusCreated, "Created successfully", data)
}

func OK(w http.ResponseWriter, data interface{}) {
	Success(w, http.StatusOK, "Success", data)
}

func BadRequest(w http.ResponseWriter, err error) {
	Error(w, http.StatusBadRequest, err)
}

func Unauthorized(w http.ResponseWriter, err error) {
	Error(w, http.StatusUnauthorized, err)
}

func InternalServerError(w http.ResponseWriter, err error) {
	Error(w, http.StatusInternalServerError, err)
}