package quran

import (
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
	suratCollection   = db.Collection(suratCollName)
	ayatCollection    = db.Collection(ayatCollName)
	catatanCollection = db.Collection(catatanCollName)
	tafsirCollection  = db.Collection(tafsirCollName)
)
