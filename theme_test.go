package fb

import (
	"testing"

	"github.com/flimzy/diff"
)

func TestThemeUnmarshalJSON(t *testing.T) {
	type Test struct {
		name     string
		input    string
		expected interface{}
		err      string
	}
	tests := []Test{
		{
			name:  "invalid json",
			input: "xxx",
			err:   "failed to unmarshal Theme: invalid character 'x' looking for beginning of value",
		},
		{
			name:  "wrong type",
			input: `{"type":"chicken"}`,
			err:   "Invalid document type for theme: chicken",
		},
		{
			name:  "no attachments",
			input: `{"type":"theme", "_id":"theme-120","created":"2017-01-01T01:01:01Z","modified":"2017-01-01T01:01:01Z"}`,
			err:   "invalid theme: no attachments",
		},
		{
			name: "with attachments",
			input: `{"type":"theme", "_id":"theme-120", "created":"2017-01-01T01:01:01Z", "modified":"2017-01-01T01:01:01Z", "_attachments":{
            "foo.txt": {"content_type":"text/plain", "content": "text"}
            }}`,
			err: "invalid theme: no file list",
		},
		{
			name:  "mismatched file list",
			input: `{"type":"theme", "_id":"theme-120", "created":"2017-01-01T01:01:01Z", "modified":"2017-01-01T01:01:01Z", "_attachments": {"foo.txt": {"content_type":"text/plain", "content": "text"}}, "files": ["foo.html"] }`,
			err:   "foo.html not found in collection",
		},
		{
			name:  "mismatched model file list",
			input: `{"type":"theme", "_id":"theme-120", "created":"2017-01-01T01:01:01Z", "modified":"2017-01-01T01:01:01Z", "_attachments": {"foo.txt": {"content_type":"text/plain", "content": "text"}}, "files":[], "models": [{"id":0, "files": ["foo.mp3"]}] }`,
			err:   "foo.mp3 not found in collection",
		},
		{
			name:  "valid",
			input: `{"type":"theme", "_id":"theme-120", "created":"2017-01-01T01:01:01Z", "modified":"2017-01-01T01:01:01Z", "_attachments": {"foo.txt": {"content_type":"text/plain", "data": "text"}}, "files":[], "models": [{"id":0, "files": ["foo.txt"]}] }`,
			expected: map[string]interface{}{
				"type":     "theme",
				"_id":      "theme-120",
				"created":  "2017-01-01T01:01:01Z",
				"modified": "2017-01-01T01:01:01Z",
				"models": []map[string]interface{}{
					{
						"id":        0,
						"modelType": "",
						"templates": nil,
						"fields":    nil,
						"files":     []string{"foo.txt"},
					},
				},
				"_attachments": map[string]interface{}{
					"foo.txt": map[string]string{
						"content_type": "text/plain",
						"data":         "text",
					},
				},
				"files":         []string{},
				"modelSequence": 0,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &Theme{}
			err := result.UnmarshalJSON([]byte(test.input))
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if d := diff.AsJSON(test.expected, result); d != "" {
				t.Error(d)
			}
		})
	}
}
