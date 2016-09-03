package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func DeepEqualJSON(t *testing.T, descr string, i1, i2 interface{}) {
	data1 := Marshal(t, descr+" (interface 1)", i1)
	data2 := Marshal(t, descr+" (interface 2)", i2)
	JSONDeepEqual(t, descr, data1, data2)
}

func JSONDeepEqual(t *testing.T, descr string, data1, data2 []byte) {
	var o1, o2 interface{}

	if err := json.Unmarshal(data1, &o1); err != nil {
		t.Errorf("%s: Error unmarshaling string 1: %s", descr, err.Error())
		return
	}
	if err := json.Unmarshal(data2, &o2); err != nil {
		t.Errorf("%s: Error unmarshaling string 2: %s", descr, err.Error())
		return
	}
	if !reflect.DeepEqual(o1, o2) {
		gotExpected(t, descr, string(data1), string(data2))
	}
}

func printDiff(got, exp string) {
	udiff := difflib.UnifiedDiff{
		A:        strings.SplitAfter(exp, "\n"),
		FromFile: "expected",
		B:        strings.SplitAfter(got, "\n"),
		ToFile:   "got",
		Context:  2,
	}
	diff, err := difflib.GetUnifiedDiffString(udiff)
	if err != nil {
		panic("Error producing diff: " + err.Error())
	}
	fmt.Print(diff)
}

func gotExpected(t *testing.T, descr, got, expected string) {
	if got[len(got)-1:] != "\n" {
		got = got + "\n"
	}
	if expected[len(expected)-1:] != "\n" {
		expected = expected + "\n"
	}
	fmt.Printf("%s\n", descr)
	printDiff(got, expected)
	t.Error()
}

func StringsEqual(t *testing.T, descr, got, expected string) {
	if got != expected {
		gotExpected(t, descr, got, expected)
	}
}

func Marshal(t *testing.T, descr string, i interface{}) []byte {
	output, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		t.Errorf("%s: Error marshaling JSON: %s\n", descr, err)
		return nil
	}
	return output
}
