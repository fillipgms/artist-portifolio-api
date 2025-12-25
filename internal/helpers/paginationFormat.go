package helpers

import "math"

type Pagination struct {
	CurrentPage int     `json:"currentPage"`
	TotalPages  int     `json:"totalPages"`
	Count       int     `json:"count"`
	Data        any   `json:"data"`
	NextPage    *string `json:"nextPage"`
	PrevPage    *string `json:"prevPage"`
	FirstPage	string	`json:"firstPage"`
	LastPage	string	`json:"lastPage"`
}

func PaginationFormat(count int64, data any, limit int64, offset int64, page int64)  Pagination {
	totalFloat := float64(count) / float64(limit)
	total := int(math.Ceil(totalFloat))

	var nextPg *string
	var prevPg *string

	if total > int(page) {
		url := ""
		nextPg = &url
	}

	if total < int(page) {
		url :=""
		prevPg = &url
	}


	return Pagination{
		CurrentPage: int(page),
		TotalPages:  int(total),
		Count:       int(count),
		Data:        data,
		NextPage:    nextPg,
		PrevPage:    prevPg,
		FirstPage:   "",
		LastPage:    "",
	}
}