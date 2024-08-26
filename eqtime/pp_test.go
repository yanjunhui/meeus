// Copyright 2013 Sonia Keys
// License: MIT

//go:build !nopp
// +build !nopp

package eqtime_test

import (
	"fmt"

	"github.com/yanjunhui/meeus/eqtime"
	"github.com/yanjunhui/meeus/julian"
	pp "github.com/yanjunhui/meeus/planetposition"
	"github.com/yanjunhui/meeus/sexa"
)

func ExampleE() {
	// Example 28.a, p. 184
	earth, err := pp.LoadPlanet(pp.Earth)
	if err != nil {
		fmt.Println(err)
		return
	}
	j := julian.CalendarGregorianToJD(1992, 10, 13)
	eq := eqtime.E(j, earth)
	fmt.Printf("%+.1d", sexa.FmtHourAngle(eq))
	// Output:
	// +13ᵐ42ˢ.6
}
