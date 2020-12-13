package quran

import (
	"fmt"

	"bitbucket.org/indoquran-api/src/handlers"
)

const (
	suratCollName   = "surat"
	ayatCollName    = "ayat"
	catatanCollName = "catatan"
	tafsirCollName  = "tafsir"
	imageURL        = "http://cdn.islamic.network/quran/images/%d_%d.png"
	audioURL        = "https://audio.qurancdn.com/Alafasy/mp3/%s%s.mp3"
)

var (
	db                = *handlers.MongoConfig()
	cache             = *handlers.RedisConfig()
	suratCollection   = db.Collection(suratCollName)
	ayatCollection    = db.Collection(ayatCollName)
	catatanCollection = db.Collection(catatanCollName)
	tafsirCollection  = db.Collection(tafsirCollName)
)

func redisKeyGeneratorAyat(search, sortBy string, surat, juz int, rowsPerPage, page int64, descending bool) string {
	return search + ":" + sortBy + ":" + fmt.Sprintf("%d", surat) + ":" + fmt.Sprintf("%d", juz) + ":" + fmt.Sprintf("%d", rowsPerPage) + ":" + fmt.Sprintf("%d", page) + ":" + fmt.Sprintf("%t", descending)
}

func redisKeyGeneratorSurat(sortBy string, rowsPerPage, page int64, descending bool, id int) string {
	return sortBy + ":" + fmt.Sprintf("%d", rowsPerPage) + ":" + fmt.Sprintf("%d", page) + ":" + fmt.Sprintf("%t", descending) + ":" + fmt.Sprintf("%d", id)
}
