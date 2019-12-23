package odbc

import (
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/option"
)

type Option struct {
	Quote

	UnderlyingSymbol         string `db:"underlying_symbol" json:"underlying_symbol"`
	UnderlyingExchangeSymbol string `db:"underlying_exchange_symbol" json:"underlying_exchange_symbol"`

	OpenInterest int     `db:"open_interest" json:"open_interest"`
	ExpireDate   int     `db:"expire_date" json:"expire_date"`
	Strike       float64 `db:"strike" json:"strike"`
}

func NewAnonOptionFromAPI(c interface{}) (o interface{}, ok bool) {
	var apiOption *finance.Option
	apiOption, ok = c.(*finance.Option)
	if !ok {
		return
	}

	o, ok = NewOptionFromAPI(apiOption)
	return
}

func GetAnonOptionFromAPI(symbol string) (interface{}, error) {
	return option.Get(symbol)
}

func NewOptionFromAPI(e *finance.Option) (o Option, ok bool) {
	ok = e != nil && (&e.Quote) != nil
	if !ok {
		return
	}

	o = Option{
		Quote: NewQuoteFromAPI(&e.Quote),

		UnderlyingSymbol:         e.UnderlyingSymbol,
		UnderlyingExchangeSymbol: e.UnderlyingExchangeSymbol,

		OpenInterest: e.OpenInterest,
		ExpireDate:   e.ExpireDate,
		Strike:       e.Strike,
	}
	return
}
