package main

import "testing"

func Test(t *testing.T) {

	// asset := "MSFT"
	// equityFlags[1] = &asset

	asset := "^GSPC"
	tickFlag = &asset
	enableFlag := true
	fifteenMinTickIntervalFlags[1] = &enableFlag

	main()
	t.FailNow()
}
