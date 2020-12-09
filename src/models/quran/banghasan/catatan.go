package banghasan

// CatatanRoot : structure for catatan root
type CatatanRoot struct {
	Catatan CatatanObj   `json:"catatan"`
	Quran   QuranCatatan `json:"quran"`
}

// CatatanObj : structure for catatan object
type CatatanObj struct {
	Nomor string `json:"nomor"`
	Teks  string `json:"teks"`
}

// QuranCatatan : structure for quran object
type QuranCatatan struct {
	ID    string `json:"id"`
	Surat string `json:"surat"`
	Ayat  string `json:"ayat"`
	Teks  string `json:"teks"`
}
