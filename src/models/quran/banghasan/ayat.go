package banghasan

// Ayat : structure for ayat
type Ayat struct {
	Ayat AyatObj `json:"ayat"`
}

// AyatObj : structure for ayat data
type AyatObj struct {
	Data AyatData `json:"data"`
}

// AyatData : structure for ayat object
type AyatData struct {
	AR  []AyatText `json:"ar"`
	IDT []AyatText `json:"idt"`
	ID  []AyatText `json:"id"`
}

// AyatText : structure for ayat text
type AyatText struct {
	ID    string `json:"id"`
	Surat string `json:"surat"`
	Ayat  string `json:"ayat"`
	Teks  string `json:"teks"`
}
