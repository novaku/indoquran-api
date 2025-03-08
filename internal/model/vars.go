package model

type (
	AyatDetail struct {
		ID         int    `json:"id"`
		Juz        int    `json:"juz"`
		Surat      int    `json:"surat"`
		Ayat       int    `json:"ayat"`
		TextIndo   string `json:"text_indo"`
		TextArabic string `json:"text_arabic"`
	}

	CountResult struct {
		Type       string `gorm:"column:type" json:"type"`
		Identifier int    `gorm:"column:identifier" json:"identifier"`
		Count      int    `gorm:"column:count" json:"count"`
	}

	SearchNewResult struct {
		Text        string `json:"text"`
		Simple      string `json:"simple"`
		Surat       int    `json:"surat"`
		Ayat        int    `json:"ayat"`
		Translation string `json:"translation"`
		AyatKey     string `json:"ayat_key"`
		AyatNumber  int    `json:"ayat_number"`
	}
)
