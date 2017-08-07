package fb

import (
	"testing"
	"time"

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

func TestNewDeck(t *testing.T) {
	type Test struct {
		name     string
		id       string
		expected *Deck
		err      string
	}
	tests := []Test{
		{
			name: "no id",
			err:  "id is required",
		},
		{
			name: "valid",
			id:   "foo id",
			expected: &Deck{
				ID:    DocID{docType: "deck", id: []byte("foo id")},
				Cards: &CardCollection{col: map[string]struct{}{}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewDeck([]byte(test.id))
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

func TestDeckMarshalJSON(t *testing.T) {
	type Test struct {
		name     string
		deck     *Deck
		expected string
		err      string
	}
	tests := []Test{
		{
			name: "full fields",
			deck: &Deck{
				ID:          DocID{docType: "deck", id: []byte("deck")},
				Created:     now(),
				Modified:    now(),
				Imported:    func() *time.Time { x := now(); return &x }(),
				Name:        func() *string { x := "test name"; return &x }(),
				Description: func() *string { x := "test description"; return &x }(),
				Cards: &CardCollection{col: map[string]struct{}{
					"card-foo": {}, "card-bar": {},
				}},
			},
			expected: `{
            "type":        "deck",
            "_id":         "deck-ZGVjaw",
            "name":        "test name",
            "description": "test description",
            "created":     "2017-01-01T00:00:00Z",
            "modified":    "2017-01-01T00:00:00Z",
            "imported":    "2017-01-01T00:00:00Z",
            "cards":       ["card-bar","card-foo"]
            }`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.deck.MarshalJSON()
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

func TestDeckAddCard(t *testing.T) {
	deck, err := NewDeck([]byte("foo"))
	if err != nil {
		t.Fatal(err)
	}
	deck.AddCard("card-jack")
	deck.AddCard("card-jill")
	expected := &Deck{
		ID: DocID{docType: "deck", id: []byte("foo")},
		Cards: &CardCollection{col: map[string]struct{}{
			"card-jack": {},
			"card-jill": {},
		}},
	}
	if d := diff.Interface(expected, deck); d != "" {
		t.Error(d)
	}
}

func TestDeckUnmarshalJSON(t *testing.T) {
	type Test struct {
		name     string
		input    string
		expected *Deck
		err      string
	}
	tests := []Test{
		{
			name:  "invalid json",
			input: "invalid json",
			err:   "invalid character 'i' looking for beginning of value",
		},
		{
			name:  "wrong type",
			input: `{"type":"chicken"}`,
			err:   "Invalid document type for deck: chicken",
		},
		{
			name: "all fields",
			input: `{
                "type":        "deck",
                "_id":         "deck-ZGVjaw",
                "name":        "test name",
                "description": "test description",
                "created":     "2017-01-01T00:00:00Z",
                "modified":    "2017-01-01T00:00:00Z",
                "imported":    "2017-01-01T00:00:00Z",
                "cards":       ["card-bar","card-foo"]
            }`,
			expected: &Deck{
				ID:          DocID{docType: "deck", id: []byte("deck")},
				Created:     now(),
				Modified:    now(),
				Imported:    func() *time.Time { x := now(); return &x }(),
				Name:        func() *string { x := "test name"; return &x }(),
				Description: func() *string { x := "test description"; return &x }(),
				Cards:       &CardCollection{col: map[string]struct{}{"card-bar": {}, "card-foo": {}}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &Deck{}
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

func TestDeckSetRev(t *testing.T) {
	deck := &Deck{}
	rev := "1-xxx"
	deck.SetRev(rev)
	if *deck.Rev != rev {
		t.Errorf("Unexpected result: %s", *deck.Rev)
	}
}

func TestDeckDocID(t *testing.T) {
	deck := &Deck{ID: DocID{docType: "deck", id: []byte("foo")}}
	expected := "deck-Zm9v"
	if id := deck.DocID(); id != expected {
		t.Errorf("Unexpected result: %s", id)
	}
}

func TestDeckImportedTime(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		deck := &Deck{}
		ts := now()
		deck.Imported = &ts
		if it := deck.ImportedTime(); it != ts {
			t.Errorf("Unexpected result: %s", it)
		}
	})
	t.Run("Unset", func(t *testing.T) {
		deck := &Deck{}
		if it := deck.ImportedTime(); !it.IsZero() {
			t.Errorf("unexpected result: %v", it)
		}
	})
}

func TestDeckModifiedTime(t *testing.T) {
	deck := &Deck{}
	ts := now()
	deck.Modified = ts
	if mt := deck.ModifiedTime(); mt != ts {
		t.Errorf("Unexpected result")
	}
}

func TestDeckMergeImport(t *testing.T) {
	type Test struct {
		name         string
		new          *Deck
		existing     *Deck
		expected     bool
		expectedDeck *Deck
		err          string
	}
	tests := []Test{
		{
			name:     "different ids",
			new:      &Deck{ID: DocID{docType: "deck", id: []byte("abcd")}},
			existing: &Deck{ID: DocID{docType: "deck", id: []byte("b")}},
			err:      "IDs don't match",
		},
		{
			name:     "created timestamps don't match",
			new:      &Deck{ID: DocID{docType: "deck", id: []byte("abcd")}, Created: parseTime("2017-01-01T01:01:01Z"), Imported: parseTimePtr("2017-01-15T00:00:00Z")},
			existing: &Deck{ID: DocID{docType: "deck", id: []byte("abcd")}, Created: parseTime("2017-02-01T01:01:01Z"), Imported: parseTimePtr("2017-01-20T00:00:00Z")},
			err:      "Created timestamps don't match",
		},
		{
			name:     "new not an import",
			new:      &Deck{ID: DocID{docType: "deck", id: []byte("abcd")}, Created: parseTime("2017-01-01T01:01:01Z")},
			existing: &Deck{ID: DocID{docType: "deck", id: []byte("abcd")}, Created: parseTime("2017-01-01T01:01:01Z"), Imported: parseTimePtr("2017-01-15T00:00:00Z")},
			err:      "not an import",
		},
		{
			name:     "existing not an import",
			new:      &Deck{ID: DocID{docType: "deck", id: []byte("abcd")}, Created: parseTime("2017-01-01T01:01:01Z"), Imported: parseTimePtr("2017-01-15T00:00:00Z")},
			existing: &Deck{ID: DocID{docType: "deck", id: []byte("abcd")}, Created: parseTime("2017-01-01T01:01:01Z")},
			err:      "not an import",
		},
		{
			name: "new is newer",
			new: &Deck{
				ID:          DocID{docType: "deck", id: []byte("abcd")},
				Name:        func() *string { x := "foo"; return &x }(),
				Description: func() *string { x := "FOO"; return &x }(),
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-02-01T01:01:01Z"),
				Imported:    parseTimePtr("2017-01-15T00:00:00Z"),
				Cards:       &CardCollection{col: map[string]struct{}{"card-foo": {}}},
			},
			existing: &Deck{
				ID:          DocID{docType: "deck", id: []byte("abcd")},
				Name:        func() *string { x := "bar"; return &x }(),
				Description: func() *string { x := "BAR"; return &x }(),
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-01-01T01:01:01Z"),
				Imported:    parseTimePtr("2017-01-20T00:00:00Z"),
				Cards:       &CardCollection{col: map[string]struct{}{"card-bar": {}}},
			},
			expected: true,
			expectedDeck: &Deck{
				ID:          DocID{docType: "deck", id: []byte("abcd")},
				Name:        func() *string { x := "foo"; return &x }(),
				Description: func() *string { x := "FOO"; return &x }(),
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-02-01T01:01:01Z"),
				Imported:    parseTimePtr("2017-01-15T00:00:00Z"),
				Cards:       &CardCollection{col: map[string]struct{}{"card-foo": {}}},
			},
		},
		{
			name: "existing is newer",
			new: &Deck{
				ID:          DocID{docType: "deck", id: []byte("abcd")},
				Name:        func() *string { x := "foo"; return &x }(),
				Description: func() *string { x := "FOO"; return &x }(),
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-01-01T01:01:01Z"),
				Imported:    parseTimePtr("2017-01-15T00:00:00Z"),
				Cards:       &CardCollection{col: map[string]struct{}{"card-foo": {}}},
			},
			existing: &Deck{
				ID:          DocID{docType: "deck", id: []byte("abcd")},
				Name:        func() *string { x := "bar"; return &x }(),
				Description: func() *string { x := "BAR"; return &x }(),
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-02-01T01:01:01Z"),
				Imported:    parseTimePtr("2017-01-20T00:00:00Z"),
				Cards:       &CardCollection{col: map[string]struct{}{"card-bar": {}}},
			},
			expected: false,
			expectedDeck: &Deck{
				ID:          DocID{docType: "deck", id: []byte("abcd")},
				Name:        func() *string { x := "bar"; return &x }(),
				Description: func() *string { x := "BAR"; return &x }(),
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-02-01T01:01:01Z"),
				Imported:    parseTimePtr("2017-01-20T00:00:00Z"),
				Cards:       &CardCollection{col: map[string]struct{}{"card-bar": {}}},
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
			if d := diff.Interface(test.expectedDeck, test.new); d != "" {
				t.Error(d)
			}
		})
	}
}

func TestCCValidate(t *testing.T) {
	tests := []validationTest{
		{
			name: "no cards",
			v:    &CardCollection{},
		},
		{
			name: "invalid card ID",
			v:    &CardCollection{col: map[string]struct{}{"foo": {}}},
			err:  "'foo': invalid ID type",
		},
	}
	testValidation(t, tests)
}

func TestDeckValidate(t *testing.T) {
	tests := []validationTest{
		{
			name: "no id",
			v:    &deckDoc{},
			err:  "id required",
		},
		{
			name: "invalid doctype",
			v:    &deckDoc{ID: DocID{docType: "chicken", id: []byte("abcd")}},
			err:  "incorrect doc type",
		},
		{
			name: "invalid doctype",
			v:    &deckDoc{ID: DocID{docType: "theme", id: []byte("abcd")}},
			err:  "incorrect doc type",
		},
		{
			name: "no created time",
			v:    &deckDoc{ID: DocID{docType: "deck", id: []byte("abcd")}},
			err:  "created time required",
		},
		{
			name: "no modified time",
			v:    &deckDoc{ID: DocID{docType: "deck", id: []byte("abcd")}, Created: now()},
			err:  "modified time required",
		},
		{
			name: "invalid card",
			v:    &deckDoc{ID: DocID{docType: "deck", id: []byte("abcd")}, Created: now(), Modified: now(), Cards: &CardCollection{col: map[string]struct{}{"foo": {}}}},
			err:  "'foo': invalid ID type",
		},
		{
			name: "valid",
			v:    &deckDoc{ID: DocID{docType: "deck", id: []byte("abcd")}, Created: now(), Modified: now(), Cards: &CardCollection{col: map[string]struct{}{"card-abcd.abcd.0": {}}}},
		},
	}
	testValidation(t, tests)
}
