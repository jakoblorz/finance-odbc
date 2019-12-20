package odbc

import (
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/forex"
)

type ForexPair struct {
	Quote
}

func NewAnonForexFromAPI(c interface{}) (p interface{}, ok bool) {
	var apiForexPair *finance.ForexPair
	apiForexPair, ok = c.(*finance.ForexPair)
	if !ok {
		return
	}

	p = NewForexFromAPI(apiForexPair)
	return
}

func GetAnonForexFromAPI(symbol string) (interface{}, error) {
	return forex.Get(symbol)
}

func NewForexFromAPI(e *finance.ForexPair) ForexPair {
	return ForexPair{
		Quote: NewQuoteFromAPI(&e.Quote),
	}
}
