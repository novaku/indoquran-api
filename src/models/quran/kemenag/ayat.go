package kemenag

type Ayat struct {
	Data []Data `json:"data"`
}

type Data struct {
	AyaID              int    `json:"aya_id"`
	AyaNumber          int    `json:"aya_number"`
	AyaText            string `json:"aya_text"`
	SuraID             int    `json:"sura_id"`
	JuzID              int    `json:"juz_id"`
	PageNumber         int    `json:"page_number"`
	TranslationAyaText string `json:"translation_aya_text"`
}
