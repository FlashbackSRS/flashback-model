package fb

import (
	"testing"

	"github.com/flimzy/diff"
)

func TestNewNote(t *testing.T) {
	type Test struct {
		name     string
		id       string
		model    *Model
		expected *Note
		err      string
	}
	tests := []Test{
		{
			name: "no model",
			id:   "foo",
			err:  "model required",
		},
		{
			name:  "no id",
			model: &Model{},
			err:   "id is required",
		},
		{
			name:  "valid",
			id:    "foo",
			model: &Model{ID: 3, Theme: &Theme{ID: "theme-Zm9v"}},
			expected: func() *Note {
				att := NewFileCollection()
				return &Note{
					ID:          DocID{docType: "note", id: []byte("foo")},
					ThemeID:     "theme-Zm9v",
					ModelID:     3,
					FieldValues: []*FieldValue{},
					Attachments: att,
					model: &Model{
						ID: 3,
						Theme: &Theme{
							ID: "theme-Zm9v",
						},
					},
				}
			}(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewNote([]byte(test.id), test.model)
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

func TestNoteSetModel(t *testing.T) {
	type Test struct {
		name     string
		note     *Note
		model    *Model
		expected *Note
		err      string
	}
	tests := []Test{
		{
			name: "nil model",
			note: &Note{},
			err:  "model required",
		},
		{
			name: "Theme IDs don't match",
			note: &Note{ThemeID: "theme-Zm9v"},
			model: &Model{
				Theme: &Theme{ID: "theme-YmFy"},
			},
			err: "Theme IDs must match",
		},
		{
			name: "fields and field values don't match",
			note: &Note{
				ThemeID:     "theme-Zm9v",
				FieldValues: []*FieldValue{{}, {}},
			},
			model: &Model{
				Theme: &Theme{ID: "theme-Zm9v"},
			},
			err: "model.Fields and node.FieldValues lengths must match",
		},
		{
			name: "no fields",
			note: &Note{ThemeID: "theme-Zm9v"},
			model: &Model{
				Theme: &Theme{ID: "theme-Zm9v"},
			},
			expected: &Note{
				ThemeID: "theme-Zm9v",
				model: &Model{
					Theme: &Theme{ID: "theme-Zm9v"},
				},
			},
		},
		{
			name: "with fields",
			note: &Note{
				ThemeID: "theme-Zm9v",
				FieldValues: []*FieldValue{
					{text: "one"},
					{text: "two"},
				},
			},
			model: &Model{
				Theme: &Theme{ID: "theme-Zm9v"},
				Fields: []*Field{
					{Name: "foo"},
					{Name: "bar"},
				},
			},
			expected: &Note{
				ThemeID: "theme-Zm9v",
				FieldValues: []*FieldValue{
					{text: "one", field: &Field{Name: "foo"}},
					{text: "two", field: &Field{Name: "bar"}},
				},
				model: &Model{
					Theme: &Theme{ID: "theme-Zm9v"},
					Fields: []*Field{
						{Name: "foo"},
						{Name: "bar"},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.note.SetModel(test.model)
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if d := diff.Interface(test.expected, test.note); d != "" {
				t.Error(d)
			}
		})
	}
}

func TestNoteModel(t *testing.T) {
	m := &Model{}
	n := &Note{model: m}
	if n.Model() != m {
		t.Error("Unexpected result")
	}
}

func TestNoteMarshalJSON(t *testing.T) {
	type Test struct {
		name     string
		note     *Note
		expected string
		err      string
	}
	tests := []Test{
		{
			name: "all fields",
			note: func() *Note {
				nowTime := now()
				att := NewFileCollection()
				view := att.NewView()
				_ = view.AddFile("foo.txt", "text/plain", []byte("some text"))
				return &Note{
					ID:       DocID{docType: "note", id: []byte("foo")},
					ThemeID:  "theme-Zm9v",
					ModelID:  3,
					Created:  now(),
					Modified: now(),
					Imported: &nowTime,
					FieldValues: []*FieldValue{
						{text: "foo", files: view},
					},
					Attachments: att,
				}
			}(),
			expected: `{
                "_id":          "note-Zm9v",
                "type":         "note",
                "created":      "2017-01-01T00:00:00Z",
                "modified":     "2017-01-01T00:00:00Z",
                "imported":     "2017-01-01T00:00:00Z",
                "fieldValues":  [{"text":"foo", "files":["foo.txt"]}],
                "model":        3,
                "theme":        "theme-Zm9v",
                "_attachments": {
                    "foo.txt": {"content_type":"text/plain", "data":"c29tZSB0ZXh0"}
                }
            }`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.note.MarshalJSON()
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

func TestNoteUnmarshalJSON(t *testing.T) {
	type Test struct {
		name     string
		input    string
		expected *Note
		err      string
	}
	tests := []Test{
		{
			name:  "invalid json",
			input: "invalid json",
			err:   "failed to unmarshal Note: invalid character 'i' looking for beginning of value",
		},
		{
			name:  "invalid type",
			input: `{"type":"chicken"}`,
			err:   "Invalid document type for note: chicken",
		},
		{
			name:  "wrong type",
			input: `{"type":"theme"}`,
			err:   "Invalid document type for note: theme",
		},
		{
			name: "null fields",
			input: `{
                "_id":          "note-Zm9v",
                "type":         "note",
                "created":      "2017-01-01T00:00:00Z",
                "modified":     "2017-01-01T00:00:00Z",
                "model":        3,
                "theme":        "theme-Zm9v"
            }`,
			expected: &Note{
				ID:       DocID{docType: "note", id: []byte("foo")},
				Created:  now(),
				Modified: now(),
				ModelID:  3,
				ThemeID:  "theme-Zm9v",
			},
		},
		{
			name: "all fields",
			input: `{
                "_id":          "note-Zm9v",
                "type":         "note",
                "created":      "2017-01-01T00:00:00Z",
                "modified":     "2017-01-01T00:00:00Z",
                "imported":     "2017-01-01T00:00:00Z",
                "fieldValues":  [{"text":"foo", "files":["foo.txt"]}],
                "model":        3,
                "theme":        "theme-Zm9v",
                "_attachments": {
                    "foo.txt": {"content_type":"text/plain", "data":"c29tZSB0ZXh0"}
                }
            }`,
			expected: func() *Note {
				nowTime := now()
				att := NewFileCollection()
				view := att.NewView()
				_ = view.AddFile("foo.txt", "text/plain", []byte("some text"))
				return &Note{
					ID:       DocID{docType: "note", id: []byte("foo")},
					ThemeID:  "theme-Zm9v",
					ModelID:  3,
					Created:  now(),
					Modified: now(),
					Imported: &nowTime,
					FieldValues: []*FieldValue{
						{text: "foo", files: view},
					},
					Attachments: att,
				}
			}(),
		},
		{
			name: "invalid file view",
			input: `{
                "_id":          "note-Zm9v",
                "type":         "note",
                "created":      "2017-01-01T00:00:00Z",
                "modified":     "2017-01-01T00:00:00Z",
                "imported":     "2017-01-01T00:00:00Z",
                "fieldValues":  [{"text":"foo", "files":["foo.html"]}],
                "model":        3,
                "theme":        "theme-Zm9v",
                "_attachments": {
                    "foo.txt": {"content_type":"text/plain", "data":"c29tZSB0ZXh0"}
                }
            }`,
			err: "foo.html not found in collection",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &Note{}
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

func TestNoteGetFieldValue(t *testing.T) {
	type Test struct {
		name     string
		note     *Note
		ord      int
		expected *FieldValue
	}
	tests := []Test{
		{
			name: "new text field",
			note: &Note{
				FieldValues: make([]*FieldValue, 1),
				model:       &Model{Fields: []*Field{{Type: TextField, Name: "text"}}},
			},
			ord: 0,
			expected: &FieldValue{
				field: &Field{
					Type: TextField,
					Name: "text",
				},
			},
		},
		{
			name: "new audio field",
			note: &Note{
				FieldValues: make([]*FieldValue, 1),
				model:       &Model{Fields: []*Field{{Type: AudioField, Name: "text"}}},
				Attachments: NewFileCollection(),
			},
			ord: 0,
			expected: func() *FieldValue {
				return &FieldValue{
					field: &Field{
						Type: AudioField,
						Name: "text",
					},
					files: NewFileCollection().NewView(),
				}
			}(),
		},
		{
			name: "existing field",
			note: &Note{
				FieldValues: []*FieldValue{{
					field: &Field{Type: TextField, Name: "foo"},
					text:  "foo text",
				}},
			},
			ord: 0,
			expected: &FieldValue{
				field: &Field{Type: TextField, Name: "foo"},
				text:  "foo text",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.note.GetFieldValue(test.ord)
			if d := diff.Interface(test.expected, result); d != "" {
				t.Error(d)
			}
		})
	}
}
