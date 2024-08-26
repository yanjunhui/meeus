// Copyright 2013 Sonia Keys
// License: MIT

package deltat_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/yanjunhui/meeus/base"
	"github.com/yanjunhui/meeus/deltat"
	"github.com/yanjunhui/meeus/julian"
	"github.com/yanjunhui/meeus/sexa"
	"github.com/yanjunhui/meeus/unit"
)

func ExampleInterp10A() {
	// Example 10.a, p. 78.
	dt := deltat.Interp10A(julian.CalendarGregorianToJD(1977, 2, 18))
	fmt.Printf("%+.1f seconds\n", dt)
	// Output:
	// +47.6 seconds
}

func ExamplePoly1900to1997() {
	// Example 10.a, p. 78.
	jd := julian.TimeToJD(time.Date(1977, 2, 18, 3, 37, 40, 0, time.UTC))
	year := base.JDEToJulianYear(jd)
	fmt.Printf("julian year %.1f\n", year)
	fmt.Printf("%+.1f seconds\n", deltat.Poly1900to1997(jd))
	// Output:
	// julian year 1977.1
	// +47.1 seconds
}

func ExamplePolyBefore948() {
	// Example 10.b, p. 80.
	ΔT := deltat.PolyBefore948(333.1)
	UT := unit.TimeFromHour(6)
	TD := UT + ΔT
	fmt.Printf("%+.0f seconds\n", ΔT)
	fmt.Printf("333 February 6 at %m TD", sexa.FmtTime(TD))
	// Output:
	// +6146 seconds
	// 333 February 6 at 7ʰ42ᵐ TD
}

// Table 10.A p. 79 provides a way to test these polynomials
func TestPoly1800to1997(t *testing.T) {
	for _, tp := range []struct {
		year int
		ΔT   unit.Time
	}{
		{1800, 13.1},
		{1900, -2.8},
		{1996, 61.6},
	} {
		jd := julian.CalendarGregorianToJD(tp.year, 0, 0)
		ΔT := deltat.Poly1800to1997(jd)
		if math.Abs((ΔT - tp.ΔT).Sec()) > 2.3 {
			t.Fatalf("%#v, got %.1f", tp, ΔT)
		}
	}
}

func TestPoly1800to1899(t *testing.T) {
	for _, tp := range []struct {
		year int
		ΔT   unit.Time
	}{
		{1800, 13.1},
		{1850, 6.8},
		{1898, -4.7},
	} {
		jd := julian.CalendarGregorianToJD(tp.year, 0, 0)
		if math.Abs((deltat.Poly1800to1899(jd) - tp.ΔT).Sec()) > 1 {
			t.Fatalf("%#v", tp)
		}
	}
}

func TestPoly1900to1997(t *testing.T) {
	for y := 1900; y < 1998; y += 2 {
		jd := julian.CalendarGregorianToJD(y, 0, 0)
		t.Logf("%d %.2f  %.1f", y, jd, deltat.Poly1900to1997(jd))
	}
	for _, tp := range []struct {
		year int
		ΔT   unit.Time
	}{
		{1900, -2.8},
		{1950, 29.1},
		{1996, 61.6},
	} {
		jd := julian.CalendarGregorianToJD(tp.year, 0, 0)
		ΔT := deltat.Poly1900to1997(jd)
		if math.Abs((ΔT - tp.ΔT).Sec()) > 1 {
			t.Fatalf("%#v, got %.1f", tp, ΔT)
		}
	}
}
