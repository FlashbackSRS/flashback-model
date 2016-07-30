package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/flimzy/flashback-model"
)

var frozenUser string = `{"_id":"user-9d11d024-a100-4045-a5b7-9f1ccf96cc9f","type":"user","password":"","salt":"","userType":"","name":"Bob Smith"}`

func TestUser(t *testing.T) {
	o := testUser()
	fmt.Printf("user = %v\n", o)
	t.Fatalf("x")
// 	output, err := json.Marshal(o)
// 	if err != nil {
// 		t.Errorf("Error marshaling bundle: %s", err)
// 	}
// 	fmt.Printf("Output = %s\n", string(output))
}

func testUser() *fbmodel.User {
	u := &fbmodel.User{}
	json.Unmarshal([]byte(frozenUser), u)
	return u
}
