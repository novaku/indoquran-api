package quran

import (
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/indoquran-api/src/config"
	"bitbucket.org/indoquran-api/src/handlers"
	"bitbucket.org/indoquran-api/src/models/quran"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
	"github.com/vmihailenco/msgpack"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetTopik : function to get topik
func GetTopik(c *gin.Context) {
	id := c.Param("id")
	topik := quran.Topik{}
	trees := []TopikTree{}
	val, err := cache.Get(c, redisKeyGeneratorTopik(id)).Result()
	if val != "" {
		mlog.Info("GET redis key : %s", redisKeyGeneratorTopik(id))

		err = msgpack.Unmarshal([]byte(val), &trees)
		if err != nil {
			panic(err)
		}
	} else {
		findoptions := options.Find()
		findoptions.SetSort(bson.D{{Key: "topicID", Value: 1}})
		findFilter := bson.M{}

		if id == "0" {
			findFilter = bson.M{"isTitle": 1}
		} else {
			parentID, _ := strconv.Atoi(id)
			findFilter = bson.M{"parentID": parentID}
		}

		cursor, err = topikCollection.Find(c, findFilter, findoptions)

		defer cursor.Close(c)

		if err := cursor.Err(); err != nil {
			mlog.Error(err)
			handlers.DefaultResponse(c, http.StatusBadRequest, "Error Get Topik Data", err)
			return
		}

		for cursor.Next(c) {
			if err = cursor.Decode(&topik); err != nil {
				mlog.Error(err)
			}

			icon := topikChildIcon
			childs := isHasChild(c, topik.TopicID)
			if childs {
				icon = topikParentIcon
			}

			tree := TopikTree{
				ID:     topik.TopicID,
				Label:  topik.Text,
				Icon:   icon,
				Header: "generic",
				Body:   topik.Isi,
				Lazy:   childs,
			}

			trees = append(trees, tree)
		}

		// redis set
		b, err := msgpack.Marshal(&trees)
		if err != nil {
			mlog.Error(err)
			handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed JSON encode", err)
			return
		}

		ttl := time.Duration(config.Config.Cache.TTL) * time.Hour

		set, err := cache.SetNX(c, redisKeyGeneratorTopik(id), string(b), ttl).Result()
		if !set || err != nil {
			mlog.Error(err)
			handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed Set Redis Key", err)
			return
		}
	}

	handlers.DefaultResponse(c, http.StatusOK, "Success Get Surat Data", trees)
}

func isHasChild(c *gin.Context, topikID int) bool {
	findFilter := bson.M{"parentID": topikID}
	rowsNumber, err := topikCollection.CountDocuments(c, findFilter)
	if err != nil {
		mlog.Error(err)
		return false
	}

	if rowsNumber > 0 {
		return true
	}

	return false
}
