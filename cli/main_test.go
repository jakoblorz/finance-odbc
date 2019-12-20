package main

import "testing"

func Test(t *testing.T) {
	asset := "MSFT"
	equityFlags[1] = &asset
	main()
	t.FailNow()
}
