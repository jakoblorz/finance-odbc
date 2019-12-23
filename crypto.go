package odbc

import (
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/crypto"
)

type CryptoPair struct {
	Quote

	Algorithm           string `db:"algorithm" json:"algorithm"`
	StartDate           int    `db:"start_date" json:"start_date"`
	MaxSupply           int    `db:"max_supply" json:"max_supply"`
	CirculatingSupply   int    `db:"circulating_supply" json:"circulating_supply"`
	VolumeLastDay       int    `db:"volume_last_day" json:"volume_last_day"`
	VolumeAllCurrencies int    `db:"volume_all_currencies" json:"volume_all_currencies"`
}

func NewAnonCryptoFromAPI(c interface{}) (p interface{}, ok bool) {
	var apiCryptoPair *finance.CryptoPair
	apiCryptoPair, ok = c.(*finance.CryptoPair)
	if !ok {
		return
	}

	p, ok = NewCryptoFromAPI(apiCryptoPair)
	return
}

func GetAnonCryptoFromAPI(symbol string) (interface{}, error) {
	return crypto.Get(symbol)
}

func NewCryptoFromAPI(c *finance.CryptoPair) (cp CryptoPair, ok bool) {
	ok = c != nil && (&c.Quote) != nil
	if !ok {
		return
	}

	cp = CryptoPair{
		Quote: NewQuoteFromAPI(&c.Quote),

		Algorithm:           c.Algorithm,
		StartDate:           c.StartDate,
		MaxSupply:           c.MaxSupply,
		CirculatingSupply:   c.CirculatingSupply,
		VolumeLastDay:       c.VolumeLastDay,
		VolumeAllCurrencies: c.VolumeAllCurrencies,
	}
	return
}
