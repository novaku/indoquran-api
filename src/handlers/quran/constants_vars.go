package quran

import (
	"fmt"

	"bitbucket.org/indoquran-api/src/handlers"
	"bitbucket.org/indoquran-api/src/models/quran/format"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	suratCollName     = "surat"
	ayatCollName      = "ayat"
	catatanCollName   = "catatan"
	tafsirCollName    = "tafsir"
	kataBijakCollName = "kata_bijak"
	topikCollName     = "topik"
	latinCollName     = "latin"
	imageURL          = "http://cdn.islamic.network/quran/images/%d_%d.png"
	audioURL          = "https://audio.qurancdn.com/Alafasy/mp3/%s%s.mp3"
	topikParentIcon   = "work_outline"
	topikChildIcon    = "content_copy"
)

var (
	pagination format.Pagination
	cursor     *mongo.Cursor
	err        error
)

var (
	db                  = *handlers.MongoInstance()
	cache               = *handlers.RedisInstance()
	suratCollection     = db.Collection(suratCollName)
	ayatCollection      = db.Collection(ayatCollName)
	catatanCollection   = db.Collection(catatanCollName)
	tafsirCollection    = db.Collection(tafsirCollName)
	kataBijakCollection = db.Collection(kataBijakCollName)
	topikCollection     = db.Collection(topikCollName)
	latinCollection     = db.Collection(latinCollName)
)

func redisKeyGeneratorAyat(search, sortBy string, suratID, ayatID, juz int, rowsPerPage, page int64, descending bool) string {
	return "ayat:" + search + ":" + sortBy + ":" + fmt.Sprintf("%d", suratID) + ":" + fmt.Sprintf("%d", ayatID) + ":" + fmt.Sprintf("%d", juz) + ":" + fmt.Sprintf("%d", rowsPerPage) + ":" + fmt.Sprintf("%d", page) + ":" + fmt.Sprintf("%t", descending)
}

func redisKeyGeneratorSurat(sortBy string, rowsPerPage, page int64, descending bool, id int) string {
	return "surat:" + sortBy + ":" + fmt.Sprintf("%d", rowsPerPage) + ":" + fmt.Sprintf("%d", page) + ":" + fmt.Sprintf("%t", descending) + ":" + fmt.Sprintf("%d", id)
}

func redisKeyGeneratorTopik(id string) string {
	return "topik:" + id
}
