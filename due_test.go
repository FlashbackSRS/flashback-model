package fb

import (
	"fmt"
	"testing"
	"time"
)

func init() {
	now = func() time.Time {
		return parseTime("2017-01-01T00:00:00Z")
	}
}

func parseTime(src string) time.Time {
	t, err := time.Parse(time.RFC3339, src)
	if err != nil {
		panic(err)
	}
	return t
}

func TestParseDue(t *testing.T) {
	if _, err := ParseDue("foobar"); err == nil {
		t.Errorf("Expected an error for invalid input to ParseDue")
	}

	expectedDue := parseTime("2017-01-01T00:00:00Z")
	result, err := ParseDue("2017-01-01")
	if err != nil {
		t.Errorf("Error parsing date-formatted Due value: %s", err)
	}
	if !expectedDue.Equal(result.Time()) {
		t.Errorf("Due = %s, expected %s\n", result, expectedDue)
	}

	expectedDue = parseTime("2017-01-01T12:30:45Z")
	result, err = ParseDue("2017-01-01 12:30:45")
	if err != nil {
		t.Errorf("Error parsing time-formatted Due value: %s", err)
	}
	if !expectedDue.Equal(result.Time()) {
		t.Errorf("Due = %s, expected %s\n", result, expectedDue)
	}
}

type StringerTest struct {
	Name     string
	I        fmt.Stringer
	Expected string
}

func TestStringer(t *testing.T) {
	tests := []StringerTest{
		StringerTest{
			Name:     "Interval seconds",
			I:        Interval(100 * Second),
			Expected: "100s",
		},
		StringerTest{
			Name:     "Interval seconds, plus nanoseconds",
			I:        Interval(10*Second + 1),
			Expected: "10s",
		},
		StringerTest{
			Name:     "Interval seconds",
			I:        Interval(5 * Minute),
			Expected: "5m",
		},
		StringerTest{
			Name:     "Interval seconds",
			I:        Interval(6 * Hour),
			Expected: "6h",
		},
		StringerTest{
			Name:     "Interval days",
			I:        Interval(100 * Hour),
			Expected: "5d",
		},
		StringerTest{
			Name:     "Due seconds",
			I:        Due(parseTime("2017-01-17T00:01:40Z")),
			Expected: "2017-01-17 00:01:40",
		},
		StringerTest{
			Name:     "Due days",
			I:        Due(parseTime("1970-04-11T00:00:00Z")),
			Expected: "1970-04-11",
		},
	}
	for _, test := range tests {
		if result := test.I.String(); result != test.Expected {
			t.Errorf("%s:\n\tExpected '%s'\n\t  Actual: '%s'\n", test.Name, test.Expected, result)
		}
	}
}

func TestDueIn(t *testing.T) {
	result := DueIn(10 * Minute)
	expected := "2017-01-01 00:10:00"
	if result.String() != expected {
		t.Errorf("Due in 10 minutes:\n\tExpected: %s\n\t  Actual: %s\n", expected, result)
	}

	result = DueIn(15 * Day)
	expected = "2017-01-16"
	if result.String() != expected {
		t.Errorf("Due in 15 days:\n\tExpected: %s\n\t  Actual: %s\n", expected, result)
	}
}

func TestAdd(t *testing.T) {
	result := DueIn(3 * Hour).Add(Interval(9000 * Second))
	expected := "2017-01-01 05:30:00"
	if result.String() != expected {
		t.Errorf("Add 9000 seconds2017-01-01 05:30:00:\n\tExpected: %s\n\t  Actual: %s\n", expected, result)
	}

	result = DueIn(3 * Hour).Add(Interval(24 * Hour))
	expected = "2017-01-02"
	if result.String() != expected {
		t.Errorf("Add 24 hours:\n\tExpected: %s\n\t  Actual: %s\n", expected, result)
	}

	result = DueIn(0).Add(Interval(24 * Hour))
	expected = "2017-01-02"
	if result.String() != expected {
		t.Errorf("Add 24 hours:\n\tExpected: %s\n\t  Actual: %s\n", expected, result)
	}

	result = DueIn(3 * Hour).Add(Interval(9000 * Hour))
	expected = "2018-01-11"
	if result.String() != expected {
		t.Errorf("Add 9000 hours:\n\tExpected: %s\n\t  Actual: %s\n", expected, result)
	}
}

func TestOn(t *testing.T) {
	ts, e := time.Parse(time.RFC3339, "2016-01-01T01:01:01+00:00")
	if e != nil {
		t.Fatal(e)
	}
	d := On(ts)
	expected, e := time.Parse("2006-01-02", "2016-01-01")
	if e != nil {
		t.Fatal(e)
	}
	if !time.Time(d).Equal(expected) {
		t.Errorf("Unexpected result: %v", d)
	}
}

func TestNow(t *testing.T) {
	now := time.Now()
	n := Now()
	if s := time.Time(n).Sub(now).Seconds(); s > 0.000001 {
		t.Errorf("Result differs by %fs", s)
	}
}

func TestToday(t *testing.T) {
	today := Today()
	expected := now().Truncate(time.Duration(Day))
	if !expected.Equal(time.Time(today)) {
		t.Errorf("Unepxected result: %v", today)
	}
}
