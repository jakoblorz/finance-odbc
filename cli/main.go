package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/piquette/finance-go"

	"github.com/briandowns/spinner"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/quote"

	"github.com/jakoblorz/dynsql/lib/go-sqlite3"
	fodbc "github.com/jakoblorz/finance-odbc"
	"github.com/jmoiron/sqlx"
)

var (
	DEBUG = false
)

var (
	cryptoCurrencyQuoteType = strings.ToLower(string(finance.QuoteTypeCryptoPair))
	equityQuoteType         = strings.ToLower(string(finance.QuoteTypeEquity))
	etfQuoteType            = strings.ToLower(string(finance.QuoteTypeETF))
	forexQuoteType          = strings.ToLower(string(finance.QuoteTypeForexPair))
	futureQuoteType         = strings.ToLower(string(finance.QuoteTypeFuture))
	indexQuoteType          = strings.ToLower(string(finance.QuoteTypeIndex))
	mutualfundQuoteType     = strings.ToLower(string(finance.QuoteTypeMutualFund))
	optionQuoteType         = strings.ToLower(string(finance.QuoteTypeOption))

	quoteTypePadding = map[string]string{
		cryptoCurrencyQuoteType: "",
		equityQuoteType:         "        ",
		etfQuoteType:            "           ",
		forexQuoteType:          "         ",
		futureQuoteType:         "        ",
		indexQuoteType:          "         ",
		mutualfundQuoteType:     "    ",
		optionQuoteType:         "        ",
	}

	quoteTypeTableNameMapping = map[string]string{
		cryptoCurrencyQuoteType: cryptoCurrencyQuoteType,
		equityQuoteType:         equityQuoteType,
		etfQuoteType:            etfQuoteType,
		forexQuoteType:          forexQuoteType,
		futureQuoteType:         futureQuoteType,
		indexQuoteType:          "indices",
		mutualfundQuoteType:     mutualfundQuoteType,
		optionQuoteType:         optionQuoteType,
	}

	quoteTypeControlFuncKeyMapping = map[string]string{
		cryptoCurrencyQuoteType: cryptoCurrencyQuoteType,
		equityQuoteType:         equityQuoteType,
		etfQuoteType:            etfQuoteType,
		forexQuoteType:          forexQuoteType,
		futureQuoteType:         futureQuoteType,
		indexQuoteType:          indexQuoteType,
		mutualfundQuoteType:     mutualfundQuoteType,
		optionQuoteType:         optionQuoteType,
		"ecnquote":              equityQuoteType,
	}

	quoteTypeControlFuncMapping = map[string][]interface{}{
		cryptoCurrencyQuoteType: []interface{}{
			fodbc.NewAnonCryptoFromAPI,
			fodbc.GetAnonCryptoFromAPI,
		},
		equityQuoteType: []interface{}{
			fodbc.NewAnonEquityFromAPI,
			fodbc.GetAnonEquityFromAPI,
		},
		etfQuoteType: []interface{}{
			fodbc.NewAnonETFFromAPI,
			fodbc.GetAnonETFFromAPI,
		},
		forexQuoteType: []interface{}{
			fodbc.NewAnonForexFromAPI,
			fodbc.GetAnonForexFromAPI,
		},
		futureQuoteType: []interface{}{
			fodbc.NewAnonFutureFromAPI,
			fodbc.GetAnonFutureFromAPI,
		},
		indexQuoteType: []interface{}{
			fodbc.NewAnonIndexFromAPI,
			fodbc.GetAnonIndexFromAPI,
		},
		mutualfundQuoteType: []interface{}{
			fodbc.NewAnonMutualFundFromAPI,
			fodbc.GetAnonMutualFundFromAPI,
		},
		optionQuoteType: []interface{}{
			fodbc.NewAnonOptionFromAPI,
			fodbc.GetAnonOptionFromAPI,
		},
	}

	tickTableName = "ticks"

	tickIntervalPadding = map[string]string{
		string(datetime.OneMin):      pad("", 1),
		string(datetime.TwoMins):     pad("", 1),
		string(datetime.FiveMins):    pad("", 1),
		string(datetime.FifteenMins): pad("", 0),
		string(datetime.ThirtyMins):  pad("", 0),
		string(datetime.SixtyMins):   pad("", 0),
		string(datetime.NinetyMins):  pad("", 0),
		string(datetime.OneHour):     pad("", 1),
		string(datetime.OneDay):      pad("", 1),
		string(datetime.FiveDay):     pad("", 1),
		string(datetime.OneMonth):    pad("", 0),
		string(datetime.ThreeMonth):  pad("", 0),
		string(datetime.SixMonth):    pad("", 0),
		string(datetime.OneYear):     pad("", 1),
		string(datetime.TwoYear):     pad("", 1),
		string(datetime.FiveYear):    pad("", 1),
		string(datetime.TenYear):     pad("", 0),
		string(datetime.YTD):         pad("", 0),
		string(datetime.Max):         pad("", 0),
	}
)

var (
	cryptoFlags     = []interface{}{cryptoCurrencyQuoteType, flag.String(cryptoCurrencyQuoteType, "", "Download Crypto Information")}
	equityFlags     = []interface{}{equityQuoteType, flag.String(equityQuoteType, "", "Download Equity Information")}
	etfFlags        = []interface{}{etfQuoteType, flag.String(etfQuoteType, "", "Download ETF Information")}
	forexFlags      = []interface{}{forexQuoteType, flag.String(forexQuoteType, "", "Download Forex Information")}
	futureFlags     = []interface{}{futureQuoteType, flag.String(futureQuoteType, "", "Download Future Information")}
	stockindexFlags = []interface{}{indexQuoteType, flag.String(indexQuoteType, "", "Download Index Information")}
	mutualfundFlags = []interface{}{mutualfundQuoteType, flag.String(mutualfundQuoteType, "", "Download Mututal Fund Information")}
	optionFlags     = []interface{}{optionQuoteType, flag.String(optionQuoteType, "", "Donwload Option Information")}

	allMetaInformationFlags = [][]interface{}{
		cryptoFlags,
		equityFlags,
		etfFlags,
		forexFlags,
		futureFlags,
		stockindexFlags,
		mutualfundFlags,
		optionFlags,
	}

	tickFlag        = flag.String(tickTableName, "", "Download Pricing Information")
	useAllTicksFlag = flag.Bool("all", false, "Download all possible time intervals")

	oneMinTickIntervalFlags     = []interface{}{string(datetime.OneMin), flag.Bool(string(datetime.OneMin), false, "Use a tick interval of 1min")}
	twoMinTickIntervalFlags     = []interface{}{string(datetime.TwoMins), flag.Bool(string(datetime.TwoMins), false, "Use a tick interval of 2min")}
	fiveMinTickIntervalFlags    = []interface{}{string(datetime.FiveMins), flag.Bool(string(datetime.FiveMins), false, "Use a tick interval of 5min")}
	fifteenMinTickIntervalFlags = []interface{}{string(datetime.FifteenMins), flag.Bool(string(datetime.FifteenMins), false, "Use a tick interval of 15min")}
	thirtyMinTickIntervalFlags  = []interface{}{string(datetime.ThirtyMins), flag.Bool(string(datetime.ThirtyMins), false, "Use a tick interval of 30min")}
	sixtyMinTickIntervalFlags   = []interface{}{string(datetime.SixtyMins), flag.Bool(string(datetime.SixtyMins), false, "Use a tick interval of 60min")}
	ninetyMinTickIntervalFlags  = []interface{}{string(datetime.NinetyMins), flag.Bool(string(datetime.NinetyMins), false, "Use a tick interval of 90min")}
	oneHourTickIntervalFlags    = []interface{}{string(datetime.OneHour), flag.Bool(string(datetime.OneHour), false, "Use a tick interval of 1h")}
	oneDayTickIntervalFlags     = []interface{}{string(datetime.OneDay), flag.Bool(string(datetime.OneDay), false, "Use a tick interval of 1d")}
	fiveDayTickIntervalFlags    = []interface{}{string(datetime.FiveDay), flag.Bool(string(datetime.FiveDay), false, "Use a tick interval of 5d")}
	oneMonthTickIntervalFlags   = []interface{}{string(datetime.OneMonth), flag.Bool(string(datetime.OneMonth), false, "Use a tick interval of 1mo")}
	threeMonthTickIntervalFlags = []interface{}{string(datetime.ThreeMonth), flag.Bool(string(datetime.ThreeMonth), false, "Use a tick interval of 3mo")}
	sixMonthTickIntervalFlags   = []interface{}{string(datetime.SixMonth), flag.Bool(string(datetime.SixMonth), false, "Use a tick interval of 6mo")}
	oneYearTickIntervalFlags    = []interface{}{string(datetime.OneYear), flag.Bool(string(datetime.OneYear), false, "Use a tick interval of 1y")}
	twoYearTickIntervalFlags    = []interface{}{string(datetime.TwoYear), flag.Bool(string(datetime.TwoYear), false, "Use a tick interval of 2y")}
	fiveYearTickIntervalFlags   = []interface{}{string(datetime.FiveYear), flag.Bool(string(datetime.FiveYear), false, "Use a tick interval of 5y")}
	tenYearTickIntervalFlags    = []interface{}{string(datetime.TenYear), flag.Bool(string(datetime.TenYear), false, "Use a tick interval of 10y")}
	ytdTickIntervalFlags        = []interface{}{string(datetime.YTD), flag.Bool(string(datetime.YTD), false, "Use a tick interval of YTD (Year-To-Date)")}
	maxTickIntervalFlags        = []interface{}{string(datetime.Max), flag.Bool(string(datetime.Max), false, "Use the maximum tick interval")}

	allPricingInformationFlags = [][]interface{}{
		oneMinTickIntervalFlags,
		twoMinTickIntervalFlags,
		fiveMinTickIntervalFlags,
		fifteenMinTickIntervalFlags,
		thirtyMinTickIntervalFlags,
		sixtyMinTickIntervalFlags,
		ninetyMinTickIntervalFlags,
		oneHourTickIntervalFlags,
		oneDayTickIntervalFlags,
		fiveDayTickIntervalFlags,
		oneMonthTickIntervalFlags,
		threeMonthTickIntervalFlags,
		sixMonthTickIntervalFlags,
		oneYearTickIntervalFlags,
		twoYearTickIntervalFlags,
		fiveYearTickIntervalFlags,
		tenYearTickIntervalFlags,
		ytdTickIntervalFlags,
		maxTickIntervalFlags,
	}
)

func parseFlagSetS(fs []interface{}) (tableName string, value string, ok bool) {
	tableName, ok = fs[0].(string)
	if !ok {
		return
	}

	var valuePtr *string
	valuePtr, ok = fs[1].(*string)
	if !ok {
		return
	}
	value = *valuePtr
	return
}

func parseFlagSetB(fs []interface{}) (tableName string, value bool, ok bool) {
	tableName, ok = fs[0].(string)
	if !ok {
		return
	}

	var valuePtr *bool
	valuePtr, ok = fs[1].(*bool)
	if !ok {
		return
	}
	value = *valuePtr
	return
}

func spin(prefix, final string) func(error) {
	sig := make(chan int, 0)
	ctx, cancel := context.WithCancel(context.Background())
	if final == "" {
		final = "DONE ✅\n"
	}

	if !DEBUG {
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Prefix = prefix
		s.Start()

		go func() {
			<-ctx.Done()
			s.Stop()
			close(sig)
		}()
	} else {
		print(prefix)
	}

	return func(err error) {
		cancel()
		if DEBUG {
			print(prefix)
		}
		if err != nil {
			print(fmt.Sprintf("FAILED ❌\n%s\n", err))
		} else {
			print(final)
		}
		if !DEBUG {
			<-sig
		}
	}
}

var (
	warnings = []string{}
)

func warn(msg string) {
	warnings = append(warnings, msg)
	if DEBUG {
		log.Printf("msg")
	}
}

func pad(s string, d int) string {
	for i := 0; i < d; i++ {
		s = fmt.Sprintf("%s ", s)
	}
	return s
}

func fatal(err error) {
	print(err)
	os.Exit(1)
}

func main() {
	sqlite3.DefaultTypeMap = map[string]string{
		"inserted_at": "DATETIME",
		"timestamp":   "INTEGER",
	}

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	go func() {
		db, err := sqlx.Connect("dyn-sqlite3", "finance.sqlite3")
		if err != nil {
			fatal(err)
		}
		defer db.Close()

		didDownloadMetaInformation := false
	ITERATE_METAINFORMATION_SOURCES:
		for _, f := range allMetaInformationFlags {
			quoteType, value, ok := parseFlagSetS(f)
			if !ok || value == "" {
				continue
			}

			tableName := quoteTypeTableNameMapping[quoteType]
			values := strings.Split(value, ",")
			cancel := spin(
				fmt.Sprintf("Downloading Metadata for %s:%s %d Download(s) required ", quoteType, quoteTypePadding[quoteType], len(values)),
				"",
			)

			for _, value := range values {

				q, err := quote.Get(value)
				if err != nil {
					cancel(err)
					continue ITERATE_METAINFORMATION_SOURCES
				}
				if q == nil {
					warnings = append(warnings, fmt.Sprintf("Parsing of response failed, skipping %s", value))
					continue
				}

				actualQuoteType := strings.ToLower(string(q.QuoteType))
				actualTableName, ok := quoteTypeTableNameMapping[actualQuoteType]
				if !ok {
					warnings = append(warnings, fmt.Sprintf("Found unregistered quote type %s, will create rogue table to accomodate symbol %s", actualQuoteType, value))
					actualTableName = actualQuoteType
				}
				controlFuncKey, ok := quoteTypeControlFuncKeyMapping[actualQuoteType]
				if !ok {
					warnings = append(warnings, fmt.Sprintf("Could not find parsing methods for quote type %s, skipping %s", actualQuoteType, value))
				}

				controlFuncs := quoteTypeControlFuncMapping[controlFuncKey]
				assetConvertAPI := controlFuncs[0].(func(interface{}) (interface{}, bool))
				assetObtainAPI := controlFuncs[1].(func(string) (interface{}, error))

				retrieved, err := assetObtainAPI(value)
				if err != nil {
					cancel(err)
					continue ITERATE_METAINFORMATION_SOURCES
				}
				asset, ok := assetConvertAPI(retrieved)
				if !ok {
					warnings = append(warnings, fmt.Sprintf("Parsing of response failed, skipping %s", value))
					continue
				}

				if actualTableName != tableName {
					warnings = append(warnings, fmt.Sprintf("Writing %s into table %s instead of %s", value, actualTableName, tableName))
				}

				data, err := json.Marshal(asset)
				if err != nil {
					cancel(err)
					continue ITERATE_METAINFORMATION_SOURCES
				}

				_, err = db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s JSON %s;", actualTableName, string(data)))
				if err != nil {
					cancel(err)
					continue ITERATE_METAINFORMATION_SOURCES
				}
			}
			cancel(nil)

			didDownloadMetaInformation = true
		}

		didDownloadPricingInformation := false
		if *tickFlag != "" {
			values := strings.Split(*tickFlag, ",")
			tUTCNow := time.Now().UTC()
			duplicates := map[string]int{}

		ITERATE_PRICING_INTERVALS:
			for _, f := range allPricingInformationFlags {
				interval, doDownload, ok := parseFlagSetB(f)
				if !(*useAllTicksFlag) && (!ok || !doDownload) {
					continue
				}

				cancel := spin(
					fmt.Sprintf("Downloading Historical Prices with an interval of %s:%s %d Download(s) required ", interval, tickIntervalPadding[interval], len(values)),
					"",
				)
				for _, value := range values {

					iter := chart.Get(&chart.Params{
						Symbol:   value,
						End:      datetime.New(&tUTCNow),
						Start:    datetime.FromUnix(int(tUTCNow.Unix() - 100*60*60)),
						Interval: datetime.Interval(interval),
					})

					ticks := []string{}
					for iter.Next() {
						tick := fodbc.NewTickFromAPI(&fodbc.MetaTick{
							ChartBar:  *iter.Bar(),
							ChartMeta: iter.Meta(),
						})

						for _, tick := range tick.PermutateGranularity() {
							data, err := json.Marshal(tick)
							if err != nil {
								cancel(err)
								continue ITERATE_PRICING_INTERVALS
							}

							if _, ok := duplicates[string(data)]; ok {
								continue
							}

							ticks = append(ticks, string(data))
							duplicates[string(data)] = 1
						}
					}
					if err := iter.Err(); err != nil {
						cancel(err)
						continue ITERATE_PRICING_INTERVALS
					}

					for _, tick := range ticks {
						_, err = db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s JSON %s;", tickTableName, tick))
						if err != nil {
							cancel(err)
							continue ITERATE_PRICING_INTERVALS
						}
					}

				}
				cancel(nil)
			}

			didDownloadPricingInformation = len(values) != 0
		}

		// print warnings
		if len(warnings) > 0 {
			for _, w := range warnings {
				print(fmt.Sprintf("⏩  WARNING: %s\n", w))
			}
		}

		if !didDownloadMetaInformation && !didDownloadPricingInformation {
			flag.PrintDefaults()
		}

		sig <- os.Interrupt
	}()

	<-sig
}
