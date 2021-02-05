package quran

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func findAyat(c *gin.Context, search, sortBy string, surat, juz int, rowsPerPage, page int64, descending bool) {
	// for pagination
	findoptions := options.Find()
	findoptions.SetLimit(rowsPerPage)
	findoptions.SetSkip(rowsPerPage * (page - 1))
	findoptions.SetSort(bson.D{
		{Key: "surat", Value: 1},
		{Key: "nomor", Value: 1},
	})

	// search with AND operator
	if search != "" {
		strRes := []string{}
		strArr := strings.Split(search, " ")
		for _, txt := range strArr {
			strRes = append(strRes, "\""+txt+"\"")
		}
		search = strings.Join(strRes, " ")
	}

	findFilter := bson.M{}
	if search != "" && surat == 0 {
		findFilter = bson.M{"$text": bson.M{"$search": search}}

		if juz != 0 {
			findFilter = bson.M{
				"$and": []bson.M{
					{"$text": bson.M{"$search": search}},
					{"juz": juz},
				},
			}
		}
	} else if search != "" && surat != 0 {
		findFilter = bson.M{
			"$and": []bson.M{
				{"$text": bson.M{"$search": search}},
				{"surat": surat},
			},
		}

		if juz != 0 {
			findFilter = bson.M{
				"$and": []bson.M{
					{"$text": bson.M{"$search": search}},
					{"surat": surat},
					{"juz": juz},
				},
			}
		}
	} else if search == "" && surat != 0 {
		findFilter = bson.M{"surat": surat}

		if juz != 0 {
			findFilter = bson.M{
				"$and": []bson.M{
					{"surat": surat},
					{"juz": juz},
				},
			}
		}
	} else if juz != 0 {
		findFilter = bson.M{"juz": juz}
	}

	mlog.Info("filter : %+v", findFilter)

	cursor, err = ayatCollection.Find(c, findFilter, findoptions)
	if err != nil {
		mlog.Error(err)
	}

	// count all document by query
	rowsNumber, err := ayatCollection.CountDocuments(c, findFilter)
	if err != nil {
		mlog.Error(err)
	}

	pagination = format.Pagination{
		SortBy:      sortBy,
		Descending:  descending,
		Page:        page,
		RowsPerPage: rowsPerPage,
		RowsNumber:  rowsNumber,
	}
}

func getCatatans(c *gin.Context, surat, ayat int) ([]format.Catatan, error) {
	catatan := quran.Catatan{}
	catatans := []format.Catatan{}

	findoptions := options.Find()
	findoptions.SetSort(bson.D{{Key: "nomor", Value: 1}})

	cursor, err := catatanCollection.Find(c, bson.M{"surat": surat, "ayat": ayat}, findoptions)
	if err != nil {
		mlog.Error(err)
		return nil, err
	}

	defer cursor.Close(c)

	if err := cursor.Err(); err != nil {
		mlog.Error(err)
		return nil, err
	}

	for cursor.Next(c) {
		if err := cursor.Decode(&catatan); err != nil {
			mlog.Error(err)
			return nil, err
		}

		cat := format.Catatan{
			ID:        catatan.Nomor,
			TeksQuran: catatan.TeksCatatan,
		}

		catatans = append(catatans, cat)
	}

	return catatans, nil
}

func getSurat(c *gin.Context, suratID int) (*quran.Surat, error) {
	surat := quran.Surat{}

	if err := suratCollection.FindOne(c, bson.M{"nomor": suratID}).Decode(&surat); err != nil {
		return nil, err
	}

	return &surat, nil
}

func createIndexAyat(c *gin.Context) {
	index := mongo.IndexModel{Keys: bson.M{"txt_id": 1}}
	if _, err := ayatCollection.Indexes().CreateOne(c, index); err != nil {
		log.Println("Could not create index:", err)
	}
}

// GetSearchAyats : endpoint for get and search ayats
func GetSearchAyats(c *gin.Context) {
	rowsPerPage, _ := strconv.ParseInt(c.DefaultQuery("rowsperpage", "5"), 10, 64)
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	sortBy := c.DefaultQuery("sortby", "nomor")
	descending, _ := strconv.ParseBool(c.DefaultQuery("descending", "false"))
	search := c.Param("searchText")
	surat, _ := strconv.Atoi(c.DefaultQuery("surat", ""))
	juz, _ := strconv.Atoi(c.DefaultQuery("juz", ""))
	result := format.Ayats{}
	redisKey := redisKeyGeneratorAyat(search, sortBy, surat, 0, juz, rowsPerPage, page, descending)

	mlog.Info("Search ayat API, url : %+v", helpers.GetCurrentURL(c))

	val, err := cache.Get(c, redisKey).Result()
	if err != nil || val == "" {
		var ayat quran.Ayat
		ayats := make([]format.Ayat, 0)

		createIndexAyat(c)
		findAyat(c, search, sortBy, surat, juz, rowsPerPage, page, descending)

		defer cursor.Close(c)

		if err := cursor.Err(); err != nil {
			mlog.Error(err)
			handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed Get All Ayat", err.Error())
			return
		}

		for cursor.Next(c) {
			if err := cursor.Decode(&ayat); err != nil {
				mlog.Error(err)
				handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed Decode Result", err.Error())
				return
			}

			catatans, err := getCatatans(c, ayat.Surat, ayat.Nomor)
			if err != nil {
				mlog.Error(err)
				handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed Get Catatan List", err.Error())
				return
			}

			surat, err := getSurat(c, ayat.Surat)
			if err != nil {
				mlog.Error(err)
				handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed Get Surat Detail", err.Error())
				return
			}

			suratPad := fmt.Sprintf("%0*d", 3, ayat.Surat)
			ayatPad := fmt.Sprintf("%0*d", 3, ayat.Nomor)

			result := format.Ayat{
				ID:        ayat.ID,
				Nomor:     ayat.Nomor,
				Surat:     ayat.Surat,
				SuratNama: surat.Nama,
				QS:        fmt.Sprintf("[%d:%d]", ayat.Surat, ayat.Nomor),
				Juz:       ayat.Juz,
				TxtAR:     ayat.TxtAR,
				TxtID:     ayat.TxtID,
				TxtIDT:    ayat.TxtIDT,
				TxtTafsir: ayat.TxtTafsir,
				Image:     ayat.Image,
				Audio:     fmt.Sprintf(audioURL, suratPad, ayatPad),
				Catatan:   catatans,
			}

			ayats = append(ayats, result)
		}

		mlog.Info("Search ayat API, search query : %s", c.Param("searchText"))

		// set the result
		result.Ayats = ayats
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

	handlers.DefaultResponse(c, http.StatusOK, "Success Get All Ayat", result)
	return
}

// GetDetailAyat : function to get detail ayat
func GetDetailAyat(c *gin.Context) {
	var ayat quran.Ayat
	detailID := c.Param("surat_ayat") // format for detail => surat_id:ayat_id => example 2:100
	exp := strings.Split(detailID, ":")
	suratID, _ := strconv.Atoi(exp[0])
	ayatID, _ := strconv.Atoi(exp[1])
	result := format.Ayat{}

	redisKey := redisKeyGeneratorAyat("", "", suratID, ayatID, 0, 0, 0, false)

	mlog.Info("Search ayat API, url : %+v", helpers.GetCurrentURL(c))

	val, err := cache.Get(c, redisKey).Result()
	if err != nil || val == "" {
		findFilter := bson.M{
			"$and": []bson.M{
				{"surat": suratID},
				{"nomor": ayatID},
			},
		}

		if err = ayatCollection.FindOne(c, findFilter).Decode(&ayat); err != nil {
			mlog.Error(err)
		}

		catatans, err := getCatatans(c, ayat.Surat, ayat.Nomor)
		if err != nil {
			mlog.Error(err)
			handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed Get Catatan List", err.Error())
			return
		}

		surat, err := getSurat(c, ayat.Surat)
		if err != nil {
			mlog.Error(err)
			handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed Get Surat Detail", err.Error())
			return
		}

		suratPad := fmt.Sprintf("%0*d", 3, ayat.Surat)
		ayatPad := fmt.Sprintf("%0*d", 3, ayat.Nomor)

		result = format.Ayat{
			ID:        ayat.ID,
			Nomor:     ayat.Nomor,
			Surat:     ayat.Surat,
			SuratNama: surat.Nama,
			QS:        fmt.Sprintf("[%d:%d]", ayat.Surat, ayat.Nomor),
			Juz:       ayat.Juz,
			TxtAR:     ayat.TxtAR,
			TxtID:     ayat.TxtID,
			TxtIDT:    ayat.TxtIDT,
			TxtTafsir: ayat.TxtTafsir,
			Image:     ayat.Image,
			Audio:     fmt.Sprintf(audioURL, suratPad, ayatPad),
			Catatan:   catatans,
		}

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
			handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed set redis", err.Error())
			return
		}
	}

	if val != "" {
		mlog.Info("GET redis key : %s", redisKey)

		err = msgpack.Unmarshal([]byte(val), &result)
		if err != nil {
			handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed JSON encode", err.Error())
			return
		}
	}

	handlers.DefaultResponse(c, http.StatusOK, "Success Get One Ayat", result)
}
