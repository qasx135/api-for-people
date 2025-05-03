package model

import (
	"log/slog"
	"net/http"
	"strconv"
)

type UserQueryParams struct {
	Page        int
	Limit       int
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
	SortBy      string
	SortOrder   string
}

func ParseQueryParams(r *http.Request) UserQueryParams {
	slog.Debug("Fetching parameters")
	queryParams := r.URL.Query()
	page, _ := strconv.Atoi(queryParams.Get("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	age, _ := strconv.Atoi(queryParams.Get("age"))
	return UserQueryParams{
		Page:        page,
		Limit:       limit,
		Name:        queryParams.Get("name"),
		Surname:     queryParams.Get("surname"),
		Patronymic:  queryParams.Get("patronymic"),
		Age:         age,
		Gender:      queryParams.Get("gender"),
		Nationality: queryParams.Get("nationality"),
		SortBy:      queryParams.Get("sort"),
		SortOrder:   queryParams.Get("order"),
	}
}
