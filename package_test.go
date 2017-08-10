package fb

import "testing"

func TestPkgValidate(t *testing.T) {
	tests := []validationTest{
		{
			name: "card without deck",
			err:  "card 'card-abcde.mViuXQThMLoh1G1Nlc4d_E8kR8o.0' found in package, but not in a deck",
			v: &Package{
				Cards: []*Card{
					&Card{
						ID:      "card-abcde.mViuXQThMLoh1G1Nlc4d_E8kR8o.0",
						ModelID: "theme-VGVzdCBUaGVtZQ/0",
					},
				},
			},
		},
		{
			name: "card missing from package",
			err:  "card 'card-12345' listed in deck, but not found in package",
			v: &Package{
				Decks: []*Deck{
					&Deck{
						ID:    "deck-AQID",
						Cards: &CardCollection{map[string]struct{}{"card-12345": {}}},
					},
				},
			},
		},
		{
			name: "valid",
			v: &Package{
				Decks: []*Deck{
					&Deck{
						ID:    "deck-AQID",
						Cards: &CardCollection{map[string]struct{}{"card-abcde.mViuXQThMLoh1G1Nlc4d_E8kR8o.0": {}}},
					},
				},
				Cards: []*Card{
					&Card{
						ID:      "card-abcde.mViuXQThMLoh1G1Nlc4d_E8kR8o.0",
						ModelID: "theme-VGVzdCBUaGVtZQ/0",
					},
				},
			},
		},
	}
	testValidation(t, tests)
}
