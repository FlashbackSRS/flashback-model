package test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pborman/uuid"

	"github.com/flimzy/anki"
	"github.com/flimzy/flashback-model"
	"github.com/flimzy/flashback-model/ankiconv"
	. "github.com/flimzy/flashback-model/test/util"
)

const ApkgFile = "Art.apkg"

var UUID []byte = []byte{0xD1, 0xC9, 0x58, 0x7D, 0x88, 0xDF, 0x4A, 0x65, 0x89, 0x23, 0xF7, 0x3C, 0xDF, 0x6D, 0x1D, 0x70}
var now time.Time = time.Unix(1469977704, 730156517).UTC()

func TestConvert(t *testing.T) {
	apkg, err := anki.ReadFile(ApkgFile)
	if err != nil {
		t.Fatalf("Error opening test file: %s", err)
	}
	u, err := fbmodel.NewUser(uuid.UUID(UUID), "testuser")
	if err != nil {
		t.Fatalf("Error creating test user: %s\n", err)
	}
	b := ankiconv.NewBundle()
	b.SetNow(now)
	if err := b.Convert("Art", u, apkg); err != nil {
		t.Fatalf("Error converting APKG: %s\n", err)
	}
	output, err := json.MarshalIndent(b, "", "    ")
	if err != nil {
		t.Fatalf("Errorf marshaling bundle: %s\n", err)
	}
	JSONDeepEqual(t, "Converted Bundle", output, frozenBundle)
}
