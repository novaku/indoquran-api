package format

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Surat struct {
	ID         primitive.ObjectID `json:"id"`
	Nomor      int                `json:"nomor"`
	Nama       string             `json:"nama"`
	Asma       string             `json:"asma"`
	JumlahAyat int                `json:"jumlah_ayat"`
	Arti       string             `json:"arti"`
	Keterangan string             `json:"keterangan"`
}

type Surats struct {
	Surats     []Surat    `json:"surats"`
	Pagination Pagination `json:"pagination"`
}
