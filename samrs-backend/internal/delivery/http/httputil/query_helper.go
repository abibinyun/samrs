package httputil

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PaginationParams struct {
	Page    int
	PerPage int
	SortBy  string
	SortDir string
}

type DateRange struct {
	From *time.Time
	To   *time.Time
}

type PaginationMeta struct {
	Total   int64 `json:"total"`
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
}

func ParsePagination(c *gin.Context) PaginationParams {
	page := ParseIntWithDefault(c.Query("page"), 1)
	perPage := ParseIntWithDefault(c.Query("per_page"), 20)

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 1
	}
	if perPage > 100 {
		perPage = 100
	}

	sortBy := strings.TrimSpace(c.Query("sort_by"))
	sortDir := strings.ToLower(strings.TrimSpace(c.Query("sort_dir")))
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	return PaginationParams{
		Page:    page,
		PerPage: perPage,
		SortBy:  sortBy,
		SortDir: sortDir,
	}
}

func ParseDateRange(c *gin.Context, fromKey, toKey string) (DateRange, error) {
	var out DateRange

	if from := strings.TrimSpace(c.Query(fromKey)); from != "" {
		parsed, err := time.ParseInLocation("2006-01-02", from, time.Local)
		if err != nil {
			return out, err
		}
		start := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.Local)
		out.From = &start
	}

	if to := strings.TrimSpace(c.Query(toKey)); to != "" {
		parsed, err := time.ParseInLocation("2006-01-02", to, time.Local)
		if err != nil {
			return out, err
		}
		end := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 0, time.Local)
		out.To = &end
	}

	return out, nil
}

func ParseIntWithDefault(value string, def int) int {
	if strings.TrimSpace(value) == "" {
		return def
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return parsed
}

func ParseUint(value string) (uint64, error) {
	return strconv.ParseUint(strings.TrimSpace(value), 10, 32)
}
