// +build !js disableunsafe

package util

import (
	"regexp"

	"github.com/davecgh/go-spew/spew"
)

var capRE *regexp.Regexp = regexp.MustCompile("cap=[0-9]+\\)")
var capRepl string = "cap=X"
var addRE = regexp.MustCompile("\\(0x[0-9a-f]{6,10}\\)")
var addRepl string = "(0xXXXXXXXXXX)"

// PrintDiff dumps a diff of the two data structures
func PrintDiff(got, expected interface{}) {
	scs := spew.ConfigState{
		Indent:         "  ",
		DisableMethods: true,
		SortKeys:       true,
	}
	gotString := scs.Sdump(got)
	expString := scs.Sdump(expected)

	gotString = capRE.ReplaceAllString(gotString, capRepl)
	expString = capRE.ReplaceAllString(expString, capRepl)
	gotString = addRE.ReplaceAllString(gotString, addRepl)
	expString = addRE.ReplaceAllString(expString, addRepl)

	printDiff(gotString, expString)
}
