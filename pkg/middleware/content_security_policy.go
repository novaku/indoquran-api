package middleware

import (
	"net/http"

	"indoquran-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Middleware CSP
func ContentSecurityPolicy() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set header Content-Security-Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'; report-uri /csp-report;")

		// Lanjutkan ke handler berikutnya
		c.Next()
	}
}

// Handler untuk menerima laporan CSP pelanggaran
func CspReportHandler(c *gin.Context) {
	var report map[string]interface{}

	// Terima JSON laporan pelanggaran
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report format"})
		return
	}

	// Cetak laporan pelanggaran (ini bisa dikirim ke sistem logging)
	logger.WriteLog(logger.LogLevelError, "CSP Violation Report: %#v", report)

	// Berikan respons OK
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
