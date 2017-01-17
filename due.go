package fb

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// Duration unit available for scheduling
const (
	Second = time.Second
	Minute = time.Minute
	Hour   = time.Hour
	Day    = time.Hour * 24
)

// This allows overriding time.Now() for tests
var now = time.Now

var epoch = time.Unix(0, 0).UTC()

// Due formatting constants
const (
	DueDays    = "2006-01-02"
	DueSeconds = "2006-01-02 15:04:05"
)

// Due represents the time/date a card is due.
type Due time.Time

// ParseDue attempts to parse the provided string as a due time.
func ParseDue(src string) (Due, error) {
	if t, err := time.Parse(DueDays, src); err == nil {
		return Due(t), nil
	}
	if t, err := time.Parse(DueSeconds, src); err == nil {
		return Due(t), nil
	}
	return Due{}, fmt.Errorf("Unrecognized input: %s", src)
}

// DueIn returns a new Due time d duration into the future. Durations greater
// than 24 hours into the future are rounded to the day.
func DueIn(dur time.Duration) Due {
	return Due(now()).Add(Interval(dur))
}

// Add returns a new Due time with the duration added to it.
func (d Due) Add(ivl Interval) Due {
	dur := time.Duration(ivl)
	if dur.Hours() < 24 {
		return Due(time.Time(d).Add(dur))
	}
	// Round up to the next whole day
	return Due(time.Time(d).Truncate(Day).Add(dur + Day - time.Nanosecond).Truncate(Day))
}

func (d Due) String() string {
	t := time.Time(d)
	if t.Truncate(Day).Equal(t) {
		return t.Format(DueDays)
	}
	return t.Format(DueSeconds)
}

// Time converts the due date to a standard time.Time
func (d Due) Time() time.Time {
	return time.Time(d)
}

func midnight(t time.Time) time.Time {
	return t.UTC().Truncate(24 * time.Hour)
}

// MarshalJSON implements the json.Marshaler interface
func (d Due) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (d *Due) UnmarshalJSON(src []byte) error {
	due, err := ParseDue(string(bytes.Trim(src, "\"")))
	*d = due
	return err
}

// Interval represents the number of days or seconds between reviews
type Interval time.Duration

var unitMap = map[string]time.Duration{
	"s": Second,
	"m": Minute,
	"h": Hour,
	"d": Day,
}

// ParseInterval parses an interval string. An interval string is string
// containing a positive integer, follwed by a unit suffix. e.g. "300s" or "15d".
// No spaces or other characters are allowed. Valid units are "s", "m", "h", and
// "d".
func ParseInterval(s string) (Interval, error) {
	q, err := strconv.ParseInt(s[:len(s)-1], 10, 64)
	if err != nil {
		return 0, err
	}
	unit, ok := unitMap[s[len(s)-1:]]
	if !ok {
		return 0, fmt.Errorf("Unknown unit in '%s'", s)
	}
	return Interval(time.Duration(q) * unit), nil
}

func (i Interval) String() string {
	d := time.Duration(i)
	if d >= 24*Hour {
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
	s := int(d.Seconds())
	if s%3600 == 0 {
		return fmt.Sprintf("%dh", s/3600)
	}
	if s%60 == 0 {
		return fmt.Sprintf("%dm", s/60)
	}
	return fmt.Sprintf("%ds", s)
}

// MarshalJSON implements the json.Marshaler interface
func (i Interval) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", i)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *Interval) UnmarshalJSON(src []byte) error {
	ivl, err := ParseInterval(string(bytes.Trim(src, "\"")))
	*i = ivl
	return err
}
