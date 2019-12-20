package odbc

import (
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/mutualfund"
)

type MutualFund struct {
	Quote

	YTDReturn                    float64 `db:"ytd_return" json:"ytd_return"`
	TrailingThreeMonthReturns    float64 `db:"trailing_three_month_returns" json:"trailing_three_month_returns"`
	TrailingThreeMonthNavReturns float64 `db:"trailing_three_month_nav_returns" json:"trailing_three_month_nav_returns"`
}

func NewAnonMutualFundFromAPI(c interface{}) (m interface{}, ok bool) {
	var apiMutualFund *finance.MutualFund
	apiMutualFund, ok = c.(*finance.MutualFund)
	if !ok {
		return
	}

	m = NewMutualFundFromAPI(apiMutualFund)
	return
}

func GetAnonMutualFundFromAPI(symbol string) (interface{}, error) {
	return mutualfund.Get(symbol)
}

func NewMutualFundFromAPI(e *finance.MutualFund) MutualFund {
	return MutualFund{
		Quote: NewQuoteFromAPI(&e.Quote),

		YTDReturn:                    e.YTDReturn,
		TrailingThreeMonthReturns:    e.TrailingThreeMonthReturns,
		TrailingThreeMonthNavReturns: e.TrailingThreeMonthNavReturns,
	}
}
