// Copyright 2013 Sonia Keys
// License: MIT

package solardisk_test

import (
	"fmt"
	"time"

	"github.com/yanjunhui/meeus/julian"
	"github.com/yanjunhui/meeus/solardisk"
)

func ExampleCycle() {
	j := solardisk.Cycle(1699)
	fmt.Printf("%.4f\n", j)
	y, m, d := julian.JDToCalendar(j)
	fmt.Printf("%d %s %.2f\n", y, time.Month(m), d)
	// Output:
	// 2444480.7230
	// 1980 August 29.22
}
