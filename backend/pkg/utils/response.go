package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// RespondWithError sends an error response to the client with the provided status code, message, and optional error details.
func RespondWithError(c *gin.Context, statusCode int, message string, err error) {
	response := ErrorResponse{
		Error: http.StatusText(statusCode),
	}

	if message != "" {
		response.Message = message
	} else if err != nil {
		response.Message = err.Error()
	}

	c.JSON(statusCode, response)
}

// RespondWithSuccess sends a JSON success response with the given status code, data payload, and optional message.
func RespondWithSuccess(c *gin.Context, statusCode int, data interface{}, message string) {
	response := SuccessResponse{
		Data:    data,
		Message: message,
	}

	c.JSON(statusCode, response)
}
