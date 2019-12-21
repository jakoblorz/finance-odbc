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

	tickFlag = flag.String(tickTableNamePrefix, "", "Download Pricing Information")

	oneMinTickIntervalFlags     = []interface{}{string(datetime.OneMin), flag.Bool(string(datetime.OneMin), false, "")}
	twoMinTickIntervalFlags     = []interface{}{string(datetime.TwoMins), flag.Bool(string(datetime.TwoMins), false, "")}
	fiveMinTickIntervalFlags    = []interface{}{string(datetime.FiveMins), flag.Bool(string(datetime.FiveMins), false, "")}
	fifteenMinTickIntervalFlags = []interface{}{string(datetime.FifteenMins), flag.Bool(string(datetime.FifteenMins), false, "")}
	thirtyMinTickIntervalFlags  = []interface{}{string(datetime.ThirtyMins), flag.Bool(string(datetime.ThirtyMins), false, "")}
	sixtyMinTickIntervalFlags   = []interface{}{string(datetime.SixtyMins), flag.Bool(string(datetime.SixtyMins), false, "")}
	ninetyMinTickIntervalFlags  = []interface{}{string(datetime.NinetyMins), flag.Bool(string(datetime.NinetyMins), false, "")}
	oneHourTickIntervalFlags    = []interface{}{string(datetime.OneHour), flag.Bool(string(datetime.OneHour), false, "")}
	oneDayTickIntervalFlags     = []interface{}{string(datetime.OneDay), flag.Bool(string(datetime.OneDay), false, "")}
	fiveDayTickIntervalFlags    = []interface{}{string(datetime.FiveDay), flag.Bool(string(datetime.FiveDay), false, "")}

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
			log.Fatal(err)
		}
		defer db.Close()

		didDownloadMetaInformation := false
		log.Println("Checking Meta Information Download")
		for _, f := range allMetaInformationFlags {
			tableName, value, ok := parseFlagSetS(f)
			if !ok || value == "" {
				log.Printf("Checking download request for %s:%s DONE", tableName, tableNamePadding[tableName])
				continue
			}

			controlFuncs := financeControlFuncMap[tableName]
			assetConvertAPI := controlFuncs[0].(func(interface{}) (interface{}, bool))
			assetObtainAPI := controlFuncs[1].(func(string) (interface{}, error))

			values := strings.Split(value, ",")
			log.Printf("Checking download request for %s:%s %d Download(s) required", tableName, tableNamePadding[tableName], len(values))

			for i, value := range values {
				log.Printf(" (%d.1/%d) %s Downloading\n", i+1, len(values), value)

				retrieved, err := assetObtainAPI(value)
				if err != nil {
					log.Fatal(err)
				}
				asset, _ := assetConvertAPI(retrieved)

				data, err := json.Marshal(asset)
				if err != nil {
					log.Fatal(err)
				}

				log.Printf(" (%d.2/%d) %s Inserting Data\n", i+1, len(values), value)
				_, err = db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s JSON %s;", tableName, string(data)))
				if err != nil {
					log.Fatal(err)
				}
			}

			didDownloadMetaInformation = true
		}

		didDownloadPricingInformation := false
		log.Println("Checking Pricing Information Download")
		if *tickFlag != "" {
			values := strings.Split(*tickFlag, ",")
			tUTCNow := time.Now().UTC()

			for _, f := range allPricingInformationFlags {
				interval, doDownload, ok := parseFlagSetB(f)
				if !ok || !doDownload {
					log.Printf("Checking download request for %s ticks: DONE", interval)
					continue
				}

				for i, value := range values {
					log.Printf(" (%d.1/%d) %s Downloading %s ticks", i+1, len(values), value, interval)

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
							log.Fatal(err)
						}

						ticks = append(ticks, string(data))
					}
					if err := iter.Err(); err != nil {
						log.Fatal(err)
					}

					log.Printf(" (%d.2/%d) %s Inserting a total of %d %s ticks", i+1, len(values), value, len(ticks), interval)
					for _, tick := range ticks {
						_, err = db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %ss_%s JSON %s;", tickTableNamePrefix, interval, tick))
						if err != nil {
							log.Fatal(err)
						}
					}
				}
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
