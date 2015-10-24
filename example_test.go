package local_test

import (
	"fmt"
	"time"

	"github.com/spkg/local"
)

func ExampleToday() {
	_, month, day := local.Today().Date()
	if month == time.November && day == 10 {
		fmt.Println("Happy Go day!")
	}
}

func ExampleDateFor() {
	d := local.DateFor(2009, time.November, 10)
	fmt.Printf("Go launched on %s\n", d)
	// Output: Go launched on 2009-11-10
}

func ExampleDateParseLayout() {
	// longForm shows by example how the reference date would be represented in
	// the desired layout.
	const longForm = "Jan 2, 2006"
	d, _ := local.DateParseLayout(longForm, "Sep 30, 2099")
	fmt.Println(d)

	// shortForm is another way the reference date would be represented
	// in the desired layout.
	const shortForm = "2006-Jan-02"
	d, _ = local.DateParseLayout(shortForm, "2092-Dec-16")
	fmt.Println(d)

	// Output:
	// 2099-09-30
	// 2092-12-16
}

func ExampleDateTimeParseLayout() {
	// longForm shows by example how the reference date would be represented in
	// the desired layout.
	const longForm = "Jan 2, 2006 03:04:05PM"
	dt, _ := local.DateTimeParseLayout(longForm, "Sep 30, 2099 06:48:30PM")
	fmt.Println(dt)

	// shortForm is another way the reference date would be represented
	// in the desired layout.
	const shortForm = "2006-Jan-02 15:04:05"
	dt, _ = local.DateTimeParseLayout(shortForm, "2092-Dec-16 11:47:00")
	fmt.Println(dt)

	// Output:
	// 2099-09-30T18:48:30
	// 2092-12-16T11:47:00
}
