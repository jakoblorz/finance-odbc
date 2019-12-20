package odbc

import (
	"github.com/piquette/finance-go"
)

// `db:"(.*)"`

type Quote struct {
	MarketID    string `db:"market_id" json:"market_id"`
	MarketState string `db:"market_state" json:"market_state"`

	Symbol      string `db:"symbol" json:"symbol"`
	Type        string `db:"type" json:"type"`
	ShortName   string `db:"short_name" json:"short_name"`
	Currency    string `db:"currency" json:"currency"`
	IsTradeable bool   `db:"is_tradeable" json:"is_tradeable"`

	Bid     float64 `db:"bid" json:"bid"`
	BidSize int     `db:"bid_size" json:"bid_size"`
	Ask     float64 `db:"ask" json:"ask"`
	AskSize int     `db:"ask_size" json:"ask_size"`

	PreMarketPrice         float64 `db:"pre_market_price" json:"pre_market_price"`
	PreMarketChange        float64 `db:"pre_market_change" json:"pre_market_change"`
	PreMarketChangePercent float64 `db:"pre_market_change_percent" json:"pre_market_change_percent"`
	PreMarketTime          int     `db:"pre_market_time" json:"pre_market_time"`

	RegularMarketChangePercent float64 `db:"regular_market_change_percent" json:"regular_market_change_percent"`
	RegularMarketPreviousClose float64 `db:"regular_market_previous_close" json:"regular_market_previous_close"`
	RegularMarketPrice         float64 `db:"regular_market_price" json:"regular_market_price"`
	RegularMarketTime          int     `db:"regular_market_time" json:"regular_market_time"`
	RegularMarketChange        float64 `db:"regular_market_change" json:"regular_market_change"`
	RegularMarketDayHigh       float64 `db:"regular_market_day_high" json:"regular_market_day_high"`
	RegularMarketDayLow        float64 `db:"regular_market_day_low" json:"regular_market_day_low"`
	RegularMarketVolume        int     `db:"regular_market_volume" json:"regular_market_volume"`

	PostMarketPrice         float64 `db:"post_market_price" json:"post_market_price"`
	PostMarketChange        float64 `db:"post_market_change" json:"post_market_change"`
	PostMarketChangePercent float64 `db:"post_market_change_percent" json:"post_market_change_percent"`
	PostMarketTime          int     `db:"post_market_time" json:"post_market_time"`

	FiftyTwoWeekLowChange         float64 `db:"fifty_two_week_low_change" json:"fifty_two_week_low_change"`
	FiftyTwoWeekLowChangePercent  float64 `db:"fifty_two_week_low_change_percent" json:"fifty_two_week_low_change_percent"`
	FiftyTwoWeekHighChange        float64 `db:"fifty_two_week_high_change" json:"fifty_two_week_high_change"`
	FiftyTwoWeekHighChangePercent float64 `db:"fifty_two_week_high_change_percent" json:"fifty_two_week_high_change_percent"`
	FiftyTwoWeekLow               float64 `db:"fifty_two_week_low" json:"fifty_two_week_low"`
	FiftyTwoWeekHigh              float64 `db:"fifty_two_week_high" json:"fifty_two_week_high"`

	FiftyDayAverage              float64 `db:"fifty_day_average" json:"fifty_day_average"`
	FiftyDayAverageChange        float64 `db:"fifty_day_average_change" json:"fifty_day_average_change"`
	FiftyDayAverageChangePercent float64 `db:"fifty_day_average_change_percent" json:"fifty_day_average_change_percent"`

	TwoHundredDayAverage              float64 `db:"two_hundred_day_average" json:"two_hundred_day_average"`
	TwoHundredDayAverageChange        float64 `db:"two_hundred_day_average_change" json:"two_hundred_day_average_change"`
	TwoHundredDayAverageChangePercent float64 `db:"two_hundred_day_average_change_percent" json:"two_hundred_day_average_change_percent"`

	AverageDailyVolumeThreeMonth int `db:"average_daily_volume_three_month" json:"average_daily_volume_three_month"`
	AverageDailyVolumeTenDay     int `db:"average_daily_volume_ten_day" json:"average_daily_volume_ten_day"`

	SourceName     string `db:"source_name" json:"source_name"`
	SourceDelay    int    `db:"source_delay" json:"source_delay"`
	SourceInterval int    `db:"source_interval" json:"source_interval"`

	ExchangeID           string `db:"exchange_id" json:"exchange_id"`
	ExchangeName         string `db:"exchange_name" json:"exchange_name"`
	ExchangeTimezoneName string `db:"exchange_timezone_name" json:"exchange_timezone_name"`
	ExchangeTimezoneCode string `db:"exchange_timezone_code" json:"exchange_timezone_code"`

	GMTOffsetMillisecond int `db:"gmt_offset_millisecond" json:"gmt_offset_millisecond"`
}

func NewQuoteFromAPI(d *finance.Quote) Quote {
	return Quote{
		MarketID:    d.MarketID,
		MarketState: string(d.MarketState),

		Symbol:      d.Symbol,
		Type:        string(d.QuoteType),
		ShortName:   d.ShortName,
		Currency:    d.CurrencyID,
		IsTradeable: d.IsTradeable,

		Bid:     d.Bid,
		BidSize: d.BidSize,
		Ask:     d.Ask,
		AskSize: d.AskSize,

		PreMarketPrice:         d.PreMarketPrice,
		PreMarketChange:        d.PreMarketChange,
		PreMarketChangePercent: d.PreMarketChangePercent,
		PreMarketTime:          d.PreMarketTime,

		RegularMarketChangePercent: d.RegularMarketChangePercent,
		RegularMarketPreviousClose: d.RegularMarketPreviousClose,
		RegularMarketPrice:         d.RegularMarketPrice,
		RegularMarketTime:          d.RegularMarketTime,
		RegularMarketChange:        d.RegularMarketChange,
		RegularMarketDayHigh:       d.RegularMarketDayHigh,
		RegularMarketDayLow:        d.RegularMarketDayLow,
		RegularMarketVolume:        d.RegularMarketVolume,

		PostMarketPrice:         d.PostMarketPrice,
		PostMarketChange:        d.PostMarketChange,
		PostMarketChangePercent: d.PostMarketChangePercent,
		PostMarketTime:          d.PostMarketTime,

		FiftyTwoWeekLowChange:         d.FiftyTwoWeekLowChange,
		FiftyTwoWeekLowChangePercent:  d.FiftyTwoWeekLowChangePercent,
		FiftyTwoWeekHighChange:        d.FiftyTwoWeekHighChange,
		FiftyTwoWeekHighChangePercent: d.FiftyTwoWeekHighChangePercent,
		FiftyTwoWeekLow:               d.FiftyTwoWeekLow,
		FiftyTwoWeekHigh:              d.FiftyTwoWeekHigh,

		FiftyDayAverage:              d.FiftyDayAverage,
		FiftyDayAverageChange:        d.FiftyDayAverageChange,
		FiftyDayAverageChangePercent: d.FiftyDayAverageChangePercent,

		TwoHundredDayAverage:              d.TwoHundredDayAverage,
		TwoHundredDayAverageChange:        d.TwoHundredDayAverageChange,
		TwoHundredDayAverageChangePercent: d.TwoHundredDayAverageChangePercent,

		AverageDailyVolumeThreeMonth: d.AverageDailyVolume3Month,
		AverageDailyVolumeTenDay:     d.AverageDailyVolume10Day,

		SourceName:     d.QuoteSource,
		SourceDelay:    d.QuoteDelay,
		SourceInterval: d.SourceInterval,

		ExchangeID:           d.ExchangeID,
		ExchangeName:         d.FullExchangeName,
		ExchangeTimezoneName: d.ExchangeTimezoneName,
		ExchangeTimezoneCode: d.ExchangeTimezoneShortName,

		GMTOffsetMillisecond: d.GMTOffSetMilliseconds,
	}
}
