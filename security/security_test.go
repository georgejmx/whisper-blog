package security

import (
	"testing"
	tp "whisper-blog/types"
)

const testHash = "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"

type MockController struct{}

/*** Defining the mock database controller methods used for all tests ***/
func (mc *MockController) Init() error { return nil }
func (mc *MockController) GrabPosts() ([]tp.Post, error) {
	testPost := tp.Post{Id: 1, Title: "test", Author: "tester",
		Contents: "test contents", Tag: 0, Descriptors: "t;t;t", Time: "4:20"}
	return []tp.Post{testPost}, nil
}
func (mc *MockController) AddPost(post tp.Post) error { return nil }
func (mc *MockController) SelectHash() (string, error) {
	return testHash, nil
}
func (mc *MockController) InsertHash(hash string) error { return nil }

/* Checks that the hash validation function behaves properly */
func TestValidateHash(t *testing.T) {
	controller := &MockController{}
	if result, _ := ValidateHash(controller, testHash); !result {
		t.Log("validating hash does not work using test setup")
		t.Fail()
	}
}

/* Checks that storing and retrieving hashes behaves properly */
func TestSetHashAndRetrieveCipher(t *testing.T) {
	controller := &MockController{}
	ciphercode, err := SetHashAndRetrieveCipher(controller)
	if err != nil {
		t.Logf("set hash function has thrown an error: %s", err)
		t.Fail()
	} else if len(ciphercode) != 32 {
		t.Logf("incorrect length cipher: %s", ciphercode)
		t.Fail()
	}
}
