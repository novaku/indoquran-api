package model

type IdMuntakhab struct {
	Index int    `gorm:"primaryKey;autoIncrement:false" json:"index"`
	Surat int    `gorm:"default:0" json:"surat"`
	Ayat  int    `gorm:"default:0" json:"ayat"`
	Text  string `json:"text"`
}

func (IdMuntakhab) TableName() string {
	return "id_muntakhab"
}
