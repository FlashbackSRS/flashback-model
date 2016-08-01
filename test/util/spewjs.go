// +build js

package util

import (
	"fmt"
)

func PrintDiff(got, expected interface{}) {
	gotString := fmt.Sprintf("%v\n", got)
	expString := fmt.Sprintf("%v\n", expected)
	fmt.Printf("GopherJS does not support go-spew. Run this test under standard Go for nicer output.\n")
	printDiff(gotString, expString)
}
