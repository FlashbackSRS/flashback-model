package fb

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// Duration unit available for scheduling
const (
	Second = Interval(time.Second)
	Minute = Interval(time.Minute)
	Hour   = Interval(time.Hour)
	Day    = Interval(time.Hour * 24)
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

// DueIn returns a new Due time i interval into the future. Durations greater
// than 24 hours into the future are rounded to the day.
func DueIn(i Interval) Due {
	return Due(now()).Add(i)
}

// Add returns a new Due time with the duration added to it.
func (d Due) Add(ivl Interval) Due {
	dur := time.Duration(ivl)
	if ivl < Day {
		return Due(time.Time(d).Add(dur))
	}
	// Round up to the next whole day
	return Due(time.Time(d).Truncate(time.Duration(Day)).Add(time.Duration(ivl + Day - 1)).Truncate(time.Duration(Day)))
}

// Sub returns the interval between d and s
func (d Due) Sub(s Due) Interval {
	return Interval(time.Time(d).Sub(time.Time(s)))
}

// Equal returns true if the two due dates are equal.
func (d Due) Equal(d2 Due) bool {
	t1 := time.Time(d)
	t2 := time.Time(d2)
	return t1.Equal(t2)
}

// After returns true if d2 is after d.
func (d Due) After(d2 Due) bool {
	t1 := time.Time(d)
	t2 := time.Time(d2)
	return t1.After(t2)
}

func (d Due) String() string {
	t := time.Time(d)
	if t.Truncate(time.Duration(Day)).Equal(t) {
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

var unitMap = map[string]Interval{
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
	return Interval(q) * unit, nil
}

func (i Interval) String() string {
	if i >= Day {
		return fmt.Sprintf("%dd", int(time.Duration(i).Hours()/24))
	}
	s := int(time.Duration(i).Seconds())
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
