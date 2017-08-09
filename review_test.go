package fb

import (
	"testing"

	"github.com/flimzy/diff"
)

func TestNewReview(t *testing.T) {
	tests := []struct {
		name     string
		cardID   string
		expected *Review
		err      string
	}{
		{
			name: "validation fails",
			err:  "card id required",
		},
		{
			name:   "valid",
			cardID: "card-Zm9v",
			expected: &Review{
				CardID:    "card-Zm9v",
				Timestamp: now(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewReview(test.cardID)
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

func TestReviewValidate(t *testing.T) {
	tests := []validationTest{
		{
			name: "no card id",
			v:    &Review{},
			err:  "card id required",
		},
		{
			name: "invalid card id",
			v:    &Review{CardID: "oink"},
			err:  "invalid DocID format",
		},
		{
			name: "wrong card id type",
			v:    &Review{CardID: "note-Zm9v"},
			err:  "incorrect doc type for card ID",
		},
		{
			name: "no timestamp",
			v:    &Review{CardID: "card-Zm9v"},
			err:  "timestamp required",
		},
		{
			name: "valid",
			v:    &Review{CardID: "card-Zm9v", Timestamp: now()},
		},
	}
	testValidation(t, tests)
}
