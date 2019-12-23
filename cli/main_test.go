package main

import "testing"

func Test(t *testing.T) {

	DEBUG = true

	cryptoAssets := "GDAXI,acb,exf"
	cryptoFlags[1] = &cryptoAssets

	equityAssets := "APPL,AAPL"
	equityFlags[1] = &equityAssets

	indexAssets := "^GDAXI"
	stockindexFlags[1] = &indexAssets

	tickAssets := "AAPL"
	tickFlag = &tickAssets

	enableFlag := true
	oneDayTickIntervalFlags[1] = &enableFlag
	oneMonthTickIntervalFlags[1] = &enableFlag

	main()
	t.FailNow()
}
