package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents the response structure
type Response struct {
	Version string      `json:"version"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

// WriteResponse writes the response to the client as default output
func WriteResponse(c *gin.Context, data interface{}, err error) {
	httpStatus := http.StatusOK
	response := Response{
		Version: "1.0",
		Data:    data,
		Error:   "",
	}

	if err != nil {
		httpStatus = http.StatusBadRequest
		response.Error = err.Error()
	}

	c.JSON(httpStatus, response)
}
