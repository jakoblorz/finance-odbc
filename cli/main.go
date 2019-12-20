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

	_ "github.com/jakoblorz/dynsql/lib/go-sqlite3"
	fodbc "github.com/jakoblorz/finance-odbc"
	"github.com/jmoiron/sqlx"
)

var (
	cryptoTableName     = "crypto"
	equityTableName     = "equity"
	etfTableName        = "etf"
	forexTableName      = "forex"
	futureTableName     = "future"
	indexTableName      = "index"
	mutualfundTableName = "mutualfund"
	optionTableName     = "option"

	tableNamePadding = map[string]string{
		cryptoTableName:     "    ",
		equityTableName:     "    ",
		etfTableName:        "       ",
		forexTableName:      "     ",
		futureTableName:     "    ",
		indexTableName:      "     ",
		mutualfundTableName: "",
		optionTableName:     "    ",
	}
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
	cryptoFlags     = []interface{}{cryptoTableName, flag.String("crypto", "", "")}
	equityFlags     = []interface{}{equityTableName, flag.String("equity", "", "")}
	etfFlags        = []interface{}{etfTableName, flag.String("etf", "", "")}
	forexFlags      = []interface{}{forexTableName, flag.String("forex", "", "")}
	futureFlags     = []interface{}{futureTableName, flag.String("future", "", "")}
	indexFlags      = []interface{}{indexTableName, flag.String("index", "", "")}
	mutualfundFlags = []interface{}{mutualfundTableName, flag.String("mutualfund", "", "")}
	optionFlags     = []interface{}{optionTableName, flag.String("option", "", "")}
)

func parseFlagSet(fs []interface{}) (tableName string, value string, ok bool) {
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

func main() {
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

		for _, f := range [][]interface{}{cryptoFlags, equityFlags, etfFlags, forexFlags, futureFlags, indexFlags, mutualfundFlags, optionFlags} {
			tableName, value, ok := parseFlagSet(f)
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
		}

		sig <- os.Interrupt
	}()

	<-sig
}
