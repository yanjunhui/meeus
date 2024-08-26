// Copyright 2013 Sonia Keys
// License: MIT

package moonposition_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/yanjunhui/meeus/base"
	"github.com/yanjunhui/meeus/julian"
	"github.com/yanjunhui/meeus/moonposition"
)

func ExamplePosition() {
	// Example 47.a, p. 342.
	λ, β, Δ := moonposition.Position(julian.CalendarGregorianToJD(1992, 4, 12))
	fmt.Printf("λ = %.6f\n", λ.Deg())
	fmt.Printf("β = %.6f\n", β.Deg())
	fmt.Printf("Δ = %.1f\n", Δ)
	// Output:
	// λ = 133.162655
	// β = -3.229126
	// Δ = 368409.7
}

func ExampleParallax() {
	// Example 47.a, p. 342.
	_, _, Δ := moonposition.Position(julian.CalendarGregorianToJD(1992, 4, 12))
	π := moonposition.Parallax(Δ)
	fmt.Printf("π = %.6f\n", π.Deg())
	// Output:
	// π = 0.991990
}

func TestParallax(t *testing.T) {
	// test case from ch 40, p. 280
	got := moonposition.Parallax(.37276 * base.AU).Sec()
	want := 23.592
	if math.Abs(got-want) > .001 {
		t.Error(got, want)
	}
}

// Test data p. 344.
var n0 = []float64{
	julian.CalendarGregorianToJD(1913, 5, 27),
	julian.CalendarGregorianToJD(1932, 1, 6),
	julian.CalendarGregorianToJD(1950, 8, 17),
	julian.CalendarGregorianToJD(1969, 3, 29),
	julian.CalendarGregorianToJD(1987, 11, 8),
	julian.CalendarGregorianToJD(2006, 6, 19),
	julian.CalendarGregorianToJD(2025, 1, 29),
	julian.CalendarGregorianToJD(2043, 9, 10),
	julian.CalendarGregorianToJD(2062, 4, 22),
	julian.CalendarGregorianToJD(2080, 12, 1),
	julian.CalendarGregorianToJD(2099, 7, 13),
}

var n180 = []float64{
	julian.CalendarGregorianToJD(1922, 9, 16),
	julian.CalendarGregorianToJD(1941, 4, 27),
	julian.CalendarGregorianToJD(1959, 12, 7),
	julian.CalendarGregorianToJD(1978, 7, 19),
	julian.CalendarGregorianToJD(1997, 2, 27),
	julian.CalendarGregorianToJD(2015, 10, 10),
	julian.CalendarGregorianToJD(2034, 5, 21),
	julian.CalendarGregorianToJD(2052, 12, 30),
	julian.CalendarGregorianToJD(2071, 8, 12),
	julian.CalendarGregorianToJD(2090, 3, 23),
	julian.CalendarGregorianToJD(2108, 11, 3),
}

func TestNode0(t *testing.T) {
	for i, j := range n0 {
		if e := math.Abs(moonposition.Node(j).Rad()); e > 1e-3 {
			t.Error(i, e)
		}
	}
}

func TestNode180(t *testing.T) {
	for i, j := range n180 {
		if e := math.Abs(moonposition.Node(j).Rad() - math.Pi); e > 1e-3 {
			t.Error(i, e)
		}
	}
}
