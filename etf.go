package odbc

import (
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/etf"
)

type ETF struct {
	Quote

	YTDReturn                    float64 `db:"ytd_return" json:"ytd_return"`
	TrailingThreeMonthReturns    float64 `db:"trailing_three_month_returns" json:"trailing_three_month_returns"`
	TrailingThreeMonthNavReturns float64 `db:"trailing_three_month_nav_returns" json:"trailing_three_month_nav_returns"`
}

func NewAnonETFFromAPI(c interface{}) (e interface{}, ok bool) {
	var apiETF *finance.ETF
	apiETF, ok = c.(*finance.ETF)
	if !ok {
		return
	}

	e, ok = NewETFFromAPI(apiETF)
	return
}

func GetAnonETFFromAPI(symbol string) (interface{}, error) {
	return etf.Get(symbol)
}

func NewETFFromAPI(e *finance.ETF) (etf ETF, ok bool) {
	ok = e != nil && (&e.Quote) != nil
	if !ok {
		return
	}

	etf = ETF{
		Quote: NewQuoteFromAPI(&e.Quote),

		YTDReturn:                    e.YTDReturn,
		TrailingThreeMonthReturns:    e.TrailingThreeMonthReturns,
		TrailingThreeMonthNavReturns: e.TrailingThreeMonthNavReturns,
	}
	return
}
