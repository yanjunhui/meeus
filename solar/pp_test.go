// Copyright 2013 Sonia Keys
// License: MIT

//go:build !nopp
// +build !nopp

package solar_test

import (
	"fmt"

	"github.com/yanjunhui/meeus/julian"
	pp "github.com/yanjunhui/meeus/planetposition"
	"github.com/yanjunhui/meeus/sexa"
	"github.com/yanjunhui/meeus/solar"
)

func ExampleApparentEquatorialVSOP87() {
	// Example 25.b, p. 169, but as this code uses the full VSOP87 theory,
	// results match those at bottom of p. 165.
	e, err := pp.LoadPlanet(pp.Earth)
	if err != nil {
		fmt.Println(err)
		return
	}
	jde := julian.CalendarGregorianToJD(1992, 10, 13)
	α, δ, _ := solar.ApparentEquatorialVSOP87(e, jde)
	fmt.Printf("α: %.3d\n", sexa.FmtRA(α))
	fmt.Printf("δ: %+.2d\n", sexa.FmtAngle(δ))
	// Output:
	// α: 13ʰ13ᵐ30ˢ.749
	// δ: -7°47′1″.74
}
