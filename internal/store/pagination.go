package store

import (
	"net/http"
	"strconv"
	"strings"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Search string   `json:"search" validate:"max=100"`
	Tags   []string `json:"tags" validate:"max=5"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

// parse the variables from the request context
func (feedQuery PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	urlQuery := r.URL.Query()

	limit := urlQuery.Get("limit")
	if limit != "" {
		limitValue, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return feedQuery, err
		}

		feedQuery.Limit = int(limitValue)
	}

	offset := urlQuery.Get("offset")
	if offset != "" {
		offsetValue, err := strconv.ParseInt(offset, 10, 64)
		if err != nil {
			return feedQuery, err
		}

		feedQuery.Offset = int(offsetValue)
	}

	sort := urlQuery.Get("sort")
	if sort != "" {
		feedQuery.Sort = sort
	}

	search := urlQuery.Get("search")
	if search != "" {
		feedQuery.Search = search
	}

	tags := urlQuery.Get("tags")
	if tags != "" {
		feedQuery.Tags = strings.Split(tags, ",")
	} else {
		feedQuery.Tags = []string{}
	}

	since := urlQuery.Get("since")
	if since != "" {
		sinceTime := parsePostgresTime(since)
		feedQuery.Since = sinceTime
	} else {
		sinceTime := parsePostgresTime(feedQuery.Since)
		feedQuery.Since = sinceTime
	}

	until := urlQuery.Get("until")
	if until != "" {
		untilTime := parsePostgresTime(until)
		feedQuery.Until = untilTime
	} else {
		untilTime := parsePostgresTime(feedQuery.Until)
		feedQuery.Until = untilTime
	}

	return feedQuery, nil
}
