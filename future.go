package odbc

import (
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/future"
)

type Future struct {
	Quote

	UnderlyingSymbol         string  `db:"underlying_symbol" json:"underlying_symbol"`
	OpenInterest             int     `db:"open_interest" json:"open_interest"`
	ExpireDate               int     `db:"expire_date" json:"expire_date"`
	Strike                   float64 `db:"strike" json:"strike"`
	UnderlyingExchangeSymbol string  `db:"underlying_exchange_symbol" json:"underlying_exchange_symbol"`
	HeadSymbolAsString       string  `db:"head_symbol_as_string" json:"head_symbol_as_string"`
	IsContractSymbol         bool    `db:"is_contract_symbol" json:"is_contract_symbol"`
}

func NewAnonFutureFromAPI(c interface{}) (f interface{}, ok bool) {
	var apiFuture *finance.Future
	apiFuture, ok = c.(*finance.Future)
	if !ok {
		return
	}

	f, ok = NewFutureFromAPI(apiFuture)
	return
}

func GetAnonFutureFromAPI(symbol string) (interface{}, error) {
	return future.Get(symbol)
}

func NewFutureFromAPI(e *finance.Future) (f Future, ok bool) {
	ok = e != nil && (&e.Quote) != nil
	if !ok {
		return
	}

	f = Future{
		Quote: NewQuoteFromAPI(&e.Quote),

		UnderlyingSymbol:         e.UnderlyingSymbol,
		OpenInterest:             e.OpenInterest,
		ExpireDate:               e.ExpireDate,
		Strike:                   e.Strike,
		UnderlyingExchangeSymbol: e.UnderlyingExchangeSymbol,
		HeadSymbolAsString:       e.HeadSymbolAsString,
		IsContractSymbol:         e.IsContractSymbol,
	}
	return
}
