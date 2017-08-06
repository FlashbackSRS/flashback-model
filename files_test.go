package fb

import (
	"encoding/json"
	"testing"

	"github.com/flimzy/diff"
)

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
		unResult := unescapeFilename(result)
		if unResult != test.Filename {
			t.Errorf("Unescape filename '%s' failed.\n\tExpected: %s\n\t  Actual: %s\n", result, test.Filename, unResult)
		}
	}
}

func TestFilesUnmarshalJSON(t *testing.T) {
	type fujTest struct {
		name     string
		input    string
		expected interface{}
		err      string
	}
	tests := []fujTest{
		{
			name:  "bogus JSON",
			input: "invalid",
			err:   "invalid character 'i' looking for beginning of value",
		},
		{
			name: "valid",
			input: `{
			"^_weirdname.txt": {
				"content_type": "audio/mpeg",
				"data": "YSBm"
			},
			"영상.jpg": {
				"content_type": "audio/mpeg",
				"data": "YSBL"
			}
		}`,
			expected: &FileCollection{
				files: map[string]*Attachment{
					"_weirdname.txt": {ContentType: "audio/mpeg", Content: []byte{0x61, 0x20, 0x66}},
					"영상.jpg":         {ContentType: "audio/mpeg", Content: []byte{0x61, 0x20, 0x4b}},
				},
				views: []*FileCollectionView{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &FileCollection{}
			err := result.UnmarshalJSON([]byte(test.input))
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if d := diff.Interface(test.expected, result); d != "" {
				t.Error(d)
			}
		})
	}
}

func TestFileCollectionViewUnmarshalJSON(t *testing.T) {
	type Test struct {
		name     string
		input    string
		expected interface{}
		err      string
	}
	tests := []Test{
		{
			name:  "invalid json",
			input: "invalid json",
			err:   "failed to unmarshal file collection view: invalid character 'i' looking for beginning of value",
		},
		{
			name:  "valid",
			input: `["foo.txt","bar.mp3"]`,
			expected: &FileCollectionView{
				members: map[string]*Attachment{
					"foo.txt": nil,
					"bar.mp3": nil,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &FileCollectionView{}
			err := result.UnmarshalJSON([]byte(test.input))
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if d := diff.Interface(test.expected, result); d != "" {
				t.Error(d)
			}
		})
	}
}

func TestFCHasMemberView(t *testing.T) {
	t.Run("Member", func(t *testing.T) {
		att := NewFileCollection()
		view := att.NewView()
		if !att.hasMemberView(view) {
			t.Errorf("Expected success")
		}
	})
	t.Run("Non-member", func(t *testing.T) {
		att := NewFileCollection()
		view := NewFileCollection().NewView()
		if att.hasMemberView(view) {
			t.Errorf("Expected failure")
		}
	})
}

func TestAddFile(t *testing.T) {
	type Test struct {
		name     string
		view     *FileCollectionView
		filename string
		err      string
		expected interface{}
	}
	tests := []Test{
		{
			name:     "valid",
			view:     NewFileCollection().NewView(),
			filename: "foo.txt",
			expected: []string{"foo.txt"},
		},
		{
			name: "duplicate",
			view: func() *FileCollectionView {
				v := NewFileCollection().NewView()
				_ = v.AddFile("foo.txt", "text/plain", []byte("foo"))
				return v
			}(),
			filename: "foo.txt",
			err:      "'foo.txt' already exists in the collection",
		},
		{
			name:     "no file name",
			view:     NewFileCollection().NewView(),
			expected: []string{""},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.view.AddFile(test.filename, "text/plain", []byte("foo"))
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if d := diff.AsJSON(test.expected, test.view); d != "" {
				t.Error(d)
			}
		})
	}
}

func TestFileCollectionMarshalJSON(t *testing.T) {
	type Test struct {
		name     string
		fc       *FileCollection
		expected string
		err      string
	}
	tests := []Test{
		{
			name:     "empty collection",
			fc:       NewFileCollection(),
			expected: `{}`,
		},
		{
			name: "two files",
			fc: func() *FileCollection {
				fc := NewFileCollection()
				view := fc.NewView()
				_ = view.AddFile("abc.txt", "text/plain", []byte("abc"))
				_ = view.AddFile("123.txt", "text/plain", []byte("123"))
				return fc
			}(),
			expected: `{
				"123.txt": {"content_type":"text/plain", "data":"MTIz"},
				"abc.txt": {"content_type":"text/plain", "data":"YWJj"}
			}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := json.Marshal(test.fc)
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if d := diff.JSON([]byte(test.expected), result); d != "" {
				t.Error(d)
			}
		})
	}
}
