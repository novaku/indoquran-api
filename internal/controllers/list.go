package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"indoquran.web.id/internal/services/list"
)

// ListSurat list of surat
func ListSurat(c *gin.Context) {
	iSurat := list.NewSurat()

	suratID := c.DefaultQuery("surat", "")
	surats, err := iSurat.GetSuratList(suratID)

	WriteResponse(c, surats, err)
}

// ListAyatInSurat list of ayat in surat
func ListAyatInSurat(c *gin.Context) {
	suratID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("n", "10"))

	iAyat := list.NewAyat()

	ayatList, err := iAyat.GetAyatList(suratID, page, pageSize)

	WriteResponse(c, ayatList, err)
}
