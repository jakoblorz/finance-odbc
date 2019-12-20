package odbc

import (
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/equity"
)

type Equity struct {
	Quote

	LongName  string `db:"long_name" json:"long_name"`
	MarketCap int64  `db:"market_cap" json:"market_cap"`

	EarningsTimestamp      int `db:"earnings_timestamp" json:"earnings_timestamp"`
	EarningsTimestampStart int `db:"earnings_timestamp_start" json:"earnings_timestamp_start"`
	EarningsTimestampEnd   int `db:"earnings_timestamp_end" json:"earnings_timestamp_end"`

	TrailingTwelveMonthsEarningsPerShare float64 `db:"trailing_twelve_months_earnings_per_share" json:"trailing_twelve_months_earnings_per_share"`
	TrailingAnnualDividendRate           float64 `db:"trailing_annual_dividend_rate" json:"trailing_annual_dividend_rate"`
	TrailingAnnualDividendYield          float64 `db:"trailing_annual_dividend_yield" json:"trailing_annual_dividend_yield"`
	TrailingPriceToEarnings              float64 `db:"trailing_price_to_earnings" json:"trailing_price_to_earnings"`

	ForwardEarningsPerShare float64 `db:"forward_earnings_per_share" json:"forward_earnings_per_share"`
	ForwardPriceToEarnings  float64 `db:"forward_price_to_earnings" json:"forward_price_to_earnings"`

	DividendRate int     `db:"dividend_rate" json:"dividend_rate"`
	BookValue    float64 `db:"book_value" json:"book_value"`
	PriceToBook  float64 `db:"price_to_book" json:"price_to_book"`

	SharesOutstanding int `db:"shares_outstanding" json:"shares_outstanding"`
}

func NewAnonEquityFromAPI(c interface{}) (e interface{}, ok bool) {
	var apiEquity *finance.Equity
	apiEquity, ok = c.(*finance.Equity)
	if !ok {
		return
	}

	e = NewEquityFromAPI(apiEquity)
	return
}

func GetAnonEquityFromAPI(symbol string) (interface{}, error) {
	return equity.Get(symbol)
}

func NewEquityFromAPI(e *finance.Equity) Equity {
	return Equity{
		Quote: NewQuoteFromAPI(&e.Quote),

		LongName:  e.LongName,
		MarketCap: e.MarketCap,

		EarningsTimestamp:      e.EarningsTimestamp,
		EarningsTimestampStart: e.EarningsTimestampStart,
		EarningsTimestampEnd:   e.EarningsTimestampEnd,

		TrailingTwelveMonthsEarningsPerShare: e.EpsTrailingTwelveMonths,
		TrailingAnnualDividendRate:           e.TrailingAnnualDividendRate,
		TrailingAnnualDividendYield:          e.TrailingAnnualDividendYield,
		TrailingPriceToEarnings:              e.TrailingPE,

		ForwardEarningsPerShare: e.EpsForward,
		ForwardPriceToEarnings:  e.ForwardPE,

		DividendRate: e.DividendDate,
		BookValue:    e.BookValue,
		PriceToBook:  e.PriceToBook,

		SharesOutstanding: e.SharesOutstanding,
	}
}
