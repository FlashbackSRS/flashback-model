package test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func jsonDeepEqual(data1, data2 []byte) (bool, error) {
	var o1, o2 interface{}
	
	if err := json.Unmarshal(data1, &o1); err != nil {
		return false, fmt.Errorf("Error unmarshaling string 1: %s", err.Error())
	}
	if err := json.Unmarshal(data2, &o2); err != nil {
		return false, fmt.Errorf("Error unmarshaling string 2: %s", err.Error())
	}
	return reflect.DeepEqual(o1,o2), nil
}

func gotExpected(t *testing.T, descr, got, expected string) {
	t.Errorf("%s\n     Got: %s\nExpected: %s\n", descr, got, expected)
}

func stringsEqual(t *testing.T, descr, got, expected string) {
	if got != expected {
		gotExpected(t, descr, got, expected)
	}
}
