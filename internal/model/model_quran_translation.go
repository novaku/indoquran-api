package model

type QuranTranslation struct {
	TranslationID  uint64 `gorm:"primaryKey;autoIncrement;column:translation_id"`
	AyatKey        string `gorm:"size:7;not null;column:ayat_key"`
	TranslationKey string `gorm:"size:20;not null;column:translation_key"`
	Surat          int    `gorm:"not null;column:surat"`
	Ayat           int    `gorm:"not null;column:ayat"`
	Translation    string `gorm:"type:text;not null;column:translation"`
}

func (QuranTranslation) TableName() string {
	return "quran_translation"
}
