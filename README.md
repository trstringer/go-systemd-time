# go-systemd-time

![ci](https://github.com/trstringer/go-systemd-time/workflows/ci/badge.svg)

_Go implementation of systemd time (`man systemd.time`)_

In command line applications, it is convenient to use the notation since `-1day`, or in `5 hours`. This package takes that string (using systemd time specs) and converts it into `time.Duration`. There is also a helper function that can take the raw string time adjustment and a `time.Time` (or `nil` for `time.Now()`) object and apply the adjustment to immutably. See below for usage.

## Usage

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/trstringer/go-systemd-time/pkg/systemdtime"
)

func main() {
	now := time.Now()
	timeFormat := "3:04 PM on January 2, 2006"

	fmt.Printf("Now is %s\n", now.Format(timeFormat))

	adjustedTime, err := systemdtime.AdjustTime(&now, "2d")
	if err != nil {
		fmt.Printf("error adjusting time: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Adjusted is %s\n", adjustedTime.Format(timeFormat))
	/*
	   Now is 8:38 PM on September 30, 2020
	   Adjusted is 8:38 PM on October 2, 2020
	*/
}
```

## Bug reports and running tests

If a bug is found, please write a failing test to uncover the bug.

To run tests, navigate to the root directory and run `go test -v ./...`.
