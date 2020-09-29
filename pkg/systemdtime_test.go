package systemdtime

import (
	"testing"
	"time"
)

func TestUnitParsingValue(t *testing.T) {
	testInputs := map[string]time.Duration{
		"usec":    time.Microsecond,
		"us":      time.Microsecond,
		"msec":    time.Millisecond,
		"ms":      time.Millisecond,
		"s":       time.Second,
		"sec":     time.Second,
		"second":  time.Second,
		"seconds": time.Second,
		"m":       time.Minute,
		"min":     time.Minute,
		"minute":  time.Minute,
		"minutes": time.Minute,
		"h":       time.Hour,
		"hr":      time.Hour,
		"hour":    time.Hour,
		"hours":   time.Hour,
		"d":       24 * time.Hour,
		"day":     24 * time.Hour,
		"days":    24 * time.Hour,
		"w":       7 * 24 * time.Hour,
		"week":    7 * 24 * time.Hour,
		"weeks":   7 * 24 * time.Hour,
		"M":       time.Duration(30.44 * float64(24) * float64(time.Hour)),
		"month":   time.Duration(30.44 * float64(24) * float64(time.Hour)),
		"months":  time.Duration(30.44 * float64(24) * float64(time.Hour)),
		"y":       time.Duration(365.25 * float64(24) * float64(time.Hour)),
		"year":    time.Duration(365.25 * float64(24) * float64(time.Hour)),
		"years":   time.Duration(365.25 * float64(24) * float64(time.Hour)),
	}

	for input, expectedDuration := range testInputs {
		duration, err := UnitToDuration(input)
		if err != nil {
			t.Errorf("Error converting unit input to duration: %v", err)
		}

		if duration == expectedDuration {
			t.Logf("Input '%s' passed test with an expected duration of %d", input, duration)
		} else {
			t.Errorf("Input '%s' failed: expected duration %d but go %d", input, expectedDuration, duration)
		}
	}
}

func TestUnitParsingSuccess(t *testing.T) {
	testInputs := map[string]bool{
		"usecs":   false,
		"usec":    true,
		"us":      true,
		"msec":    true,
		"ms":      true,
		"xmsecs":  false,
		"s":       true,
		"sec":     true,
		"second":  true,
		"sonds":   false,
		"seconds": true,
		"m":       true,
		"min":     true,
		"minute":  true,
		"minutes": true,
		"minutos": false,
		"h":       true,
		"hr":      true,
		"hour":    true,
		"hours":   true,
		"hrour":   false,
		"d":       true,
		"day":     true,
		"days":    true,
		"dayz":    false,
		"w":       true,
		"week":    true,
		"weeks":   true,
		"wek":     false,
		"M":       true,
		"month":   true,
		"months":  true,
		"moths":   false,
		"y":       true,
		"year":    true,
		"years":   true,
		"yars":    false,
	}

	for input, shouldSucceed := range testInputs {
		_, err := UnitToDuration(input)
		if (err == nil) == shouldSucceed {
			t.Logf("Input '%s' passed test with shouldSucceed set to %t", input, shouldSucceed)
		} else {
			t.Errorf("Input '%s' failed: err is '%v', and shouldSucceed set to %t", input, err, shouldSucceed)
		}
	}
}

func TestToDuration(t *testing.T) {
	testInputs := map[string]time.Duration{
		" 2 min 23sec": 2*time.Minute + 23*time.Second,
		"1sec":         1 * time.Second,
		"1 hour":       1 * time.Hour,
		"-2day":        -2 * 24 * time.Hour,
		"10 minutes":   10 * time.Minute,
	}

	for input, output := range testInputs {
		duration, err := ToDuration(input)
		if err != nil {
			t.Errorf("%v", err)
		}

		if duration != output {
			t.Errorf("Duration %d is not equal to output %d", duration, output)
		}
	}
}

func TestAdjustTime(t *testing.T) {
	time1 := time.Date(2012, time.May, 12, 5, 0, 0, 0, time.UTC)
	time1Mod, err := AdjustTime(&time1, " 4 days 2 hr")
	if err != nil {
		t.Error(err)
	}
	expectedTime := time.Date(2012, time.May, 16, 7, 0, 0, 0, time.UTC)

	if time1Mod != expectedTime {
		t.Errorf("Expected %s but got %s", expectedTime, time1Mod)
	}
}
