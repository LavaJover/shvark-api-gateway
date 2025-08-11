package common

import (
	
	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, HTTPError{
		Code:    code,
		Message: message,
	})
}