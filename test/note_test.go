package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/FlashbackSRS/flashback-model"
	. "github.com/FlashbackSRS/flashback-model/test/util"
)

var frozenNote = []byte(`
{
    "type": "note",
    "_id": "note-VGVzdCBOb3Rl",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "imported": "2016-08-02T15:08:24.730156517Z",
    "model": "VGVzdCBUaGVtZQ.1",
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
	th := &fb.Theme{}
	json.Unmarshal(frozenTheme, th)
	m := th.Models[1]
	n, err := fb.NewNote([]byte("Test Note"), m)
	if err != nil {
		t.Fatalf("Unable to create Note: %s\n", err)
	}
	n.Created = now
	n.Modified = now
	imp := now.AddDate(0, 0, 2)
	n.Imported = &imp
	fv1 := n.GetFieldValue(0)
	fv1.SetText("cat")
	fv2 := n.GetFieldValue(1)
	if err := fv2.SetText("A vicious beast"); err == nil {
		t.Fatal("Should not be permitted to add text to an audio field")
	}
	if err := fv1.AddFile("foo.jpg", "image/jpeg", []byte("not really an image")); err == nil {
		t.Fatal("Should not be permitted to add an attachment to a text field")
	}
	if err := fv2.AddFile("foo.mp3", "audio/mpeg", []byte("not a real MP3")); err != nil {
		t.Fatal("Error attaching audio file")
	}
	JSONDeepEqual(t, "Create Note", Marshal(t, "Create Note", n), frozenNote)

	n2 := &fb.Note{}
	if err := json.Unmarshal(frozenNote, n2); err != nil {
		t.Fatalf("Error thawing note: %s", err)
	}
	// We have to set the model explicitly for the next test to pass
	n2.SetModel(m)
	JSONDeepEqual(t, "Thawed Note", Marshal(t, "Thaw Note", n2), frozenNote)

	if !reflect.DeepEqual(n, n2) {
		PrintDiff(n2, n)
		t.Fatal("Thawed and created Notes don't match")
	}

}

var frozenExistingNote = []byte(`
{
    "type": "note",
    "_id": "note-VGVzdCBOb3Rl",
    "_rev": "1-6e1b6fb5352429cf3013eab5d692aac8",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-15T15:07:24.730156517Z",
    "imported": "2016-08-01T15:08:24.730156517Z",
    "model": "VGVzdCBUaGVtZQ.1",
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
    "model": "VGVzdCBUaGVtZQ.1",
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
	th := &fb.Theme{}
	json.Unmarshal(frozenTheme, th)
	m := th.Models[1]
	n := &fb.Note{}
	if err := json.Unmarshal(frozenNote, n); err != nil {
		t.Fatalf("Error thawing Note: %s", err)
	}
	n.SetModel(m)
	e := &fb.Note{}
	if err := json.Unmarshal(frozenExistingNote, e); err != nil {
		t.Fatalf("Error thawing ExistingNote: %s", err)
	}
	e.SetModel(m)
	changed, err := n.MergeImport(e)
	if err != nil {
		t.Fatalf("Error merging Note: %s\n", err)
	}
	if !changed {
		t.Fatalf("No change!")
	}
	JSONDeepEqual(t, "Merged Note", Marshal(t, "Merge Note", n), frozenMergedNote)
}
