package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseWriterInterface defines the contract for writing responses
type ResponseWriterInterface interface {
	WriteResponse(c *gin.Context, data interface{}, err error)
}

// ResponseFormatterInterface defines the contract for formatting responses
type ResponseFormatterInterface interface {
	FormatResponse(data interface{}, err error) (int, Response)
}

// Response represents the response structure
type Response struct {
	Version string      `json:"version"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

// DefaultResponseWriter implements ResponseWriterInterface
type DefaultResponseWriter struct {
	formatter ResponseFormatterInterface
}

// DefaultResponseFormatter implements ResponseFormatterInterface
type DefaultResponseFormatter struct{}

// NewDefaultResponseWriter creates a new instance of DefaultResponseWriter
func NewDefaultResponseWriter(formatter ResponseFormatterInterface) ResponseWriterInterface {
	return &DefaultResponseWriter{
		formatter: formatter,
	}
}

// NewDefaultResponseFormatter creates a new instance of DefaultResponseFormatter
func NewDefaultResponseFormatter() ResponseFormatterInterface {
	return &DefaultResponseFormatter{}
}

// FormatResponse implements ResponseFormatterInterface
func (f *DefaultResponseFormatter) FormatResponse(data interface{}, err error) (int, Response) {
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

	return httpStatus, response
}

// WriteResponse implements ResponseWriterInterface
func (w *DefaultResponseWriter) WriteResponse(c *gin.Context, data interface{}, err error) {
	httpStatus, response := w.formatter.FormatResponse(data, err)
	c.JSON(httpStatus, response)
}

// GetDefaultResponseWriter returns a new DefaultResponseWriter with default implementations
func GetDefaultResponseWriter() ResponseWriterInterface {
	return NewDefaultResponseWriter(NewDefaultResponseFormatter())
}

// WriteResponse writes the response to the client as default output
func WriteResponse(c *gin.Context, data interface{}, err error) {
	GetDefaultResponseWriter().WriteResponse(c, data, err)
}
