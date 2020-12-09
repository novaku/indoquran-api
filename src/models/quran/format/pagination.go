package format

type Pagination struct {
	SortBy      string `json:"sortBy"`     // field name
	Descending  bool   `json:"descending"` // true & false
	Page        int64  `json:"page"`
	RowsPerPage int64  `json:"rowsPerPage"`
	RowsNumber  int64  `json:"rowsNumber"`
}
