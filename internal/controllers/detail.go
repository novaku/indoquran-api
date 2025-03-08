package controllers

import (
	"github.com/gin-gonic/gin"
	"indoquran.web.id/internal/services/detail"
)

// DetailAyat handles the GET request for retrieving ayat details
func DetailAyat(c *gin.Context) {
	iDetail := detail.NewDetail()

	ayatID := c.Param("id")
	ayat, err := iDetail.GetAyat(ayatID)

	WriteResponse(c, ayat, err)
}
