package systemdtime

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var units = []struct {
	Regex    *regexp.Regexp
	Duration time.Duration
}{
	{regexp.MustCompile(`^us(ec)?$`), time.Microsecond},
	{regexp.MustCompile(`^ms(ec)?$`), time.Millisecond},
	{regexp.MustCompile(`^s(ec(onds?)?)?$`), time.Second},
	{regexp.MustCompile(`^m(in(utes?)?)?$`), time.Minute},
	{regexp.MustCompile(`^(hr|h(ours?)?)$`), time.Hour},
	{regexp.MustCompile(`^d(ays?)?$`), 24 * time.Hour},
	{regexp.MustCompile(`^w(eeks?)?$`), 7 * 24 * time.Hour},
	{regexp.MustCompile(`^(M|months?)$`), time.Duration(30.44 * 24 * float64(time.Hour))},
	{regexp.MustCompile(`^y(ears?)?$`), time.Duration(365.25 * 24 * float64(time.Hour))},
}

// UnitToDuration converts a systemd unit (e.g. "day") to time.Duration
func UnitToDuration(unit string) (time.Duration, error) {
	for _, u := range units {
		if u.Regex.MatchString(unit) {
			return u.Duration, nil
		}
	}
	return 0, fmt.Errorf("Unit %s did not match", unit)
}

// ToDuration converts a systemd relative time string into time.Duration
func ToDuration(raw string) (time.Duration, error) {
	re, err := regexp.Compile(`^\s*-?\s*(\d+\s*[a-z]+)`)
	if err != nil {
		return 0, err
	}

	if !re.MatchString(raw) {
		return 0, fmt.Errorf("ToDuration: incorrect format for raw input %s", raw)
	}

	reNegative, err := regexp.Compile(`^\s*-.*`)
	if err != nil {
		return 0, err
	}
	isNegative := reNegative.MatchString(raw)

	reGroups, err := regexp.Compile(`\d+\s*[a-z]+`)
	if err != nil {
		return 0, err
	}

	matches := reGroups.FindAllString(raw, -1)

	totalDuration := time.Duration(0)
	reSubGroup, err := regexp.Compile(`^(\d+)\s*([a-z]+)$`)
	if err != nil {
		return 0, err
	}
	for _, match := range matches {
		matchTrimmed := strings.Replace(match, " ", "", -1)
		subGroupMatches := reSubGroup.FindStringSubmatch(matchTrimmed)

		// if we run into a case where there aren't exactly two matches
		// then that means this is an unexpected string and we should error out
		if len(subGroupMatches) != 3 {
			return 0, fmt.Errorf("Unexpected match count for '%s': expected 2 and got %d", matchTrimmed, len(subGroupMatches))
		}

		subGroupMatchValue, err := strconv.Atoi(subGroupMatches[1])
		if err != nil {
			return 0, err
		}

		subGroupMatchUnit, err := UnitToDuration(subGroupMatches[2])
		if err != nil {
			return 0, err
		}

		totalDuration += time.Duration(subGroupMatchValue) * subGroupMatchUnit
	}

	if isNegative {
		totalDuration *= -1
	}

	return totalDuration, nil
}

// AdjustTime takes a systemd time adjustment string and uses it to modify a time.Time
func AdjustTime(original *time.Time, adjustment string) (time.Time, error) {
	duration, err := ToDuration(adjustment)
	if err != nil {
		return time.Time{}, err
	}

	return adjustTimeByDuration(original, duration), nil
}

func adjustTimeByDuration(original *time.Time, adjustment time.Duration) time.Time {
	if original == nil {
		rightNow := time.Now()
		original = &rightNow
	}

	return original.Add(adjustment)
}
