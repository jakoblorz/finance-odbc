package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"

	"github.com/jakoblorz/dynsql/lib/go-sqlite3"
	fodbc "github.com/jakoblorz/finance-odbc"
	"github.com/jmoiron/sqlx"
)

var (
	cryptoTableName     = "crypto"
	equityTableName     = "equity"
	etfTableName        = "etf"
	forexTableName      = "forex"
	futureTableName     = "future"
	stockindexTableName = "stockindex"
	mutualfundTableName = "mutualfund"
	optionTableName     = "option"

	tableNamePadding = map[string]string{
		cryptoTableName:     "    ",
		equityTableName:     "    ",
		etfTableName:        "       ",
		forexTableName:      "     ",
		futureTableName:     "    ",
		stockindexTableName: "",
		mutualfundTableName: "",
		optionTableName:     "    ",
	}

	tickTableNamePrefix = "tick"
)

type financeConvertAPIFunc func(interface{}) (interface{}, bool)
type financeObtainAPIFunc func(string) (interface{}, error)

var (
	financeControlFuncMap = map[string][]interface{}{
		cryptoTableName: []interface{}{
			fodbc.NewAnonCryptoFromAPI,
			fodbc.GetAnonCryptoFromAPI,
		},
		equityTableName: []interface{}{
			fodbc.NewAnonEquityFromAPI,
			fodbc.GetAnonEquityFromAPI,
		},
		etfTableName: []interface{}{
			fodbc.NewAnonETFFromAPI,
			fodbc.GetAnonETFFromAPI,
		},
		forexTableName: []interface{}{
			fodbc.NewAnonForexFromAPI,
			fodbc.GetAnonForexFromAPI,
		},
		futureTableName: []interface{}{
			fodbc.NewAnonFutureFromAPI,
			fodbc.GetAnonFutureFromAPI,
		},
		stockindexTableName: []interface{}{
			fodbc.NewAnonIndexFromAPI,
			fodbc.GetAnonIndexFromAPI,
		},
		mutualfundTableName: []interface{}{
			fodbc.NewAnonMutualFundFromAPI,
			fodbc.GetAnonMutualFundFromAPI,
		},
		optionTableName: []interface{}{
			fodbc.NewAnonOptionFromAPI,
			fodbc.GetAnonOptionFromAPI,
		},
	}
)

var (
	cryptoFlags     = []interface{}{cryptoTableName, flag.String(cryptoTableName, "", "Download Crypto Information")}
	equityFlags     = []interface{}{equityTableName, flag.String(equityTableName, "", "Download Equity Information")}
	etfFlags        = []interface{}{etfTableName, flag.String(etfTableName, "", "Download ETF Information")}
	forexFlags      = []interface{}{forexTableName, flag.String(forexTableName, "", "Download Forex Information")}
	futureFlags     = []interface{}{futureTableName, flag.String(futureTableName, "", "Download Future Information")}
	stockindexFlags = []interface{}{stockindexTableName, flag.String(stockindexTableName, "", "Download Index Information")}
	mutualfundFlags = []interface{}{mutualfundTableName, flag.String(mutualfundTableName, "", "Download Mututal Fund Information")}
	optionFlags     = []interface{}{optionTableName, flag.String(optionTableName, "", "Donwload Option Information")}

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

	tickFlag        = flag.String(tickTableNamePrefix, "", "Download Pricing Information")
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
	ctx, cancel := context.WithCancel(context.Background())
	if final == "" {
		final = "DONE ✅\n"
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Prefix = prefix
	s.Start()

	sig := make(chan int, 0)
	go func() {
		<-ctx.Done()
		s.Stop()
		close(sig)
	}()

	return func(err error) {
		cancel()
		if err != nil {
			print(fmt.Sprintf("FAILED ❌\n%s\n", err))
		} else {
			print(final)
		}
		<-sig
	}
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
		for _, f := range allMetaInformationFlags {
			tableName, value, ok := parseFlagSetS(f)
			if !ok || value == "" {
				continue
			}

			controlFuncs := financeControlFuncMap[tableName]
			assetConvertAPI := controlFuncs[0].(func(interface{}) (interface{}, bool))
			assetObtainAPI := controlFuncs[1].(func(string) (interface{}, error))

			values := strings.Split(value, ",")
			cancel := spin(
				fmt.Sprintf("Downloading Metadata for %s:%s %d Download(s) required ", tableName, tableNamePadding[tableName], len(values)),
				"",
			)

			for _, value := range values {
				retrieved, err := assetObtainAPI(value)
				if err != nil {
					cancel(err)
					os.Exit(1)
				}
				asset, _ := assetConvertAPI(retrieved)

				data, err := json.Marshal(asset)
				if err != nil {
					cancel(err)
					os.Exit(1)
				}

				_, err = db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s JSON %s;", tableName, string(data)))
				if err != nil {
					cancel(err)
					os.Exit(1)
				}
			}
			cancel(nil)

			didDownloadMetaInformation = true
		}

		didDownloadPricingInformation := false
		if *tickFlag != "" {
			values := strings.Split(*tickFlag, ",")
			tUTCNow := time.Now().UTC()

			for _, f := range allPricingInformationFlags {
				interval, doDownload, ok := parseFlagSetB(f)
				if !(*useAllTicksFlag) && (!ok || !doDownload) {
					continue
				}

				cancel := spin(
					fmt.Sprintf("Downloading Historical Prices with an interval of %s: %d Download(s) required ", interval, len(values)),
					"",
				)
				for _, value := range values {

					iter := chart.Get(&chart.Params{
						Symbol:   value,
						End:      datetime.New(&tUTCNow),
						Start:    datetime.FromUnix(int(tUTCNow.Unix() - 1000*60*60)),
						Interval: datetime.Interval(interval),
					})

					ticks := []string{}
					for iter.Next() {
						tick := fodbc.NewTickFromAPI(&fodbc.MetaTick{
							ChartBar:  *iter.Bar(),
							ChartMeta: iter.Meta(),
						})

						data, err := json.Marshal(tick)
						if err != nil {
							cancel(err)
							os.Exit(1)
						}

						ticks = append(ticks, string(data))
					}
					if err := iter.Err(); err != nil {
						cancel(err)
						os.Exit(1)
					}

					for _, tick := range ticks {
						_, err = db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %ss_%s JSON %s;", tickTableNamePrefix, interval, tick))
						if err != nil {
							cancel(err)
							os.Exit(1)
						}
					}

				}
				cancel(nil)
			}

			didDownloadPricingInformation = len(values) != 0
		}

		if !didDownloadMetaInformation && !didDownloadPricingInformation {
			flag.PrintDefaults()
		}

		sig <- os.Interrupt
	}()

	<-sig
}
