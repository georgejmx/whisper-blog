package security

import (
	"testing"
	u "whisper-blog/utils"
)

/* Checks that the hash validation function behaves properly */
func TestValidateHash(t *testing.T) {
	controller := &u.MockController{}
	if result, _ := ValidateHash(controller, u.TestHash); !result {
		t.Log("validating hash does not work using test setup")
		t.Fail()
	}
}

/* Checks that storing and retrieving hashes behaves properly */
func TestSetHashAndRetrieveCipher(t *testing.T) {
	controller := &u.MockController{}
	ciphercode, err := SetHashAndRetrieveCipher(controller)
	if err != nil {
		t.Logf("set hash function has thrown an error: %s", err)
		t.Fail()
	} else if len(ciphercode) != 32 {
		t.Logf("incorrect length cipher: %s", ciphercode)
		t.Fail()
	}
}
