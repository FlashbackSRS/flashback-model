package fb

import (
	"testing"
	"time"

	"github.com/flimzy/diff"
)

func TestNewBundle(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		owner    *User
		expected *Bundle
		err      string
	}{
		{
			name: "no id",
			err:  "id required",
		},
		{
			name: "no owner",
			id:   "foo",
			err:  "owner required",
		},
		{
			name:  "valid",
			id:    "foo",
			owner: &User{},
			expected: &Bundle{
				ID:    DbID{docType: "bundle", id: []byte("foo")},
				Owner: &User{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewBundle([]byte(test.id), test.owner)
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

func TestBundleMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		bundle   *Bundle
		expected string
		err      string
	}{
		{
			name: "null fields",
			bundle: &Bundle{
				ID: DbID{docType: "bundle", id: []byte("foo")},
				Owner: &User{
					ID: DbID{docType: "user", id: []byte("bob")},
				},
				Created:  now(),
				Modified: now(),
			},
			expected: `{
                "_id":      "bundle-mzxw6",
                "type":     "bundle",
                "owner":    "mjxwe",
                "created":  "2017-01-01T00:00:00Z",
                "modified": "2017-01-01T00:00:00Z"
            }`,
		},
		{
			name: "all fields",
			bundle: &Bundle{
				ID: DbID{docType: "bundle", id: []byte("foo")},
				Owner: &User{
					ID: DbID{docType: "user", id: []byte("bob")},
				},
				Created:     now(),
				Modified:    now(),
				Imported:    func() *time.Time { x := now(); return &x }(),
				Name:        func() *string { x := "foo name"; return &x }(),
				Description: func() *string { x := "foo description"; return &x }(),
			},
			expected: `{
                "_id":         "bundle-mzxw6",
                "type":        "bundle",
                "owner":       "mjxwe",
                "name":        "foo name",
                "description": "foo description",
                "created":     "2017-01-01T00:00:00Z",
                "modified":    "2017-01-01T00:00:00Z",
                "imported":    "2017-01-01T00:00:00Z"
            }`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.bundle.MarshalJSON()
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

func TestBundleUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Bundle
		err      string
	}{
		{
			name:  "invalid json",
			input: "invalid json",
			err:   "failed to unmarshal Bundle: invalid character 'i' looking for beginning of value",
		},
		{
			name:  "wrong doc type",
			input: `{"type":"chicken"}`,
			err:   "Invalid document type for bundle: chicken",
		},
		{
			name:  "invalid user",
			input: `{"type":"bundle","owner":"unf"}`,
			err:   "invalid user for bundle: invalid DbID: illegal base32 data at input byte 3",
		},
		{
			name: "null fiels",
			input: `{
                "_id":      "bundle-mzxw6",
                "type":     "bundle",
                "owner":    "mjxwe",
                "created":  "2017-01-01T00:00:00Z",
                "modified": "2017-01-01T00:00:00Z"
            }`,
			expected: &Bundle{
				ID: DbID{docType: "bundle", id: []byte("foo")},
				Owner: &User{
					ID:   DbID{docType: "user", id: []byte("bob")},
					uuid: []byte("bob"),
				},
				Created:  now(),
				Modified: now(),
			},
		},
		{
			name: "all fields",
			input: `{
                "_id":         "bundle-mzxw6",
                "type":        "bundle",
                "owner":       "mjxwe",
                "name":        "foo name",
                "description": "foo description",
                "created":     "2017-01-01T00:00:00Z",
                "modified":    "2017-01-01T00:00:00Z",
                "imported":    "2017-01-01T00:00:00Z"
            }`,
			expected: &Bundle{
				ID: DbID{docType: "bundle", id: []byte("foo")},
				Owner: &User{
					ID:   DbID{docType: "user", id: []byte("bob")},
					uuid: []byte("bob"),
				},
				Created:     now(),
				Modified:    now(),
				Imported:    func() *time.Time { x := now(); return &x }(),
				Name:        func() *string { x := "foo name"; return &x }(),
				Description: func() *string { x := "foo description"; return &x }(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &Bundle{}
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

func TestBundleSetRev(t *testing.T) {
	bundle := &Bundle{}
	rev := "1-xxx"
	bundle.SetRev(rev)
	if *bundle.Rev != rev {
		t.Errorf("failed to set rev")
	}
}

func TestBundleID(t *testing.T) {
	expected := "bundle-mzxw6"
	bundle := &Bundle{ID: DbID{docType: "bundle", id: []byte("foo")}}
	if id := bundle.DocID(); id != expected {
		t.Errorf("unexpected id: %s", id)
	}
}

func TestBundleImportedTime(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		bundle := &Bundle{}
		ts := now()
		bundle.Imported = &ts
		if it := bundle.ImportedTime(); it != ts {
			t.Errorf("Unexpected result: %s", it)
		}
	})
	t.Run("Unset", func(t *testing.T) {
		bundle := &Bundle{}
		if it := bundle.ImportedTime(); !it.IsZero() {
			t.Errorf("unexpected result: %v", it)
		}
	})
}

func TestBundleModifiedTime(t *testing.T) {
	bundle := &Bundle{}
	ts := now()
	bundle.Modified = ts
	if mt := bundle.ModifiedTime(); mt != ts {
		t.Errorf("Unexpected result")
	}
}

func TestBundleMergeImport(t *testing.T) {
	type Test struct {
		name           string
		new            *Bundle
		existing       *Bundle
		expected       bool
		expectedBundle *Bundle
		err            string
	}
	tests := []Test{
		{
			name:     "different ids",
			new:      &Bundle{ID: DbID{docType: "bundle", id: []byte("foo")}},
			existing: &Bundle{ID: DbID{docType: "bundle", id: []byte("bar")}},
			err:      "IDs don't match",
		},
		{
			name:     "created timestamps don't match",
			new:      &Bundle{ID: DbID{docType: "bundle", id: []byte("foo")}, Created: parseTime("2017-01-01T01:01:01Z"), Imported: parseTimePtr("2017-01-15T00:00:00Z")},
			existing: &Bundle{ID: DbID{docType: "bundle", id: []byte("foo")}, Created: parseTime("2017-02-01T01:01:01Z"), Imported: parseTimePtr("2017-01-20T00:00:00Z")},
			err:      "Created timestamps don't match",
		},
		{
			name:     "owners don't match",
			new:      &Bundle{ID: DbID{docType: "bundle", id: []byte("foo")}, Owner: &User{uuid: []byte("bob")}, Created: parseTime("2017-01-01T01:01:01Z"), Imported: parseTimePtr("2017-01-15T00:00:00Z")},
			existing: &Bundle{ID: DbID{docType: "bundle", id: []byte("foo")}, Owner: &User{uuid: []byte("alice")}, Created: parseTime("2017-01-01T01:01:01Z"), Imported: parseTimePtr("2017-01-20T00:00:00Z")},
			err:      "Cannot change bundle ownership",
		},
		{
			name:     "new not an import",
			new:      &Bundle{ID: DbID{docType: "bundle", id: []byte("foo")}, Owner: &User{uuid: []byte("bob")}, Created: parseTime("2017-01-01T01:01:01Z")},
			existing: &Bundle{ID: DbID{docType: "bundle", id: []byte("foo")}, Owner: &User{uuid: []byte("bob")}, Created: parseTime("2017-01-01T01:01:01Z"), Imported: parseTimePtr("2017-01-15T00:00:00Z")},
			err:      "not an import",
		},
		{
			name:     "existing not an import",
			new:      &Bundle{ID: DbID{docType: "bundle", id: []byte("foo")}, Owner: &User{uuid: []byte("bob")}, Created: parseTime("2017-01-01T01:01:01Z"), Imported: parseTimePtr("2017-01-15T00:00:00Z")},
			existing: &Bundle{ID: DbID{docType: "bundle", id: []byte("foo")}, Owner: &User{uuid: []byte("bob")}, Created: parseTime("2017-01-01T01:01:01Z")},
			err:      "not an import",
		},
		{
			name: "new is newer",
			new: &Bundle{
				ID:          DbID{docType: "bundle", id: []byte("foo")},
				Owner:       &User{uuid: []byte("bob")},
				Name:        func() *string { x := "foo"; return &x }(),
				Description: func() *string { x := "FOO"; return &x }(),
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-02-01T01:01:01Z"),
				Imported:    parseTimePtr("2017-01-15T00:00:00Z"),
			},
			existing: &Bundle{
				ID:          DbID{docType: "bundle", id: []byte("foo")},
				Owner:       &User{uuid: []byte("bob")},
				Name:        func() *string { x := "bar"; return &x }(),
				Description: func() *string { x := "BAR"; return &x }(),
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-01-01T01:01:01Z"),
				Imported:    parseTimePtr("2017-01-20T00:00:00Z"),
			},
			expected: true,
			expectedBundle: &Bundle{
				ID:          DbID{docType: "bundle", id: []byte("foo")},
				Owner:       &User{uuid: []byte("bob")},
				Name:        func() *string { x := "foo"; return &x }(),
				Description: func() *string { x := "FOO"; return &x }(),
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-02-01T01:01:01Z"),
				Imported:    parseTimePtr("2017-01-15T00:00:00Z"),
			},
		},
		{
			name: "existing is newer",
			new: &Bundle{
				ID:          DbID{docType: "bundle", id: []byte("foo")},
				Owner:       &User{uuid: []byte("bob")},
				Name:        func() *string { x := "foo"; return &x }(),
				Description: func() *string { x := "FOO"; return &x }(),
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-01-01T01:01:01Z"),
				Imported:    parseTimePtr("2017-01-15T00:00:00Z"),
			},
			existing: &Bundle{
				ID:          DbID{docType: "bundle", id: []byte("foo")},
				Owner:       &User{uuid: []byte("bob")},
				Name:        func() *string { x := "bar"; return &x }(),
				Description: func() *string { x := "BAR"; return &x }(),
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-02-01T01:01:01Z"),
				Imported:    parseTimePtr("2017-01-20T00:00:00Z"),
			},
			expected: false,
			expectedBundle: &Bundle{
				ID:          DbID{docType: "bundle", id: []byte("foo")},
				Owner:       &User{uuid: []byte("bob")},
				Name:        func() *string { x := "bar"; return &x }(),
				Description: func() *string { x := "BAR"; return &x }(),
				Created:     parseTime("2017-01-01T01:01:01Z"),
				Modified:    parseTime("2017-02-01T01:01:01Z"),
				Imported:    parseTimePtr("2017-01-20T00:00:00Z"),
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
			if d := diff.Interface(test.expectedBundle, test.new); d != "" {
				t.Error(d)
			}
		})
	}
}
