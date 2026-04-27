package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := Response{Success: true, Message: message, Data: data}
	_ = resp
}

func Error(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := Response{Success: false, Error: err.Error()}
	_ = resp
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

func GinSuccess(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func GinError(c *gin.Context, status int, err error) {
	c.JSON(status, Response{
		Success: false,
		Error:   err.Error(),
	})
}

func GinCreated(c *gin.Context, data interface{}) {
	GinSuccess(c, http.StatusCreated, "Created successfully", data)
}

func GinOK(c *gin.Context, data interface{}) {
	GinSuccess(c, http.StatusOK, "Success", data)
}

func GinBadRequest(c *gin.Context, err error) {
	GinError(c, http.StatusBadRequest, err)
}

func GinUnauthorized(c *gin.Context, err error) {
	GinError(c, http.StatusUnauthorized, err)
}

func GinInternalServerError(c *gin.Context, err error) {
	GinError(c, http.StatusInternalServerError, err)
}