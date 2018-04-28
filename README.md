# systemdtime

Go implementation of systemd time (`man systemd.time`)

## Installation

You can use `dep` to install this vendor dependency, or `go get github.com/trstring/go-systemd-time/systemdtime`.

## Usage

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/trstringer/go-systemd-time/systemdtime"
)

func main() {
	now := time.Now()
	fmt.Printf("Current date/time: %s\n", now)

	twoDaysThreeHoursAgo, err := systemdtime.AdjustTime(&now, "-2 days 3 hr")
	if err != nil {
		fmt.Printf("Error converting: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Two days and three hours ago: %s\n", twoDaysThreeHoursAgo)

	eighteenDaysTwelveMinutesFromNow, err := systemdtime.AdjustTime(&now, "18d12min")
	if err != nil {
		fmt.Printf("Error converting: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Eighteen days and twelve minutes from now: %s\n", eighteenDaysTwelveMinutesFromNow)
}
```

Output from the above sample...

```
Current date/time: 2018-04-28 10:14:10.545534601 -0400 EDT m=+0.000220259      
Two days and three hours ago: 2018-04-26 07:14:10.545534601 -0400 EDT m=-183599.999779741                                                                      
Eighteen days and twelve minutes from now: 2018-05-16 10:26:10.545534601 -0400 EDT m=+1555920.000220259
```

## Bug reports and running tests

If a bug is found, please write a failing test to uncover the bug.

To run tests, navigate to the root directory run `go test ./test/`.
