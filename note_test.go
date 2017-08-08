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
