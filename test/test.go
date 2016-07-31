package test

import (
	"encoding/json"
	"reflect"
	"testing"
)

func jsonDeepEqual(t *testing.T, descr string, data1, data2 []byte) {
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
	t.Errorf("%s\n     Got: %s\nExpected: %s\n", descr, got, expected)
}

func stringsEqual(t *testing.T, descr, got, expected string) {
	if got != expected {
		gotExpected(t, descr, got, expected)
	}
}

func marshal(t *testing.T, descr string, i interface{}) []byte {
	output, err := json.Marshal(i)
	if err != nil {
		t.Errorf("%s: Error marshaling JSON: %s\n", descr, err)
		return nil
	}
	return output
}
