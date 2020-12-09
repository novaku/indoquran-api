package banghasan

// SuratHasil : structure for hasil
type SuratHasil struct {
	Hasil []SuratHasilArray `json:"hasil"`
}

// SuratHasilArray : structure for hasil array
type SuratHasilArray struct {
	Arti       string `json:"arti"`
	Asma       string `json:"asma"`
	Ayat       string `json:"ayat"`
	Keterangan string `json:"keterangan"`
	Nama       string `json:"nama"`
	Nomor      string `json:"nomor"`
}
