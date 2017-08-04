package fb

import (
	"testing"

	"github.com/flimzy/diff"
)

func TestParseID(t *testing.T) {
	type pidTest struct {
		name              string
		input             string
		bundle, note, err string
		template          uint32
	}
	tests := []pidTest{
		{
			name: "empty input",
			err:  "invalid ID type",
		},
		{
			name:  "invalid ID",
			input: "card-the quick brown fox",
			err:   "invalid ID format",
		},
		{
			name:  "invalid template",
			input: "card-krsxg5baij2w4zdmmu.mViuXQThMLoh1G1Nlc4d_E8kR8o.boo",
			err:   `invalid TemplateID: strconv.Atoi: parsing "boo": invalid syntax`,
		},
		{
			name:  "wrong id type",
			input: "foo-card-krsxg5baij2w4zdmmu.mViuXQThMLoh1G1Nlc4d_E8kR8o.0",
			err:   "invalid ID type",
		},
		{
			name:     "valid",
			input:    "card-krsxg5baij2w4zdmmu.mViuXQThMLoh1G1Nlc4d_E8kR8o.0",
			bundle:   "krsxg5baij2w4zdmmu",
			note:     "mViuXQThMLoh1G1Nlc4d_E8kR8o",
			template: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bundle, note, template, err := parseID(test.input)
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if bundle != test.bundle || note != test.note || template != test.template {
				t.Errorf("Unexpected result: %s %s %d", bundle, note, template)
			}
		})
	}
}

func TestNewCard(t *testing.T) {
	type ncTest struct {
		name     string
		theme    string
		model    uint32
		id       string
		expected *Card
		err      string
	}
	tests := []ncTest{
		{
			name: "invalid id",
			id:   "chicken man",
			err:  "error parsing card ID: invalid ID type",
		},
		{
			name:  "valid",
			theme: "theme-foo",
			id:    "card-krsxg5baij2w4zdmmu.mViuXQThMLoh1G1Nlc4d_E8kR8o.1",
			expected: &Card{
				bundleID:   "krsxg5baij2w4zdmmu",
				noteID:     "mViuXQThMLoh1G1Nlc4d_E8kR8o",
				templateID: 1,
				themeID:    "theme-foo",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewCard(test.theme, test.model, test.id)
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
