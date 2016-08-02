package test

import (
	"encoding/json"
	"testing"

	"github.com/flimzy/flashback-model"
	. "github.com/flimzy/flashback-model/test/util"
)

var frozenNote []byte = []byte(`
{
    "type": "note",
    "_id": "note-VGVzdCBOb3RlCg==",
    "model": "NVXGa7SD7zl4CpU_-R7o-qwAZs8=.1",
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

func TestCreateNote(t *testing.T) {
	th := &fb.Theme{}
	json.Unmarshal(frozenTheme, th)
	m := th.Models[1]
	n, err := fb.NewNote("VGVzdCBOb3RlCg==", m)
	if err != nil {
		t.Fatalf("Unable to create Note: %s\n", err)
	}
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
}
