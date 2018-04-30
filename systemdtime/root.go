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

var (
	reAdjustment = regexp.MustCompile(`^\s*-?\s*(\d+\s*[a-z]+)`)
	reNegative   = regexp.MustCompile(`^\s*-.*`)
	reGroups     = regexp.MustCompile(`\d+\s*[a-z]+`)
	reParts      = regexp.MustCompile(`^(\d+)\s*([a-z]+)$`)
)

// ToDuration converts a systemd relative time string into time.Duration
func ToDuration(raw string) (time.Duration, error) {
	if !reAdjustment.MatchString(raw) {
		return 0, fmt.Errorf("ToDuration: incorrect format for raw input %s", raw)
	}
	var total time.Duration
	for _, group := range reGroups.FindAllString(raw, -1) {
		group := strings.TrimSpace(group)
		parts := reParts.FindStringSubmatch(group)

		// if we run into a case where there aren't exactly two matches
		// then that means this is an unexpected string and we should error out
		if len(parts) != 3 {
			return 0, fmt.Errorf("Unexpected match count for '%s': expected 2 and got %d", group, len(parts))
		}
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		unit, err := UnitToDuration(parts[2])
		if err != nil {
			return 0, err
		}
		total += time.Duration(value) * unit
	}
	if reNegative.MatchString(raw) {
		total = -total
	}
	return total, nil
}

// AdjustTime takes a systemd time adjustment string and uses it to modify a time.Time
func AdjustTime(t time.Time, adjustment string) (time.Time, error) {
	duration, err := ToDuration(adjustment)
	if err != nil {
		return time.Time{}, err
	}
	if t.IsZero() {
		t = time.Now()
	}
	return t.Add(duration), nil
}
