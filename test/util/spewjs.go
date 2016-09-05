// +build js,!disableunsafe

package util

func init() {
	panic("Must run with 'disableunsafe' build flag")
}

func PrintDiff(got, expected interface{}) {}
