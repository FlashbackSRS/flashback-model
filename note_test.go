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
