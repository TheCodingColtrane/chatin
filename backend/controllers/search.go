package controllers

import (
	"chatin/config/keys"
	searchQueries "chatin/queries/search"
	"chatin/server/request"
	"chatin/server/response"
	"chatin/services"
	"net/http"
)

type searchControlller struct{}

var searchService = services.NewSearchService()

func NewSearchController() *searchControlller {
	return &searchControlller{}
}

func (sc *searchControlller) GetSearchFoundItemsNumber(res http.ResponseWriter, req *http.Request) {
	var searchArgs searchQueries.SearchArguments
	var err = request.GetBody(req, &searchArgs)
	if err != nil {
		http.Error(res, "Something went wrong", 500)
	}
	if searchArgs.Term == "" && searchArgs.MessageCode == "" {
		http.Error(res, "one or more arguments were not provided", 400)
	}

	var id = req.Context().Value(keys.UserContextKey).(float64)
	searchArgs.UserID = id
	itemsCount, err := searchService.GetSearchResultRowNumber(searchArgs)
	if err != nil {
		http.Error(res, "Something went wrong", 500)
	}

	var items struct {
		Count int `json:"itemsCount"`
	}
	items.Count = itemsCount
	response.JSON(res, items)

}

func (sc *searchControlller) GetSearchFoundItems(res http.ResponseWriter, req *http.Request) {
	var searchArgs searchQueries.SearchArguments
	var err = request.GetBody(req, &searchArgs)
	if err != nil {
		http.Error(res, "Something went wrong", 500)
	}
	if searchArgs.Term == "" && searchArgs.MessageCode == "" {
		http.Error(res, "one or more arguments were not provided", 400)
	}
	searchArgs.UserID = req.Context().Value(keys.UserContextKey).(float64)
	items, err := searchService.GetSearchResult(searchArgs)
	if err != nil {
		http.Error(res, "Something went wrong", 500)
	}
	response.JSON(res, items)
}
