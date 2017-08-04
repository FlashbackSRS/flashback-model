package fb

import (
	"encoding/json"
	"testing"
	"time"

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

func TestMarshalJSON(t *testing.T) {
	card := &Card{
		bundleID:   "foo",
		themeID:    "bar",
		noteID:     "baz",
		templateID: 1,
		modelID:    2,
		Created:    parseTime("2017-01-01T01:01:01Z"),
		Modified:   parseTime("2017-01-01T01:01:01Z"),
		Suspended:  true,
	}
	expected := []byte(`{
        "_id": "card-foo.baz.1",
        "created": "2017-01-01T01:01:01Z",
        "model": "bar/2",
        "modified": "2017-01-01T01:01:01Z",
        "suspended": true,
        "type": "card"
    }`)
	result, err := json.Marshal(card)
	checkErr(t, nil, err)
	if d := diff.JSON(expected, result); d != "" {
		t.Error(d)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type ujTest struct {
		name     string
		input    string
		expected *Card
		err      string
	}
	tests := []ujTest{
		{
			name: "no input",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "wrong type",
			input: `{"type":"chicken"}`,
			err:   "invalid document type for card: chicken",
		},
		{
			name:  "invalid id",
			input: `{"type":"card","_id":"oink"}`,
			err:   "invalid ID type",
		},
		{
			name:  "invalid model id",
			input: `{"type":"card", "_id":"card-krsxg5baij2w4zdmmu.mViuXQThMLoh1G1Nlc4d_E8kR8o.1", "model": "foo/chicken"}`,
			err:   `invalid model ID: strconv.Atoi: parsing "chicken": invalid syntax`,
		},
		{
			name:  "valid",
			input: `{"type":"card", "_id":"card-krsxg5baij2w4zdmmu.mViuXQThMLoh1G1Nlc4d_E8kR8o.1", "model": "foo/2", "suspended":true}`,
			expected: &Card{
				bundleID:   "krsxg5baij2w4zdmmu",
				noteID:     "mViuXQThMLoh1G1Nlc4d_E8kR8o",
				templateID: 1,
				themeID:    "foo",
				modelID:    2,
				Suspended:  true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &Card{}
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

func TestIdentity(t *testing.T) {
	card := &Card{bundleID: "bundle", noteID: "note", templateID: 2}
	expected := "bundle.note.2"
	result := card.Identity()
	if result != expected {
		t.Errorf("Unexpected result: %s", result)
	}
}

func TestSetRev(t *testing.T) {
	card := &Card{}
	rev := "1-xxx"
	card.SetRev(rev)
	if *card.Rev != rev {
		t.Errorf("Unexpected rev: %s", *card.Rev)
	}
}

func TestDocID(t *testing.T) {
	card := &Card{bundleID: "bundle", noteID: "note", templateID: 2}
	expected := "card-bundle.note.2"
	result := card.DocID()
	if result != expected {
		t.Errorf("Unexpected result: %s", result)
	}
}

func TestImportedTime(t *testing.T) {
	ts := time.Now()
	card := &Card{Imported: &ts}
	if it := card.ImportedTime(); !it.Equal(ts) {
		t.Errorf("Unexpected result: %v", *it)
	}
}

func TestModifiedTime(t *testing.T) {
	ts := time.Now()
	card := &Card{Modified: ts}
	if mt := card.ModifiedTime(); !mt.Equal(ts) {
		t.Errorf("Unexpected result: %v", *mt)
	}
}

func TestMergeImport(t *testing.T) {
	type miTest struct {
		name         string
		card         *Card
		i            interface{}
		expected     bool
		expectedCard *Card
		err          string
	}
	tests := []miTest{
		{
			name: "no input",
			card: &Card{},
			i:    nil,
			err:  "i is <nil>, not *fb.Card",
		},
		{
			name: "mismatched identities",
			card: &Card{bundleID: "foo", noteID: "bar", templateID: 1},
			i:    &Card{bundleID: "foo", noteID: "bar", templateID: 2},
			err:  "IDs don't match",
		},
		{
			name: "different timestamps",
			card: &Card{bundleID: "foo", noteID: "bar", Created: parseTime("2017-01-01T01:01:01Z")},
			i:    &Card{bundleID: "foo", noteID: "bar", Created: parseTime("2017-02-01T01:01:01Z")},
			err:  "Created timestamps don't match",
		},
		{
			name:         "existing is newer",
			card:         &Card{bundleID: "foo", noteID: "bar", Created: parseTime("2017-01-01T01:01:01Z"), Modified: parseTime("2017-01-01T01:01:01Z")},
			i:            &Card{bundleID: "foo", noteID: "bar", Created: parseTime("2017-01-01T01:01:01Z"), Modified: parseTime("2017-01-02T01:01:01Z")},
			expectedCard: &Card{bundleID: "foo", noteID: "bar", Created: parseTime("2017-01-01T01:01:01Z"), Modified: parseTime("2017-01-02T01:01:01Z")},
		},
		{
			name:         "new is newer",
			card:         &Card{bundleID: "foo", noteID: "bar", Created: parseTime("2017-01-01T01:01:01Z"), Modified: parseTime("2017-01-02T01:01:01Z")},
			i:            &Card{bundleID: "foo", noteID: "bar", Created: parseTime("2017-01-01T01:01:01Z"), Modified: parseTime("2017-01-01T01:01:01Z")},
			expected:     true,
			expectedCard: &Card{bundleID: "foo", noteID: "bar", Created: parseTime("2017-01-01T01:01:01Z"), Modified: parseTime("2017-01-02T01:01:01Z")},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.card.MergeImport(test.i)
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if result != test.expected {
				t.Errorf("Unexpected result: %t", result)
			}
			if d := diff.Interface(test.expectedCard, test.card); d != "" {
				t.Error(d)
			}
		})
	}
}

func TestBundleID(t *testing.T) {
	card := &Card{bundleID: "foo"}
	expected := "bundle-foo"
	if id := card.BundleID(); id != expected {
		t.Errorf("Unexpected result: %s", id)
	}
}

func TestTemplateID(t *testing.T) {
	expected := uint32(3)
	card := &Card{templateID: expected}
	if id := card.TemplateID(); id != expected {
		t.Errorf("Unexpected result: %d", id)
	}
}

func TestModelID(t *testing.T) {
	expected := 4
	card := &Card{modelID: uint32(expected)}
	if id := card.ModelID(); id != expected {
		t.Errorf("Unexpected result: %d", id)
	}
}

func TestNoteID(t *testing.T) {
	card := &Card{noteID: "bar"}
	expected := "note-bar"
	if id := card.NoteID(); id != expected {
		t.Errorf("Unexpected result: %s", id)
	}
}

func TestCardValidate(t *testing.T) {
	type cvTest struct {
		name string
		card *cardDoc
		err  string
	}
	tests := []cvTest{
		{
			name: "empty card",
			card: &cardDoc{},
			err:  "id required",
		},
		{
			name: "invalid id",
			card: &cardDoc{ID: "chicken"},
			err:  "invalid ID type",
		},
		{
			name: "zero created time",
			card: &cardDoc{ID: "card-foo.bar.0"},
			err:  "created time required",
		},
		{
			name: "zero modified time",
			card: &cardDoc{ID: "card-foo.bar.0", Created: parseTime("2017-01-01T01:01:01Z")},
			err:  "modified time required",
		},
		{
			name: "missing model id",
			card: &cardDoc{ID: "card-foo.bar.0", Created: parseTime("2017-01-01T01:01:01Z"), Modified: parseTime("2017-01-01T01:01:01Z")},
			err:  "model id required",
		},
		{
			name: "invalid model id",
			card: &cardDoc{ID: "card-foo.bar.0", Created: parseTime("2017-01-01T01:01:01Z"), Modified: parseTime("2017-01-01T01:01:01Z"),
				ModelID: "chicken"},
			err: "invalid type in model ID",
		},
		{
			name: "valid",
			card: &cardDoc{ID: "card-foo.bar.0", Created: parseTime("2017-01-01T01:01:01Z"), Modified: parseTime("2017-01-01T01:01:01Z"),
				ModelID: "model-foo"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.card.Validate()
			checkErr(t, test.err, err)
		})
	}
}
