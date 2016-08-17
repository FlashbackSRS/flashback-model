package test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/flimzy/flashback-model"
	. "github.com/flimzy/flashback-model/test/util"
)

var frozenPackage []byte = []byte(`
{
    "bundle": {
        "type": "bundle",
        "_id": "bundle-VjMOV9J35iuH1lXdM_lgQPOYx9I=",
        "owner": "nRHQJKEAQEWlt58cz5bMnw=="
    },
    "cards": [
        {
            "type": "card",
            "_id": "card-mViuXQThMLoh1G1Nlc4d_E8kR8o=.0"
        },
        {
            "type": "card",
            "_id": "card-mViuXQThMLoh1G1Nlc4d_E8kR8o=.1"
        },
        {
            "type": "card",
            "_id": "card-mViuXQThMLoh1G1Nlc4d_E8kR8o=.2"
        }
    ],
    "notes": [
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
    ],
    "decks": [
        {
            "type": "deck",
            "_id": "deck-AO1yee9FPLVtU3h0M5pcYy3AOTQ=",
            "created": "2016-07-31T15:08:24.730156517Z",
            "modified": "2016-07-31T15:08:24.730156517Z",
            "name": "Test Deck",
            "description": "Deck for testing",
            "cards": []
        }
    ],
    "themes": [
        {
            "type": "theme",
            "_id": "theme-NVXGa7SD7zl4CpU_-R7o-qwAZs8=",
            "created": "2016-07-31T15:08:24.730156517Z",
            "modified": "2016-07-31T15:08:24.730156517Z",
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
                    "content-type": "text/css",
                    "data": "LyogYW4gZW1wdHkgQ1NTIGZpbGUgKi8="
                },
                "m1.html": {
                    "content-type": "text/html",
                    "data": "PGh0bWw+PC9odG1sPg=="
                },
                "m1.txt": {
                    "content-type": "text/plain",
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
            "cardID": "mViuXQThMLoh1G1Nlc4d_E8kR8o=.0",
            "timestamp": null,
            "ease": 0,
            "interval": null,
            "previousInterval": null,
            "srsFactor": 0,
            "reviewTime": null,
            "reviewType": 0
        }
    ]
}
`)

func TestPackage(t *testing.T) {
	u, _ := testUser()
	b, _ := fb.NewBundle("VjMOV9J35iuH1lXdM_lgQPOYx9I=", u)
	th := &fb.Theme{}
	json.Unmarshal(frozenTheme, th)
	d := &fb.Deck{}
	json.Unmarshal(frozenDeck, d)
	n := &fb.Note{}
	fmt.Printf("model = %v\n", th.Models[1])
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
		c, _ := fb.NewCard("mViuXQThMLoh1G1Nlc4d_E8kR8o=", i)
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
