package controllers

import (
	"math"
	"net/http"
	"strconv"

	"indoquran-api/internal/services/search"
	"indoquran-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ResultJsonFormat search full text search result
func SearchHandler(c *gin.Context) {
	// Parse query parameters
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	juz, _ := strconv.Atoi(c.DefaultQuery("juz", "0"))
	surat, _ := strconv.Atoi(c.DefaultQuery("surat", "0"))
	rowsPerPage, _ := strconv.Atoi(c.DefaultQuery("n", "10"))

	results, totalCount, aggreggate, err := search.FullTextSearch(query, juz, surat, page, rowsPerPage)
	if err != nil {
		logger.WriteLog(logger.LogLevelError, "Error fetching from database: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(rowsPerPage)))

	res := ResultJsonFormat{
		Aggregate: aggreggate,
		Pagination: &Pagination{
			CurrentPage: page,
			RowsPerPage: rowsPerPage,
			TotalPages:  totalPages,
			TotalRows:   int(totalCount),
		},
		Results: results,
	}

	WriteResponse(c, res, err)
}
