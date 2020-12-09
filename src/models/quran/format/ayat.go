package format

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ayat struct {
	ID        primitive.ObjectID `json:"id"`
	Nomor     int                `json:"nomor"`
	Surat     int                `json:"surat"`
	SuratNama string             `json:"surat_nama"`
	QS        string             `json:"qs"`
	Juz       int                `json:"juz"`
	TxtAR     string             `json:"txt_ar"`
	TxtID     string             `json:"txt_id"`
	TxtIDT    string             `json:"txt_idt"`
	TxtTafsir string             `json:"txt_tafsir"`
	Image     string             `json:"image"`
	Audio     string             `json:"audio"`
	Catatan   []Catatan          `json:"catatan"`
}

type Catatan struct {
	ID        int    `json:"id"`
	TeksQuran string `json:"teks_quran"`
}

type Ayats struct {
	Ayats      []Ayat     `json:"ayats"`
	Pagination Pagination `json:"pagination"`
}
