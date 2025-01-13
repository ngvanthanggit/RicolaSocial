package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
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
		sinceTime := parseTime(since)
		feedQuery.Since = sinceTime
	}

	until := urlQuery.Get("until")
	if until != "" {
		untilTime := parseTime(until)
		feedQuery.Until = untilTime
	}

	return feedQuery, nil
}

func parseTime(s string) string {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}
