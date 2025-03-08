package model

type QuranAyat struct {
	AyatKey    string  `gorm:"primaryKey;size:7;not null"`
	AyatNumber int     `gorm:"not null"`
	Surat      int     `gorm:"default:0;not null"`
	Ayat       int     `gorm:"default:0;not null"`
	Text       string  `gorm:"type:text;not null"`
	Simple     *string `gorm:"type:text"`     // Nullable
	Juz        *int    `gorm:"type:smallint"` // Nullable
	Hezb       *int    `gorm:"type:smallint"` // Nullable
	Page       *int    `gorm:"type:smallint"` // Nullable
	Rub        *int    `gorm:"default:NULL"`  // Nullable
}

func (QuranAyat) TableName() string {
	return "quran_ayat"
}
