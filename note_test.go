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
			model: &Model{ID: 3, Theme: &Theme{ID: "theme-Zm9v"}},
			err:   "id required",
		},
		{
			name:  "valid",
			id:    "note-Zm9v",
			model: &Model{ID: 3, Theme: &Theme{ID: "theme-Zm9v"}},
			expected: func() *Note {
				att := NewFileCollection()
				return &Note{
					ID:          "note-Zm9v",
					ThemeID:     "theme-Zm9v",
					ModelID:     3,
					Created:     now(),
					Modified:    now(),
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
			result, err := NewNote(test.id, test.model)
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
				att := NewFileCollection()
				view := att.NewView()
				_ = view.AddFile("foo.txt", "text/plain", []byte("some text"))
				return &Note{
					ID:       "note-Zm9v",
					ThemeID:  "theme-Zm9v",
					ModelID:  3,
					Created:  now(),
					Modified: now(),
					Imported: now(),
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
				ID:       "note-Zm9v",
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
				att := NewFileCollection()
				view := att.NewView()
				_ = view.AddFile("foo.txt", "text/plain", []byte("some text"))
				return &Note{
					ID:       "note-Zm9v",
					ThemeID:  "theme-Zm9v",
					ModelID:  3,
					Created:  now(),
					Modified: now(),
					Imported: now(),
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

func TestFieldValueType(t *testing.T) {
	expected := AudioField
	fv := &FieldValue{field: &Field{Type: expected}}
	if ft := fv.Type(); ft != expected {
		t.Errorf("Unexpected result: %v", ft)
	}
}

func TestFieldValueUnmarshalJSON(t *testing.T) {
	type Test struct {
		name     string
		input    string
		expected *FieldValue
		err      string
	}
	tests := []Test{
		{
			name:  "invalid json",
			input: "invalid json",
			err:   "failed to unmarshal FieldValue: invalid character 'i' looking for beginning of value",
		},
		{
			name:     "empty field",
			input:    `{}`,
			expected: &FieldValue{},
		},
		{
			name:     "text field",
			input:    `{"text":"foo"}`,
			expected: &FieldValue{text: "foo"},
		},
		{
			name:  "files field",
			input: `{"text":"foo","files":["foo.txt","main.css"]}`,
			expected: &FieldValue{text: "foo", files: &FileCollectionView{
				members: map[string]*Attachment{"foo.txt": nil, "main.css": nil},
			}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &FieldValue{}
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

func TestFieldValueSetText(t *testing.T) {
	type Test struct {
		name     string
		fv       *FieldValue
		text     string
		err      string
		expected *FieldValue
	}
	tests := []Test{
		{
			name: "Audio field",
			fv:   &FieldValue{field: &Field{Type: AudioField}},
			text: "foo",
			err:  "Text field not permitted",
		},
		{
			name:     "Text field",
			fv:       &FieldValue{field: &Field{Type: TextField}},
			text:     "foo",
			expected: &FieldValue{field: &Field{Type: TextField}, text: "foo"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.fv.SetText(test.text)
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if d := diff.Interface(test.expected, test.fv); d != "" {
				t.Error(d)
			}
		})
	}
}

func TestFieldViewText(t *testing.T) {
	type Test struct {
		name     string
		fv       *FieldValue
		expected string
		err      string
	}
	tests := []Test{
		{
			name: "Audio field",
			fv:   &FieldValue{field: &Field{Type: AudioField}},
			err:  "FieldValue has no text field",
		},
		{
			name:     "Text field",
			fv:       &FieldValue{field: &Field{Type: TextField}, text: "foo"},
			expected: "foo",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.fv.Text()
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if test.expected != result {
				t.Errorf("Unexpected result: %s", result)
			}
		})
	}
}

func TestFieldViewAddFile(t *testing.T) {
	type Test struct {
		name     string
		fv       *FieldValue
		filename string
		expected *FieldValue
		err      string
	}
	tests := []Test{
		{
			name:     "text field",
			fv:       &FieldValue{field: &Field{Type: TextField}},
			filename: "foo.txt",
			err:      "Text fields do not support attachments",
		},
		{
			name:     "anki field",
			fv:       &FieldValue{field: &Field{Type: AnkiField}, files: NewFileCollection().NewView()},
			filename: "foo.txt",
			expected: func() *FieldValue {
				view := NewFileCollection().NewView()
				_ = view.AddFile("foo.txt", "text/plain", []byte("some text"))
				return &FieldValue{field: &Field{Type: AnkiField}, files: view}
			}(),
		},
		{
			name: "duplicate file",
			fv: func() *FieldValue {
				view := NewFileCollection().NewView()
				_ = view.AddFile("foo.txt", "text/plain", []byte("some text"))
				return &FieldValue{field: &Field{Type: AnkiField}, files: view}
			}(),
			filename: "foo.txt",
			err:      "'foo.txt' already exists in the collection",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.fv.AddFile(test.filename, "text/plain", []byte("some text"))
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if d := diff.Interface(test.expected, test.fv); d != "" {
				t.Error(d)
			}
		})
	}
}

func TestNoteSetRev(t *testing.T) {
	note := &Note{}
	rev := "1-xxx"
	note.SetRev(rev)
	if note.Rev != rev {
		t.Errorf("failed to set rev")
	}
}

func TestNoteDocID(t *testing.T) {
	note := &Note{ID: "note-Zm9v"}
	expected := "note-Zm9v"
	if id := note.DocID(); id != expected {
		t.Errorf("unexpected id: %s", id)
	}
}

func TestNoteImportedTime(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		note := &Note{}
		ts := now()
		note.Imported = ts
		if it := note.ImportedTime(); it != ts {
			t.Errorf("Unexpected result: %s", it)
		}
	})
	t.Run("Unset", func(t *testing.T) {
		note := &Note{}
		if it := note.ImportedTime(); !it.IsZero() {
			t.Errorf("unexpected result: %v", it)
		}
	})
}

func TestNoteModifiedTime(t *testing.T) {
	note := &Note{}
	ts := now()
	note.Modified = ts
	if mt := note.ModifiedTime(); mt != ts {
		t.Errorf("Unexpected result")
	}
}

func TestNoteMergeImport(t *testing.T) {
	type Test struct {
		name         string
		new          *Note
		existing     *Note
		expected     bool
		expectedNote *Note
		err          string
	}
	tests := []Test{
		{
			name:     "different ids",
			new:      &Note{ID: "note-Zm9v"},
			existing: &Note{ID: "note-YmFy"},
			err:      "IDs don't match",
		},
		{
			name:     "created timestamps don't match",
			new:      &Note{ID: "note-Zm9v", Created: parseTime("2017-01-01T01:01:01Z"), Imported: parseTime("2017-01-15T00:00:00Z")},
			existing: &Note{ID: "note-Zm9v", Created: parseTime("2017-02-01T01:01:01Z"), Imported: parseTime("2017-01-20T00:00:00Z")},
			err:      "Created timestamps don't match",
		},
		{
			name:     "new not an import",
			new:      &Note{ID: "note-Zm9v", Created: parseTime("2017-01-01T01:01:01Z")},
			existing: &Note{ID: "note-Zm9v", Created: parseTime("2017-01-01T01:01:01Z"), Imported: parseTime("2017-01-15T00:00:00Z")},
			err:      "not an import",
		},
		{
			name:     "existing not an import",
			new:      &Note{ID: "note-Zm9v", Created: parseTime("2017-01-01T01:01:01Z"), Imported: parseTime("2017-01-15T00:00:00Z")},
			existing: &Note{ID: "note-Zm9v", Created: parseTime("2017-01-01T01:01:01Z")},
			err:      "not an import",
		},
		{
			name: "new is newer",
			new: &Note{
				ID:          "note-Zm9v",
				ThemeID:     "theme-Zm9v",
				ModelID:     1,
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-02-01T01:01:01Z"),
				Imported:    parseTime("2017-01-15T00:00:00Z"),
				FieldValues: []*FieldValue{},
				Attachments: NewFileCollection(),
				model:       &Model{ID: 1},
			},
			existing: &Note{
				ID:          "note-Zm9v",
				ThemeID:     "theme-YmFy",
				ModelID:     2,
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-01-01T01:01:01Z"),
				Imported:    parseTime("2017-01-20T00:00:00Z"),
				FieldValues: []*FieldValue{{}},
				model:       &Model{ID: 2},
			},
			expected: true,
			expectedNote: &Note{
				ID:          "note-Zm9v",
				ThemeID:     "theme-Zm9v",
				ModelID:     1,
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-02-01T01:01:01Z"),
				Imported:    parseTime("2017-01-15T00:00:00Z"),
				FieldValues: []*FieldValue{},
				Attachments: NewFileCollection(),
				model:       &Model{ID: 1},
			},
		},
		{
			name: "existing is newer",
			new: &Note{
				ID:          "note-Zm9v",
				ThemeID:     "theme-Zm9v",
				ModelID:     1,
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-01-01T01:01:01Z"),
				Imported:    parseTime("2017-01-15T00:00:00Z"),
				FieldValues: []*FieldValue{},
				Attachments: NewFileCollection(),
				model:       &Model{ID: 1},
			},
			existing: &Note{
				ID:          "note-Zm9v",
				ThemeID:     "theme-Zm9v",
				ModelID:     2,
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-02-01T01:01:01Z"),
				Imported:    parseTime("2017-01-20T00:00:00Z"),
				FieldValues: []*FieldValue{{}},
				model:       &Model{ID: 2},
			},
			expected: false,
			expectedNote: &Note{ID: "note-Zm9v",
				ThemeID:     "theme-Zm9v",
				ModelID:     2,
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-02-01T01:01:01Z"),
				Imported:    parseTime("2017-01-20T00:00:00Z"),
				FieldValues: []*FieldValue{{}},
				model:       &Model{ID: 2},
			},
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
			if d := diff.Interface(test.expectedNote, test.new); d != "" {
				t.Error(d)
			}
		})
	}
}

func TestNoteValidate(t *testing.T) {
	tests := []validationTest{
		{
			name: "no ID",
			v:    &Note{},
			err:  "id required",
		},
		{
			name: "invalid doctype",
			v:    &Note{ID: "chicken-foo"},
			err:  "incorrect doc type",
		},
		{
			name: "wrong doctype",
			v:    &Note{ID: "deck-foo"},
			err:  "incorrect doc type",
		},
		{
			name: "no created time",
			v:    &Note{ID: "note-Zm9v"},
			err:  "created time required",
		},
		{
			name: "no modified time",
			v:    &Note{ID: "note-Zm9v", Created: now()},
			err:  "modified time required",
		},
		{
			name: "nil attachments collection",
			v:    &Note{ID: "note-Zm9v", Created: now(), Modified: now()},
			err:  "attachments collection must not be nil",
		},
		{
			name: "invalid field file list",
			v:    &Note{ID: "note-Zm9v", Created: now(), Modified: now(), Attachments: NewFileCollection(), FieldValues: []*FieldValue{{files: NewFileCollection().NewView()}}},
			err:  "field 0 file list must be member of attachments collection",
		},
		{
			name: "valid",
			v: func() *Note {
				att := NewFileCollection()
				view := att.NewView()
				return &Note{ID: "note-Zm9v", Created: now(), Modified: now(), Attachments: att, FieldValues: []*FieldValue{{files: view}}}
			}(),
		},
	}
	testValidation(t, tests)
}
