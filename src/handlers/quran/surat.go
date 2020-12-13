package quran

import (
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/indoquran-api/src/config"
	"bitbucket.org/indoquran-api/src/handlers"
	"bitbucket.org/indoquran-api/src/helpers"
	"bitbucket.org/indoquran-api/src/models/quran"
	"bitbucket.org/indoquran-api/src/models/quran/format"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
	"github.com/vmihailenco/msgpack"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func findSurat(c *gin.Context, rowsPerPage, page int64, sortBy string, descending bool, id int) {
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

	rowsPerPage, _ := strconv.ParseInt(c.DefaultQuery("rowsperpage", ""), 10, 64)
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	sortBy := c.DefaultQuery("sortby", "nomor")
	descending, _ := strconv.ParseBool(c.DefaultQuery("descending", "false"))
	id, _ := strconv.Atoi(c.Param("id"))
	result := format.Surats{}
	redisKey := redisKeyGeneratorSurat(sortBy, rowsPerPage, page, descending, id)

	mlog.Info("Search surat API, url : %+v", helpers.GetCurrentURL(c))

	val, err := cache.Get(c, redisKey).Result()
	if err != nil || val == "" {
		findSurat(c, rowsPerPage, page, sortBy, descending, id)

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

		result.Surats = surats
		result.Pagination = pagination

		// redis set
		b, err := msgpack.Marshal(&result)
		if err != nil {
			mlog.Error(err)
		}

		ttl := time.Duration(config.Config.Cache.TTL) * time.Hour

		mlog.Info("SET redis key : %s", redisKey)

		set, err := cache.SetNX(c, redisKey, string(b), ttl).Result()
		if !set || err != nil {
			mlog.Error(err)
		}
	}

	if val != "" {
		mlog.Info("GET redis key : %s", redisKey)

		err = msgpack.Unmarshal([]byte(val), &result)
		if err != nil {
			panic(err)
		}
	}

	handlers.DefaultResponse(c, http.StatusOK, "Success Get Surat Data", result)
	return
}
