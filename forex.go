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

	p, ok = NewForexFromAPI(apiForexPair)
	return
}

func GetAnonForexFromAPI(symbol string) (interface{}, error) {
	return forex.Get(symbol)
}

func NewForexFromAPI(e *finance.ForexPair) (fp ForexPair, ok bool) {
	ok = e != nil && (&e.Quote) != nil
	if !ok {
		return
	}

	fp = ForexPair{
		Quote: NewQuoteFromAPI(&e.Quote),
	}
	return
}
