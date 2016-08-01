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
    "_id": "theme-NVXGa7SD7zl4CpU_-R7o-qwAZs8=",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "name": "Test Theme",
    "description": "Theme for testing",
    "models": [
        {
            "id": 0,
            "name": "Model A",
            "files": [
                "m1.html"
            ]
        },
        {
            "id": 1,
            "name": "Model 2",
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
	th, err := fbmodel.NewTheme("NVXGa7SD7zl4CpU_-R7o-qwAZs8=")
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
	m1, _ := th.NewModel()
	m2, _ := th.NewModel()
	name1 := "Model A"
	name2 := "Model 2"
	m1.Name = &name1
	m2.Name = &name2
	m1.AddFile("m1.html", "text/html", []byte("<html></html>"))
	m2.AddFile("m1.txt", "text/plain", []byte("Test text file"))
	output := Marshal(t, "Create Theme", th)
	JSONDeepEqual(t, "Create Theme", output, frozenTheme)

	th2 := &fbmodel.Theme{}
	if err := json.Unmarshal(frozenTheme, th2); err != nil {
		t.Fatalf("Error thawing theme: %s\n", err)
	}
	JSONDeepEqual(t, "Thawed Theme", Marshal(t, "Thaw Theme", th2), frozenTheme)

	if !reflect.DeepEqual(th, th2) {
		PrintDiff(th2, th)
		t.Fatalf("Thawed and created Themes don't match")
	}
}
