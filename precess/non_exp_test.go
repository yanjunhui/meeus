// Copyright 2013 Sonia Keys
// License: MIT

package precess

import (
	"math"
	"testing"

	"github.com/yanjunhui/meeus/base"
	"github.com/yanjunhui/meeus/coord"
	"github.com/yanjunhui/meeus/nutation"
	"github.com/yanjunhui/unit"
)

// test data from p. 132.
func TestMn(t *testing.T) {
	epochFrom := 2000.0
	for _, e := range []struct{ epoch, m, na, nd float64 }{
		{1700, 3.069, 1.338, 20.07},
		{1800, 3.071, 1.337, 20.06},
		{1900, 3.073, 1.337, 20.05},
		{2000, 3.075, 1.336, 20.04},
		{2100, 3.077, 1.336, 20.03},
		{2200, 3.079, 1.335, 20.03},
	} {
		m, na, nd := mn(epochFrom, e.epoch)
		if math.Abs(m.Sec()-e.m) > 1e-3 {
			t.Fatal("m:", m, e.m)
		}
		if math.Abs(na.Sec()-e.na) > 1e-3 {
			t.Fatal("na:", na, e.na)
		}
		if math.Abs(nd.Sec()-e.nd) > 1e-2 {
			t.Fatal("nd:", nd, e.nd)
		}
	}
}

// Test with proper motion of Regulus, with equatorial motions given
// in Example 21.a, p. 132, and ecliptic motions given in table 21.A,
// p. 138.
func TestEqProperMotionToEcl(t *testing.T) {
	ε := coord.NewObliquity(nutation.MeanObliquity(base.J2000))
	mλ, mβ := eqProperMotionToEcl(
		// eq motions from p. 132.
		unit.NewHourAngle('-', 0, 0, 0.0169),
		unit.NewAngle(' ', 0, 0, 0.006),
		2000.0,
		// eq coordinates from p. 132.
		new(coord.Ecliptic).EqToEcl(&coord.Equatorial{
			RA:  unit.NewRA(10, 8, 22.3),
			Dec: unit.NewAngle(' ', 11, 58, 2),
		}, ε))
	d := math.Abs((mλ - unit.AngleFromSec(-.2348)).Rad() / mλ.Rad())
	if d*169 > 1 { // 169 = significant digits of given lon
		t.Fatal("mλ")
	}
	d = math.Abs((mβ - unit.AngleFromSec(-.0813)).Rad() / mβ.Rad())
	if d*6 > 1 { // 6 = significant digit of given lat
		t.Fatal("mβ")
	}
}
