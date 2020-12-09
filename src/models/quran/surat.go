package quran

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Surat : struct of surat
type Surat struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Nomor      int                `bson:"nomor" json:"nomor"`
	Nama       string             `bson:"nama" json:"nama"`
	Asma       string             `bson:"asma" json:"asma"`
	JumlahAyat int                `bson:"jumlah_ayat" json:"jumlah_ayat"`
	Arti       string             `bson:"arti" json:"arti"`
	Keterangan string             `bson:"keterangan" json:"keterangan"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

//Surats : collections of surat
type Surats []Surat
