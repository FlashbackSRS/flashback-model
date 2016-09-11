package test

import (
	"encoding/json"
	"testing"

	"github.com/flimzy/testify/require"

	"github.com/FlashbackSRS/flashback-model"
)

var frozenNote = []byte(`
{
    "type": "note",
    "_id": "note-VGVzdCBOb3Rl",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "imported": "2016-08-02T15:08:24.730156517Z",
    "theme": "theme-VGVzdCBUaGVtZQ",
    "model": 1,
    "fieldValues": [
        {
            "text": "cat"
        },
        {
            "files": [
                "foo.mp3"
            ]
        }
    ],
    "_attachments": {
        "foo.mp3": {
            "content_type": "audio/mpeg",
            "data": "bm90IGEgcmVhbCBNUDM="
        }
    }
}
`)

func TestNote(t *testing.T) {
	require := require.New(t)
	th := &fb.Theme{}
	json.Unmarshal(frozenTheme, th)
	m := th.Models[1]
	n, err := fb.NewNote([]byte("Test Note"), m)
	require.Nil(err, "Unable to create Note: %s", err)

	n.Created = now
	n.Modified = now
	imp := now.AddDate(0, 0, 2)
	n.Imported = &imp
	fv1 := n.GetFieldValue(0)
	fv1.SetText("cat")
	fv2 := n.GetFieldValue(1)

	err = fv2.SetText("A vicious beast")
	require.NotNil(err, "Should not be permitted to add text to an audio field")

	err = fv1.AddFile("foo.jpg", "image/jpeg", []byte("not really an image"))
	require.NotNil(err, "Should not be permitted to add an attachment to a text field")

	err = fv2.AddFile("foo.mp3", "audio/mpeg", []byte("not a real MP3"))
	require.Nil(err, "Error attaching audio file: %s", err)

	require.MarshalsToJSON(frozenNote, n, "Create Note")

	n2 := &fb.Note{}
	if err := json.Unmarshal(frozenNote, n2); err != nil {
		t.Fatalf("Error thawing note: %s", err)
	}
	// We have to set the model explicitly for the next test to pass
	n2.SetModel(m)
	require.MarshalsToJSON(frozenNote, n2, "Thaed Note")

	require.DeepEqual(n, n2, "Thawed vs. Created Notes")
}

var frozenExistingNote = []byte(`
{
    "type": "note",
    "_id": "note-VGVzdCBOb3Rl",
    "_rev": "1-6e1b6fb5352429cf3013eab5d692aac8",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-15T15:07:24.730156517Z",
    "imported": "2016-08-01T15:08:24.730156517Z",
    "theme": "theme-VGVzdCBUaGVtZQ",
    "model": 1,
    "fieldValues": [
        {
            "text": "Cat"
        },
        {
            "files": [
                "foo.mp3"
            ]
        }
    ],
    "_attachments": {
        "foo.mp3": {
            "content_type": "audio/mpeg",
            "data": "bm90IGEgcmVhbCBNUDM="
        }
    }
}
`)

var frozenMergedNote = []byte(`
{
    "type": "note",
    "_id": "note-VGVzdCBOb3Rl",
    "_rev": "1-6e1b6fb5352429cf3013eab5d692aac8",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "imported": "2016-08-02T15:08:24.730156517Z",
    "theme": "theme-VGVzdCBUaGVtZQ",
    "model": 1,
    "fieldValues": [
        {
            "text": "cat"
        },
        {
            "files": [
                "foo.mp3"
            ]
        }
    ],
    "_attachments": {
        "foo.mp3": {
            "content_type": "audio/mpeg",
            "data": "bm90IGEgcmVhbCBNUDM="
        }
    }
}
`)

func TestNoteMergeImport(t *testing.T) {
	require := require.New(t)
	th := &fb.Theme{}
	json.Unmarshal(frozenTheme, th)
	m := th.Models[1]
	n := &fb.Note{}
	err := json.Unmarshal(frozenNote, n)
	require.Nil(err, "Error thawing Note: %s", err)

	n.SetModel(m)
	e := &fb.Note{}
	err = json.Unmarshal(frozenExistingNote, e)
	require.Nil(err, "Error thawing ExistingNote: %s", err)

	e.SetModel(m)
	changed, err := n.MergeImport(e)
	require.Nil(err, "Error merging Note: %s", err)
	require.True(changed, "No change!")
	require.MarshalsToJSON(frozenMergedNote, n, "Merged Note")
}
