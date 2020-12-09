package quran

import (
	"net/http"
	"strconv"

	"bitbucket.org/indoquran-api/src/handlers"
	"bitbucket.org/indoquran-api/src/helpers"
	"bitbucket.org/indoquran-api/src/models/quran"
	"bitbucket.org/indoquran-api/src/models/quran/format"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func findSurat(c *gin.Context) {
	rowsPerPage, _ := strconv.ParseInt(c.DefaultQuery("rowsperpage", ""), 10, 64)
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	sortBy := c.DefaultQuery("sortby", "nomor")
	descending, _ := strconv.ParseBool(c.DefaultQuery("descending", "false"))
	id, _ := strconv.Atoi(c.Param("id"))

	// for pagination
	findoptions := options.Find()
	if rowsPerPage != 0 {
		findoptions.SetLimit(rowsPerPage)
		findoptions.SetSkip(rowsPerPage * (page - 1))
	}

	findoptions.SetSort(bson.D{{Key: "nomor", Value: 1}})

	findFilter := bson.M{}
	if id != 0 {
		findFilter = bson.M{"nomor": id}
	}

	cursor, err = suratCollection.Find(c, findFilter, findoptions)
	if err != nil {
		mlog.Error(err)
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Error Get Surat Data", err)
		return
	}

	// count all document by query
	rowsNumber, err := suratCollection.CountDocuments(c, findFilter)
	if err != nil {
		mlog.Error(err)
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Error Get Surat Data", err)
		return
	}

	pagination = format.Pagination{
		SortBy:      sortBy,
		Descending:  descending,
		Page:        page,
		RowsPerPage: rowsPerPage,
		RowsNumber:  rowsNumber,
	}
}

// GetSurats : Get All Surat Endpoint
func GetSurats(c *gin.Context) {
	mlog.Info("Get Surat API, url : %+v", helpers.GetCurrentURL(c))
	var (
		surat  quran.Surat
		surats []format.Surat
	)

	findSurat(c)

	defer cursor.Close(c)

	for cursor.Next(c) {
		if err := cursor.Decode(&surat); err != nil {
			mlog.Error(err)
		}

		result := format.Surat{
			ID:         surat.ID,
			Nomor:      surat.Nomor,
			Nama:       surat.Nama,
			Asma:       surat.Asma,
			JumlahAyat: surat.JumlahAyat,
			Arti:       surat.Arti,
			Keterangan: surat.Keterangan,
		}

		surats = append(surats, result)
	}

	if err := cursor.Err(); err != nil {
		mlog.Error(err)
		handlers.DefaultResponse(c, http.StatusBadRequest, "Error Get Surat Data", err)
		return
	}

	handlers.DefaultResponse(c, http.StatusOK, "Success Get Surat Data", &format.Surats{
		Surats:     surats,
		Pagination: pagination,
	})
	return
}
