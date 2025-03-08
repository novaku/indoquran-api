package model

type ArabicQuran struct {
	ID    int    `gorm:"primary_key;auto_increment"`
	Surat int    `gorm:"not null"`
	Ayat  int    `gorm:"not null"`
	Text  string `gorm:"type:text;not null"`
}

func (ArabicQuran) TableName() string {
	return "arabic_quran"
}
