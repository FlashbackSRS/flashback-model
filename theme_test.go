package fb

import (
	"encoding/json"
	"testing"

	"github.com/flimzy/diff"
)

func TestNewTheme(t *testing.T) {
	type Test struct {
		name     string
		id       []byte
		expected interface{}
		err      string
	}
	tests := []Test{
		{
			name: "no id",
			err:  "failed to create DocID for Theme: id is required",
		},
		{
			name: "valid",
			id:   []byte("theme id"),
			expected: func() *Theme {
				t := &Theme{
					ID:          DocID{docType: "theme", id: []byte("theme id")},
					Models:      make([]*Model, 0, 1),
					Attachments: NewFileCollection(),
				}
				t.Files = t.Attachments.NewView()
				return t
			}(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewTheme(test.id)
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

func TestSetFile(t *testing.T) {
	att := NewFileCollection()
	view := att.NewView()
	_ = view.AddFile("foo.mp3", "audio/mpeg", []byte("foo"))
	theme, _ := NewTheme([]byte("foo"))
	theme.SetFile("foo.mp3", "audio/mpeg", []byte("foo"))
	expected := &Theme{
		ID:          DocID{docType: "theme", id: []byte("foo")},
		Models:      []*Model{},
		Attachments: att,
		Files:       view,
	}
	if d := diff.Interface(expected, theme); d != "" {
		t.Error(d)
	}
}

func TestThemeMarshalJSON(t *testing.T) {
	type Test struct {
		name     string
		theme    *Theme
		expected string
		err      string
	}
	tests := []Test{
		{
			name: "valid",
			theme: func() *Theme {
				theme, _ := NewTheme([]byte("test theme id"))
				theme.SetFile("file.txt", "text/plain", []byte("some text"))
				theme.Created = now()
				theme.Modified = now()
				return theme
			}(),
			expected: `{
                "_id":           "theme-dGVzdCB0aGVtZSBpZA",
                "type":          "theme",
                "created":       "2017-01-01T00:00:00Z",
                "modified":      "2017-01-01T00:00:00Z",
                "modelSequence": 0,
                "files":         ["file.txt"],
                "_attachments":  {
                    "file.txt": {
                        "content_type": "text/plain",
                        "data":         "c29tZSB0ZXh0"
                    }
                }
            }`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := json.Marshal(test.theme)
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

func TestThemeUnmarshalJSON(t *testing.T) {
	type Test struct {
		name     string
		input    string
		expected interface{}
		err      string
	}
	tests := []Test{
		{
			name:  "invalid json",
			input: "xxx",
			err:   "failed to unmarshal Theme: invalid character 'x' looking for beginning of value",
		},
		{
			name:  "wrong type",
			input: `{"type":"chicken"}`,
			err:   "Invalid document type for theme: chicken",
		},
		{
			name:  "no attachments",
			input: `{"type":"theme", "_id":"theme-120","created":"2017-01-01T01:01:01Z","modified":"2017-01-01T01:01:01Z"}`,
			err:   "invalid theme: no attachments",
		},
		{
			name: "with attachments",
			input: `{"type":"theme", "_id":"theme-120", "created":"2017-01-01T01:01:01Z", "modified":"2017-01-01T01:01:01Z", "_attachments":{
            "foo.txt": {"content_type":"text/plain", "content": "text"}
            }}`,
			err: "invalid theme: no file list",
		},
		{
			name:  "mismatched file list",
			input: `{"type":"theme", "_id":"theme-120", "created":"2017-01-01T01:01:01Z", "modified":"2017-01-01T01:01:01Z", "_attachments": {"foo.txt": {"content_type":"text/plain", "content": "text"}}, "files": ["foo.html"] }`,
			err:   "foo.html not found in collection",
		},
		{
			name:  "mismatched model file list",
			input: `{"type":"theme", "_id":"theme-120", "created":"2017-01-01T01:01:01Z", "modified":"2017-01-01T01:01:01Z", "_attachments": {"foo.txt": {"content_type":"text/plain", "content": "text"}}, "files":[], "models": [{"id":0, "files": ["foo.mp3"]}] }`,
			err:   "foo.mp3 not found in collection",
		},
		{
			name:  "valid",
			input: `{"type":"theme", "_id":"theme-120", "created":"2017-01-01T01:01:01Z", "modified":"2017-01-01T01:01:01Z", "_attachments": {"foo.txt": {"content_type":"text/plain", "data": "text"}}, "files":[], "models": [{"id":0, "files": ["foo.txt"]}] }`,
			expected: map[string]interface{}{
				"type":     "theme",
				"_id":      "theme-120",
				"created":  "2017-01-01T01:01:01Z",
				"modified": "2017-01-01T01:01:01Z",
				"models": []map[string]interface{}{
					{
						"id":        0,
						"modelType": "",
						"templates": nil,
						"fields":    nil,
						"files":     []string{"foo.txt"},
					},
				},
				"_attachments": map[string]interface{}{
					"foo.txt": map[string]string{
						"content_type": "text/plain",
						"data":         "text",
					},
				},
				"files":         []string{},
				"modelSequence": 0,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &Theme{}
			err := result.UnmarshalJSON([]byte(test.input))
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if d := diff.AsJSON(test.expected, result); d != "" {
				t.Error(d)
			}
		})
	}
}

func TestThemeNewModel(t *testing.T) {
	type Test struct {
		name      string
		theme     *Theme
		modelType string
		expected  *Model
		err       string
	}
	tests := []Test{
		{
			name:  "no type",
			theme: &Theme{},
			err:   "failed to create model: model type is required",
		},
		{
			name: "success",
			theme: func() *Theme {
				theme, _ := NewTheme([]byte("foo"))
				return theme
			}(),
			modelType: "chicken",
			expected: func() *Model {
				theme, _ := NewTheme([]byte("foo"))
				// att := NewFileCollection()
				// theme.Files = att.NewView()
				theme.modelSequence = 1
				model := &Model{
					Type:      "chicken",
					Templates: []string{},
					Fields:    []*Field{},
					Files:     theme.Attachments.NewView(),
					Theme:     theme,
				}
				theme.Models = []*Model{model}
				return model
			}(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.theme.NewModel(test.modelType)
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

func TestThemeSetRev(t *testing.T) {
	theme := &Theme{}
	rev := "1-xxx"
	theme.SetRev(rev)
	if *theme.Rev != rev {
		t.Errorf("failed to set rev")
	}
}

func TestThemeDocID(t *testing.T) {
	theme, _ := NewTheme([]byte("foo"))
	expected := "theme-Zm9v"
	if id := theme.DocID(); id != expected {
		t.Errorf("unexpected id: %s", id)
	}
}

func TestThemeImportedTime(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		theme := &Theme{}
		ts := now()
		theme.Imported = &ts
		if it := theme.ImportedTime(); *it != ts {
			t.Errorf("Unexpected result")
		}
	})
	t.Run("Unset", func(t *testing.T) {
		theme := &Theme{}
		if it := theme.ImportedTime(); it != nil {
			t.Errorf("unexpected result")
		}
	})
}

func TestThemeModifiedTime(t *testing.T) {
	theme := &Theme{}
	ts := now()
	theme.Modified = ts
	if mt := theme.ModifiedTime(); *mt != ts {
		t.Errorf("Unexpected result")
	}
}

func TestThemeMergeImport(t *testing.T) {
	type Test struct {
		name          string
		new           *Theme
		existing      *Theme
		expected      bool
		expectedTheme *Theme
		err           string
	}
	tests := []Test{
		{
			name:     "different ids",
			new:      &Theme{ID: DocID{docType: "theme", id: []byte("a")}},
			existing: &Theme{ID: DocID{docType: "theme", id: []byte("b")}},
			err:      "IDs don't match",
		},
		{
			name:     "created timestamps don't match",
			new:      &Theme{ID: DocID{docType: "theme", id: []byte("a")}, Created: parseTime("2017-01-01T01:01:01Z"), Imported: parseTimePtr("2017-01-15T00:00:00Z")},
			existing: &Theme{ID: DocID{docType: "theme", id: []byte("a")}, Created: parseTime("2017-02-01T01:01:01Z"), Imported: parseTimePtr("2017-01-20T00:00:00Z")},
			err:      "Created timestamps don't match",
		},
		{
			name:     "new not an import",
			new:      &Theme{ID: DocID{docType: "theme", id: []byte("a")}, Created: parseTime("2017-01-01T01:01:01Z")},
			existing: &Theme{ID: DocID{docType: "theme", id: []byte("a")}, Created: parseTime("2017-02-01T01:01:01Z"), Imported: parseTimePtr("2017-01-15T00:00:00Z")},
			err:      "not an import",
		},
		{
			name:     "existing not an import",
			new:      &Theme{ID: DocID{docType: "theme", id: []byte("a")}, Created: parseTime("2017-01-01T01:01:01Z"), Imported: parseTimePtr("2017-01-15T00:00:00Z")},
			existing: &Theme{ID: DocID{docType: "theme", id: []byte("a")}, Created: parseTime("2017-02-01T01:01:01Z")},
			err:      "not an import",
		},
		{
			name: "new is newer",
			new: &Theme{ID: DocID{docType: "theme", id: []byte("a")},
				Name:     func() *string { x := "foo"; return &x }(),
				Created:  parseTime("2017-01-01T01:01:01Z"),
				Modified: parseTime("2017-02-01T01:01:01Z"),
				Imported: parseTimePtr("2017-01-15T00:00:00Z")},
			existing: &Theme{ID: DocID{docType: "theme", id: []byte("a")},
				Name:     func() *string { x := "bar"; return &x }(),
				Created:  parseTime("2017-01-01T01:01:01Z"),
				Modified: parseTime("2017-01-01T01:01:01Z"),
				Imported: parseTimePtr("2017-01-20T00:00:00Z")},
			expected: true,
			expectedTheme: &Theme{ID: DocID{docType: "theme", id: []byte("a")},
				Name:     func() *string { x := "foo"; return &x }(),
				Created:  parseTime("2017-01-01T01:01:01Z"),
				Modified: parseTime("2017-02-01T01:01:01Z"),
				Imported: parseTimePtr("2017-01-15T00:00:00Z")},
		},
		{
			name: "existing is newer",
			new: &Theme{ID: DocID{docType: "theme", id: []byte("a")},
				Name:     func() *string { x := "foo"; return &x }(),
				Created:  parseTime("2017-01-01T01:01:01Z"),
				Modified: parseTime("2017-01-01T01:01:01Z"),
				Imported: parseTimePtr("2017-01-15T00:00:00Z")},
			existing: &Theme{ID: DocID{docType: "theme", id: []byte("a")},
				Name:     func() *string { x := "bar"; return &x }(),
				Created:  parseTime("2017-01-01T01:01:01Z"),
				Modified: parseTime("2017-02-01T01:01:01Z"),
				Imported: parseTimePtr("2017-01-20T00:00:00Z")},
			expected: false,
			expectedTheme: &Theme{ID: DocID{docType: "theme", id: []byte("a")},
				Name:     func() *string { x := "bar"; return &x }(),
				Created:  parseTime("2017-01-01T01:01:01Z"),
				Modified: parseTime("2017-02-01T01:01:01Z"),
				Imported: parseTimePtr("2017-01-20T00:00:00Z")},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.new.MergeImport(test.existing)
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if test.expected != result {
				t.Errorf("Unexpected result: %t", result)
			}
			if d := diff.Interface(test.expectedTheme, test.new); d != "" {
				t.Error(d)
			}
		})
	}
}
