package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

// parse the variables from the request context
func (pq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	urlQuery := r.URL.Query()

	limit := urlQuery.Get("limit")
	if limit != "" {
		limitValue, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return pq, err
		}

		pq.Limit = int(limitValue)
	}

	offset := urlQuery.Get("offset")
	if offset != "" {
		offsetValue, err := strconv.ParseInt(offset, 10, 64)
		if err != nil {
			return pq, err
		}

		pq.Offset = int(offsetValue)
	}

	sort := urlQuery.Get("sort")
	if sort != "" {
		pq.Sort = sort
	}

	return pq, nil
}
