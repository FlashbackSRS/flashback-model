package fb

import "testing"

type escapeFilenameTest struct {
	Filename string
	Expected string
}

func TestEscapeFilename(t *testing.T) {
	tests := []escapeFilenameTest{
		{
			Filename: "foobar.jpg",
			Expected: "foobar.jpg",
		},
		{
			Filename: "_foobar.jpg",
			Expected: "^_foobar.jpg",
		},
		{
			Filename: "^foobar.jpg",
			Expected: "^^foobar.jpg",
		},
		{
			Filename: "foo^bar_baz.jpg",
			Expected: "foo^bar_baz.jpg",
		},
		{
			Filename: "영상.jpg",
			Expected: "영상.jpg",
		},
	}
	for _, test := range tests {
		result := escapeFilename(test.Filename)
		if result != test.Expected {
			t.Errorf("Escape filename '%s' failed.\n\tExpected: %s\n\t  Actual: %s\n", test.Filename, test.Expected, result)
		}
		unResult, err := unescapeFilename(result)
		if err != nil {
			t.Errorf("Unexpected error unescaping filename '%s': %s\n", unResult, err)
		}
		if unResult != test.Filename {
			t.Errorf("Unescape filename '%s' failed.\n\tExpected: %s\n\t  Actual: %s\n", result, test.Filename, unResult)
		}
	}
}
