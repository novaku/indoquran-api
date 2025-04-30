package model

import (
	"gorm.io/gorm"
)

// QuranAyat represents a row in the quran_ayat table
type QuranAyat struct {
	gorm.Model
	Surat      int     `json:"surat"`
	Ayat       int     `json:"ayat"`
	Text       string  `json:"text"`
	Simple     *string `json:"simple,omitempty"`
	AyatKey    string  `json:"ayat_key"`
	AyatNumber int     `json:"ayat_number"`
}

// GetAyat retrieves a specific ayat by surat and ayat number
func GetAyat(db *gorm.DB, suratID, ayatNumber int) (*QuranAyat, error) {
	var ayat QuranAyat
	result := db.Where("surat = ? AND ayat = ?", suratID, ayatNumber).First(&ayat)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ayat, nil
}

// GetAyatInSurat retrieves all ayat in a specific surat
func GetAyatInSurat(db *gorm.DB, suratID int) ([]QuranAyat, error) {
	var ayats []QuranAyat
	result := db.Where("surat = ?", suratID).Find(&ayats)
	if result.Error != nil {
		return nil, result.Error
	}
	return ayats, nil
}
