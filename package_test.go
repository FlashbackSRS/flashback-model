package fb

import "testing"

func TestPkgValidate(t *testing.T) {
	type pvTest struct {
		name string
		pkg  *Package
		err  string
	}
	tests := []pvTest{
		{
			name: "card without deck",
			err:  "card 'abcde.mViuXQThMLoh1G1Nlc4d_E8kR8o.0' found in package, but not in a deck",
			pkg: &Package{
				Cards: []*Card{
					func() *Card {
						c, err := NewCard("theme-VGVzdCBUaGVtZQ", 0, "card-abcde.mViuXQThMLoh1G1Nlc4d_E8kR8o.0")
						if err != nil {
							t.Fatal(err)
						}
						return c
					}(),
				},
			},
		},
		{
			name: "card missing from package",
			err:  "card 'card-12345' listed in deck, but not found in package",
			pkg: &Package{
				Decks: []*Deck{
					func() *Deck {
						id, _ := NewDocID("deck", []byte{1, 2, 3})
						return &Deck{
							ID:    id,
							Cards: &CardCollection{map[string]struct{}{"card-12345": struct{}{}}},
						}
					}(),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var msg string
			if err := test.pkg.Validate(); err != nil {
				msg = err.Error()
			}
			if test.err != msg {
				t.Errorf("Unexpected error: %s", msg)
			}
		})
	}
}