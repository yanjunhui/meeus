// Copyright 2013 Sonia Keys
// License: MIT

package planetelements_test

import (
	"fmt"
	"testing"

	"github.com/yanjunhui/meeus/julian"
	pe "github.com/yanjunhui/meeus/planetelements"
)

func ExampleMean() {
	// Example 31.a, p. 211
	j := julian.CalendarGregorianToJD(2065, 6, 24)
	var e pe.Elements
	pe.Mean(pe.Mercury, j, &e)
	fmt.Printf("L: %.6f\n", e.Lon.Deg())
	fmt.Printf("a: %.9f\n", e.Axis)
	fmt.Printf("e: %.8f\n", e.Ecc)
	fmt.Printf("i: %.6f\n", e.Inc.Deg())
	fmt.Printf("Ω: %.6f\n", e.Node.Deg())
	fmt.Printf("ϖ: %.6f\n", e.Peri.Deg())
	// Output:
	// L: 203.494701
	// a: 0.387098310
	// e: 0.20564510
	// i: 7.006171
	// Ω: 49.107650
	// ϖ: 78.475382
}

func TestInc(t *testing.T) {
	j := julian.CalendarGregorianToJD(2065, 6, 24)
	var e pe.Elements
	pe.Mean(pe.Mercury, j, &e)
	if i := pe.Inc(pe.Mercury, j); i != e.Inc {
		t.Fatal(i, "!=", e.Inc)
	}
}

func TestNode(t *testing.T) {
	j := julian.CalendarGregorianToJD(2065, 6, 24)
	var e pe.Elements
	pe.Mean(pe.Mercury, j, &e)
	if Ω := pe.Node(pe.Mercury, j); Ω != e.Node {
		t.Fatal(Ω, "!=", e.Node)
	}
}
