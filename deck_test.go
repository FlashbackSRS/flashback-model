package fb

import (
	"testing"

	"github.com/flimzy/diff"
)

func TestCCMarshalJSON(t *testing.T) {
	type Test struct {
		name     string
		cc       *CardCollection
		expected string
		err      string
	}
	tests := []Test{
		{
			name:     "empty",
			cc:       &CardCollection{},
			expected: "[]",
		},
		{
			name: "some cards",
			cc: &CardCollection{
				col: map[string]struct{}{
					"card-foo": {},
					"card-bar": {},
				},
			},
			expected: `["card-bar","card-foo"]`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.cc.MarshalJSON()
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

func TestNewCardCollection(t *testing.T) {
	cc := NewCardCollection()
	expected := &CardCollection{
		col: map[string]struct{}{},
	}
	if d := diff.Interface(expected, cc); d != "" {
		t.Error(d)
	}
}

func TestCCUnmarshalJSON(t *testing.T) {
	type Test struct {
		name     string
		input    string
		expected *CardCollection
		err      string
	}
	tests := []Test{
		{
			name:  "invalid json",
			input: "invalid json",
			err:   "invalid character 'i' looking for beginning of value",
		},
		{
			name:  "valid",
			input: `["card-foo","card-bar"]`,
			expected: &CardCollection{col: map[string]struct{}{
				"card-foo": {},
				"card-bar": {},
			}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &CardCollection{}
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

func TestCCAll(t *testing.T) {
	cc := &CardCollection{col: map[string]struct{}{
		"card-foo": {},
		"card-bar": {},
	}}
	expected := []string{"card-bar", "card-foo"}
	result := cc.All()
	if d := diff.Interface(expected, result); d != "" {
		t.Error(d)
	}
}
