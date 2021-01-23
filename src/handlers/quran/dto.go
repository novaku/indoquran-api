package quran

// TopikTree : structure of topik tree for the front-end
type TopikTree struct {
	ID     int    `json:"id"`
	Label  string `json:"label"`
	Avatar string `json:"avatar"`
	Header string `json:"header"`
	Icon   string `json:"icon"`
	Body   string `json:"body"`
	Lazy   bool   `json:"lazy"`
}
