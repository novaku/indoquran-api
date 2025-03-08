package controllers

import "indoquran-api/internal/model"

type (
	ResultJsonFormat struct {
		Aggregate  []*model.CountResult `json:"aggregate"`
		Pagination *Pagination          `json:"pagination"`
		Results    []*model.AyatDetail  `json:"results"`
	}
	Pagination struct {
		CurrentPage int `json:"current_page"`
		RowsPerPage int `json:"rows_per_page"`
		TotalPages  int `json:"total_pages"`
		TotalRows   int `json:"total_rows"`
	}
)
