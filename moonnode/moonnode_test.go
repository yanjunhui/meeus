// Copyright 2013 Sonia Keys
// License: MIT

package moonnode_test

import (
	"fmt"
	"math"
	"time"

	"github.com/soniakeys/sexagesimal"
	"github.com/yanjunhui/meeus/julian"
	"github.com/yanjunhui/meeus/moonnode"
	"github.com/yanjunhui/unit"
)

func ExampleAscending() {
	// Example 51.a, p. 365.
	j := moonnode.Ascending(1987.37)
	fmt.Printf("%.5f\n", j)
	y, m, d := julian.JDToCalendar(j)
	d, f := math.Modf(d)
	fmt.Printf("%d %s %d, at %d TD\n", y, time.Month(m), int(d),
		sexa.FmtTime(unit.TimeFromDay(f)))
	// Output:
	// 2446938.76803
	// 1987 May 23, at 6ʰ25ᵐ58ˢ TD
}
