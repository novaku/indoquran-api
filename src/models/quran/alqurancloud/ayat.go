package alqurancloud

// Ayat : structure of ayat
type Ayat struct {
	Data struct {
		Text          string `json:"text"`
		NumberInSurah int    `json:"numberInSurah"`
		Juz           int    `json:"juz"`
		Sajda         bool   `json:"sajda"`
	} `json:"data"`
}
