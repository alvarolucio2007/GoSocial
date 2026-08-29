package store

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int        `json:"limit" validate:"gte=1,lte=20"`
	Offset int        `json:"offset" validate:"gte=0"`
	Sort   string     `json:"sort" validate:"oneof=asc desc"`
	Tags   []string   `json:"tags" validate:"max=5"`
	Search string     `json:"search" validate:"max=100"`
	Since  *time.Time `json:"since"`
	Until  *time.Time `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	if limit := qs.Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil || l < 0 {
			return fq, fmt.Errorf("invalid limit parameter: %w", err)
		}
		fq.Limit = l
	}

	if offset := qs.Get("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil || o < 0 {
			return fq, fmt.Errorf("invalid offset parameter: %w", err)
		}
		fq.Offset = o
	}

	if sort := qs.Get("sort"); sort != "" {
		sortLower := strings.ToLower(strings.TrimSpace(sort))
		if sortLower != "asc" && sortLower != "desc" {
			return fq, fmt.Errorf("invalid sort parameter: must be ASC or DESC")
		}
		fq.Sort = sortLower
	} else {
		fq.Sort = "desc"
	}

	if tags := qs.Get("tags"); tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}
	if search := qs.Get("search"); search != "" {
		fq.Search = strings.TrimSpace(search)
	}

	if since := qs.Get("since"); since != "" {
		if t, err := parseTimeString(since); err == nil {
			fq.Since = &t
		}
	}
	if until := qs.Get("until"); until != "" {
		if t, err := parseTimeString(until); err == nil {
			fq.Until = &t
		}
	}

	return fq, nil
}

func parseTimeString(s string) (time.Time, error) {
	// Tenta formato RFC3339 primeiro (padrão de APIs)
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	// Fallback para apenas data (YYYY-MM-DD)
	return time.Parse("2006-01-02", s)
}
