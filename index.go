package odbc

import (
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/index"
)

type Index struct {
	Quote
}

func NewAnonIndexFromAPI(c interface{}) (i interface{}, ok bool) {
	var apiIndex *finance.Index
	apiIndex, ok = c.(*finance.Index)
	if !ok {
		return
	}

	i, ok = NewIndexFromAPI(apiIndex)
	return
}

func GetAnonIndexFromAPI(symbol string) (interface{}, error) {
	return index.Get(symbol)
}

func NewIndexFromAPI(e *finance.Index) (i Index, ok bool) {
	ok = e != nil && (&e.Quote) != nil
	if !ok {
		return
	}

	i = Index{
		Quote: NewQuoteFromAPI(&e.Quote),
	}
	return
}
