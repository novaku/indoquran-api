package quran

import (
	"net/http"

	"bitbucket.org/indoquran-api/src/handlers"
	"bitbucket.org/indoquran-api/src/models/quran"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
	"go.mongodb.org/mongo-driver/bson"
)

// GetKataBijak : get one random kata bijak
func GetKataBijak(c *gin.Context) {
	kata := quran.KataBijak{}
	pipeline := []bson.M{
		{"$sample": bson.M{"size": 1}},
	}

	cur, err := kataBijakCollection.Aggregate(c, pipeline)
	if err != nil {
		mlog.Error(err)
	}
	defer cur.Close(c)

	for cur.Next(c) {
		if err := cur.Decode(&kata); err != nil {
			mlog.Error(err)
		}
	}

	handlers.DefaultResponse(c, http.StatusOK, "Success Get Kata Bijak", kata)
	return
}
