package model

type IdIndonesian struct {
	ID    int    `gorm:"primary_key"`
	Surat int    `gorm:"not null;default:0"`
	Ayat  int    `gorm:"not null;default:0"`
	Text  string `gorm:"not null"`
}

func (IdIndonesian) TableName() string {
	return "id_indonesian"
}
