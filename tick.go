package odbc

import (
	"time"

	yfin "github.com/piquette/finance-go"
	"github.com/shopspring/decimal"
)

type MetaTick struct {
	yfin.ChartMeta
	yfin.ChartBar
}

type Tick struct {
	DBEntry

	Open      decimal.Decimal `db:"open" json:"open"`
	Low       decimal.Decimal `db:"low" json:"low"`
	High      decimal.Decimal `db:"high" json:"high"`
	Close     decimal.Decimal `db:"close" json:"close"`
	AdjClose  decimal.Decimal `db:"adj_close" json:"adj_close"`
	PrvClose  decimal.Decimal `db:"prv_close" json:"prv_close"`
	Volume    int             `db:"volume" json:"volume"`
	Timestamp int             `db:"timestamp" json:"timestamp"`
	Currency  string          `db:"currency" json:"currency"`
	Symbol    string          `db:"symbol" json:"symbol"`
	Type      string          `db:"type" json:"type"`

	FirstTradeDate   int    `db:"first_trade_day" json:"first_trade_day"`
	GMTOffset        int    `db:"gmt_offset" json:"gmt_offset"`
	Timezone         string `db:"timezone" json:"timezone"`
	ExchangeName     string `db:"exchange_name" json:"exchange_name"`
	ExchangeTimezone string `db:"exchange_timezone" json:"exchange_timezone"`

	PreTimezone  string `db:"pre_timezone" json:"pre_timezone"`
	PreStart     int    `db:"pre_start" json:"pre_start"`
	PreEnd       int    `db:"pre_end" json:"pre_end"`
	PreGMTOffset int    `db:"pre_gmt_offset" json:"pre_gmt_offset"`

	RegularTimezone  string `db:"regular_timezone" json:"regular_timezone"`
	RegularStart     int    `db:"regular_start" json:"regular_start"`
	RegularEnd       int    `db:"regular_end" json:"regular_end"`
	RegularGMTOffset int    `db:"regular_gmt_offset" json:"regular_gmt_offset"`

	PostTimezone  string `db:"post_timezone" json:"post_timezone"`
	PostStart     int    `db:"post_start" json:"post_start"`
	PostEnd       int    `db:"post_end" json:"post_end"`
	PostGMTOffset int    `db:"post_gmt_offset" json:"post_gmt_offset"`

	Granularity string   `db:"granularity" json:"granularity"`
	Ranges      []string `db:"-" json:"-"`
}

func NewTickFromAPI(x *MetaTick) Tick {
	t, m := x.ChartBar, x.ChartMeta
	return Tick{
		DBEntry: DBEntry{
			InsertedAt: time.Now().UTC(),
		},

		Open:      t.Open,
		Low:       t.Low,
		High:      t.High,
		Close:     t.Close,
		AdjClose:  t.AdjClose,
		PrvClose:  decimal.NewFromFloat(m.ChartPreviousClose),
		Volume:    t.Volume,
		Timestamp: t.Timestamp,
		Currency:  m.Currency,
		Symbol:    m.Symbol,
		Type:      string(m.QuoteType),

		FirstTradeDate:   m.FirstTradeDate,
		GMTOffset:        m.Gmtoffset,
		Timezone:         m.Timezone,
		ExchangeName:     m.ExchangeName,
		ExchangeTimezone: m.ExchangeTimezoneName,

		PreTimezone:  m.CurrentTradingPeriod.Pre.Timezone,
		PreStart:     m.CurrentTradingPeriod.Pre.Start,
		PreEnd:       m.CurrentTradingPeriod.Pre.End,
		PreGMTOffset: m.CurrentTradingPeriod.Pre.Gmtoffset,

		RegularTimezone:  m.CurrentTradingPeriod.Regular.Timezone,
		RegularStart:     m.CurrentTradingPeriod.Regular.Start,
		RegularEnd:       m.CurrentTradingPeriod.Regular.End,
		RegularGMTOffset: m.CurrentTradingPeriod.Regular.Gmtoffset,

		PostTimezone:  m.CurrentTradingPeriod.Post.Timezone,
		PostStart:     m.CurrentTradingPeriod.Post.Start,
		PostEnd:       m.CurrentTradingPeriod.Post.End,
		PostGMTOffset: m.CurrentTradingPeriod.Post.Gmtoffset,

		Granularity: m.DataGranularity,
		Ranges:      m.ValidRanges,
	}
}
