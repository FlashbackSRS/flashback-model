package fb

import (
	"testing"

	"github.com/flimzy/diff"
)

func TestNewModel(t *testing.T) {
	type Test struct {
		name      string
		theme     *Theme
		modelType string
		expected  *Model
		err       string
	}
	tests := []Test{
		{
			name: "nil theme",
			err:  "theme is required",
		},
		{
			name:  "no type",
			theme: &Theme{},
			err:   "model type is required",
		},
		{
			name: "valid",
			theme: func() *Theme {
				theme, _ := NewTheme([]byte("foo"))
				return theme
			}(),
			modelType: "foo",
			expected: func() *Model {
				theme, _ := NewTheme([]byte("foo"))
				// att := NewFileCollection()
				// theme.Files = att.NewView()
				theme.modelSequence = 1
				model := &Model{
					Type:      "foo",
					Templates: []string{},
					Fields:    []*Field{},
					Files:     theme.Attachments.NewView(),
					Theme:     theme,
				}
				return model
			}(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewModel(test.theme, test.modelType)
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
