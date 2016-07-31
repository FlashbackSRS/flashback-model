package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

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

func gotExpected(t *testing.T, descr, got, expected string) {
	udiff := difflib.UnifiedDiff{
		A:        strings.SplitAfter(expected, "\n"),
		FromFile: "expected",
		B:        strings.SplitAfter(got, "\n"),
		ToFile:   "got",
	}
	diff, err := difflib.GetUnifiedDiffString(udiff)
	if err != nil {
		t.Fatal("Error producing diff: %s\n", err)
	}
	fmt.Printf("%s\n")
	fmt.Print(diff)
	t.Error()
	// 	t.Errorf("%s\n%s", descr, diff)
	// 	t.Errorf("%s\n     Got: %s\nExpected: %s\n", descr, got, expected)
}

func StringsEqual(t *testing.T, descr, got, expected string) {
	if got != expected {
		gotExpected(t, descr, got, expected)
	}
}

func Marshal(t *testing.T, descr string, i interface{}) []byte {
	output, err := json.Marshal(i)
	if err != nil {
		t.Errorf("%s: Error marshaling JSON: %s\n", descr, err)
		return nil
	}
	return output
}
