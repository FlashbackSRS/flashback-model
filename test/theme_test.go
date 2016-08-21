package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/flimzy/flashback-model"
	. "github.com/flimzy/flashback-model/test/util"
)

var frozenTheme []byte = []byte(`
{
    "type": "theme",
    "_id": "theme-0NVXGa7SD7zl4CpU_-R7o-qwAZs8",
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
        "m1.html": {
            "content-type": "text/html",
            "data": "PGh0bWw+PC9odG1sPg=="
        },
        "m1.txt": {
            "content-type": "text/plain",
            "data": "VGVzdCB0ZXh0IGZpbGU="
        },
        "$main.css": {
            "content-type": "text/css",
            "data": "LyogYW4gZW1wdHkgQ1NTIGZpbGUgKi8="
        }
    },
    "files": [
        "$main.css"
    ],
    "modelSequence": 2
}
`)

func TestCreateTheme(t *testing.T) {
	th, err := fb.NewTheme("0NVXGa7SD7zl4CpU_-R7o-qwAZs8")
	if err != nil {
		t.Fatalf("Error creating theme: %s\n", err)
	}
	name := "Test Theme"
	th.Name = &name
	descr := "Theme for testing"
	th.Description = &descr
	th.Created = &now
	th.Modified = &now
	th.SetFile("$main.css", "text/css", []byte("/* an empty CSS file */"))
	m1, _ := th.NewModel(fb.ModelType(0))
	m2, _ := th.NewModel(fb.ModelType(1))
	m1.AddField(fb.TextField, "Word")
	m1.AddField(fb.TextField, "Definition")
	m2.AddField(fb.TextField, "Word")
	m2.AddField(fb.AudioField, "Audio")
	name1 := "Model A"
	name2 := "Model 2"
	m1.Name = &name1
	m2.Name = &name2
	m1.AddFile("m1.html", "text/html", []byte("<html></html>"))
	m2.AddFile("m1.txt", "text/plain", []byte("Test text file"))
	JSONDeepEqual(t, "Create Theme", Marshal(t, "Create Theme", th), frozenTheme)

	th2 := &fb.Theme{}
	if err := json.Unmarshal(frozenTheme, th2); err != nil {
		t.Fatalf("Error thawing theme: %s\n", err)
	}
	JSONDeepEqual(t, "Thawed Theme", Marshal(t, "Thaw Theme", th2), frozenTheme)

	if !reflect.DeepEqual(th, th2) {
		PrintDiff(th2, th)
		t.Fatal("Thawed and created Themes don't match")
	}
}
