package test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/flimzy/testify/require"

	"github.com/FlashbackSRS/flashback-model"
)

var frozenPackage = []byte(`
{
    "version": 0,
    "bundle": {
        "type": "bundle",
        "_id": "bundle-krsxg5baij2w4zdmmu",
        "created": "2016-07-31T15:08:24.730156517Z",
        "modified": "2016-07-31T15:08:24.730156517Z",
        "owner": "user-tui5ajfbabaeljnxt4om7fwmt4"
    },
    "cards": [
        {
            "type": "card",
            "_id": "card-krsxg5baij2w4zdmmu.VGVzdCBOb3Rl.0",
            "created": "2016-07-31T15:08:24.730156517Z",
            "modified": "2016-07-31T15:08:24.730156517Z",
            "model": "theme-VGVzdCBUaGVtZQ/0"
        },
        {
            "type": "card",
            "_id": "card-krsxg5baij2w4zdmmu.VGVzdCBOb3Rl.1",
            "created": "2016-07-31T15:08:24.730156517Z",
            "modified": "2016-07-31T15:08:24.730156517Z",
            "model": "theme-VGVzdCBUaGVtZQ/0"
        },
        {
            "type": "card",
            "_id": "card-krsxg5baij2w4zdmmu.VGVzdCBOb3Rl.2",
            "created": "2016-07-31T15:08:24.730156517Z",
            "modified": "2016-07-31T15:08:24.730156517Z",
            "model": "theme-VGVzdCBUaGVtZQ/0"
        }
    ],
    "notes": [
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
                        "_weirdname.txt",
                        "foo.mp3",
                        "영상.jpg"
                    ]
                }
            ],
            "_attachments": {
                "^_weirdname.txt": {
                    "content_type": "audio/mpeg",
                    "data": "YSBmaWxlIHdpdGggYSBzdHJhbmdlIG5hbWU="
                },
                "영상.jpg": {
                    "content_type": "audio/mpeg",
                    "data": "YSBLb3JlYW4gZmlsZW5hbWU="
                },
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
                    "modelType": "anki-basic",
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
                    "modelType": "anki-cloze",
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
            "cardID": "card-krsxg5baij2w4zdmmu.mViuXQThMLoh1G1Nlc4d_E8kR8o.0",
            "timestamp": "2017-01-01T00:00:00Z"
        }
    ]
}
`)

func TestPackage(t *testing.T) {
	require := require.New(t)
	b, err := fb.NewBundle("bundle-krsxg5baij2w4zdmmu", "user-tui5ajfbabaeljnxt4om7fwmt4")
	require.Nil(err, "Error creating bundle: %s", err)

	b.Created = now
	b.Modified = now
	th := &fb.Theme{}
	json.Unmarshal(frozenTheme, th)
	d := &fb.Deck{}
	json.Unmarshal(frozenDeck, d)
	n := &fb.Note{}
	json.Unmarshal(frozenNote, n)
	n.Model = th.Models[1]
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
		c, e := fb.NewCard("theme-VGVzdCBUaGVtZQ", 0, fmt.Sprintf("card-%s.%s.%d", strings.TrimPrefix(b.ID, "bundle-"), strings.TrimPrefix(n.ID, "note-"), i))
		require.Nil(e, "Error creating new card: %s", err)
		c.Created = now
		c.Modified = now
		p.Cards = append(p.Cards, c)
	}
	require.MarshalsToJSON(frozenPackage, p, "Create Package")

	p2 := &fb.Package{}
	err = json.Unmarshal(frozenPackage, p2)
	require.Nil(err, "Error thawing package: %s", err)
	require.MarshalsToJSON(frozenPackage, p2, "Thawed Package")

	// We have to set the username explicitly for the next test to pass, as a simple unmarshaling
	// of a bundle doesn't know user details (nor should it)
	require.DeepEqual(p2, p, "Thawed vs Created Packages")
}
