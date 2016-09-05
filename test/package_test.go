package test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/FlashbackSRS/flashback-model"
	. "github.com/FlashbackSRS/flashback-model/test/util"
)

var frozenPackage = []byte(`
{
    "bundle": {
        "type": "bundle",
        "_id": "bundle-krsxg5baij2w4zdmmu",
        "created": "2016-07-31T15:08:24.730156517Z",
        "modified": "2016-07-31T15:08:24.730156517Z",
        "owner": "tui5ajfbabaeljnxt4om7fwmt4"
    },
    "cards": [
        {
            "type": "card",
            "_id": "card-krsxg5baij2w4zdmmu.VGVzdCBOb3Rl.0",
            "created": "2016-07-31T15:08:24.730156517Z",
            "modified": "2016-07-31T15:08:24.730156517Z"
        },
        {
            "type": "card",
            "_id": "card-krsxg5baij2w4zdmmu.VGVzdCBOb3Rl.1",
            "created": "2016-07-31T15:08:24.730156517Z",
            "modified": "2016-07-31T15:08:24.730156517Z"
        },
        {
            "type": "card",
            "_id": "card-krsxg5baij2w4zdmmu.VGVzdCBOb3Rl.2",
            "created": "2016-07-31T15:08:24.730156517Z",
            "modified": "2016-07-31T15:08:24.730156517Z"
        }
    ],
    "notes": [
        {
            "type": "note",
            "_id": "note-VGVzdCBOb3Rl",
            "created": "2016-07-31T15:08:24.730156517Z",
            "modified": "2016-07-31T15:08:24.730156517Z",
            "imported": "2016-08-02T15:08:24.730156517Z",
            "theme": "VGVzdCBUaGVtZQ",
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
    ],
    "decks": [
        {
            "type": "deck",
            "_id": "deck-VGVzdCBEZWNr",
            "created": "2016-07-31T15:08:24.730156517Z",
            "modified": "2016-07-31T15:08:24.730156517Z",
            "imported": "2016-08-02T15:08:24.730156517Z",
            "name": "Test Deck",
            "description": "Deck for testing",
            "cards": []
        }
    ],
    "themes": [
        {
            "type": "theme",
            "_id": "theme-VGVzdCBUaGVtZQ",
            "created": "2016-07-31T15:08:24.730156517Z",
            "modified": "2016-07-31T15:08:24.730156517Z",
            "imported": "2016-08-02T15:08:24.730156517Z",
            "name": "Test Theme",
            "description": "Theme for testing",
            "models": [
                {
                    "id": 0,
                    "modelType": 0,
                    "name": "Model A",
                    "templates": [],
                    "fields": [
                        {
                            "fieldType": 0,
                            "name": "Word"
                        },
                        {
                            "fieldType": 0,
                            "name": "Definition"
                        }
                    ],
                    "files": [
                        "m1.html"
                    ]
                },
                {
                    "id": 1,
                    "modelType": 1,
                    "name": "Model 2",
                    "templates": [],
                    "fields": [
                        {
                            "fieldType": 0,
                            "name": "Word"
                        },
                        {
                            "fieldType": 2,
                            "name": "Audio"
                        }
                    ],
                    "files": [
                        "m1.txt"
                    ]
                }
            ],
            "_attachments": {
                "$main.css": {
                    "content_type": "text/css",
                    "data": "LyogYW4gZW1wdHkgQ1NTIGZpbGUgKi8="
                },
                "m1.html": {
                    "content_type": "text/html",
                    "data": "PGh0bWw+PC9odG1sPg=="
                },
                "m1.txt": {
                    "content_type": "text/plain",
                    "data": "VGVzdCB0ZXh0IGZpbGU="
                }
            },
            "files": [
                "$main.css"
            ],
            "modelSequence": 2
        }
    ],
    "reviews": [
        {
            "cardID": "VGVzdCBOb3Rl.0",
            "timestamp": null
        }
    ]
}
`)

func TestPackage(t *testing.T) {
	u, _ := testUser()
	b, err := fb.NewBundle([]byte("Test Bundle"), u)
	if err != nil {
		t.Fatalf("Error creating bundle: %s", err)
	}
	b.Created = now
	b.Modified = now
	th := &fb.Theme{}
	json.Unmarshal(frozenTheme, th)
	d := &fb.Deck{}
	json.Unmarshal(frozenDeck, d)
	n := &fb.Note{}
	json.Unmarshal(frozenNote, n)
	r := &fb.Review{}
	json.Unmarshal(frozenReview, r)
	p := &fb.Package{
		Bundle:  b,
		Cards:   make([]*fb.Card, 0),
		Themes:  []*fb.Theme{th},
		Decks:   []*fb.Deck{d},
		Notes:   []*fb.Note{n},
		Reviews: []*fb.Review{r},
	}

	for i := 0; i < 3; i++ {
		c, err := fb.NewCard(fmt.Sprintf("%s.%s.%d", b.ID.Identity(), n.ID.Identity(), i))
		if err != nil {
			t.Fatalf("Error creating new card: %s", err)
		}
		c.Created = now
		c.Modified = now
		p.Cards = append(p.Cards, c)
	}
	JSONDeepEqual(t, "Create Package", Marshal(t, "Create Package", p), frozenPackage)

	p2 := &fb.Package{}
	if err := json.Unmarshal(frozenPackage, p2); err != nil {
		t.Fatalf("Error thawing package: %s", err)
	}
	JSONDeepEqual(t, "Thawed Package", Marshal(t, "Thaw Package", p2), frozenPackage)

	// We have to set the username explicitly for the next test to pass, as a simple unmarshaling
	// of a bundle doesn't know user details (nor should it)
	p2.Bundle.Owner.Username = "mrsmith"
	if !reflect.DeepEqual(p, p2) {
		PrintDiff(p2, p)
		t.Fatal("Thawed and created Packages don't match")
	}
}
