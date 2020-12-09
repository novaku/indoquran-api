package helpers

import "github.com/gin-gonic/gin"

// GetCurrentURL : get current url
func GetCurrentURL(c *gin.Context) string {
	return c.Request.URL.RequestURI()
}
