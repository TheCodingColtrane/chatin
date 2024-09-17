package services

import (
	"chatin/auth"
	"chatin/data"
	searchQueries "chatin/queries/search"
)

var searchRepo = data.NewSearchData()

type searchService struct{}

func NewSearchService() *searchService {
	return &searchService{}
}

func (sv *searchService) GetSearchResultRowNumber(searchArgs searchQueries.SearchArguments) (int, error) {
	var foundRowsChan = make(chan int)
	var errChan = make(chan error)
	go searchRepo.GetSearchResultRowNumber(searchArgs, foundRowsChan, errChan)
	select {
	case rows := <-foundRowsChan:
		return rows, nil
	case err := <-errChan:
		return 0, err
	}
}

func (sv *searchService) GetSearchResult(searchArgs searchQueries.SearchArguments) ([]searchQueries.Items, error) {
	var searchResultChan searchQueries.SearchResultChannel
	searchResultChan.Err = make(chan error)
	searchResultChan.FoundItems = make(chan []searchQueries.Items)
	searchArgs.ChatID, _ = auth.DecodeUserID(searchArgs.ChatCode)
	searchArgs.MessageID, _ = auth.DecodeUserID(searchArgs.MessageCode)
	if searchArgs.ChatID == 0 {
		go searchRepo.GetSearchResultsByText(searchArgs, searchResultChan)
	} else {
		go searchRepo.GetSearchResultsByChatID(searchArgs, searchResultChan)
	}

	select {
	case foundItems := <-searchResultChan.FoundItems:
		return foundItems, nil
	case err := <-searchResultChan.Err:
		return []searchQueries.Items{}, err
	}

}
