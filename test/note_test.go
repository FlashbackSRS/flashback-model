package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/flimzy/flashback-model"
	. "github.com/flimzy/flashback-model/test/util"
)

var frozenNote []byte = []byte(`
{
    "type": "note",
    "_id": "note-0VGVzdCBOb3RlCg",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "model": "0NVXGa7SD7zl4CpU_-R7o-qwAZs8.1",
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
            "content-type": "audio/mpeg",
            "data": "bm90IGEgcmVhbCBNUDM="
        }
    }
}
`)

func TestNote(t *testing.T) {
	th := &fb.Theme{}
	json.Unmarshal(frozenTheme, th)
	m := th.Models[1]
	n, err := fb.NewNote("0VGVzdCBOb3RlCg", m)
	if err != nil {
		t.Fatalf("Unable to create Note: %s\n", err)
	}
	n.Created = &now
	n.Modified = &now
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
	JSONDeepEqual(t, "Create Note", Marshal(t, "Create Theme", n), frozenNote)

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
