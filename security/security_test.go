package security

import (
	"testing"
	mock "whisper-blog/utils"
)

/* Checks that the hash validation function behaves properly */
func TestValidateHash(t *testing.T) {
	controller := &mock.MockController{}

	// Latest hash will always succeed with no error
	isValid, err := ValidateHash(controller, mock.MockHashes[0])
	if err != nil {
		t.Logf("execution failed with error %v", err)
		t.Fail()
	} else if !isValid {
		t.Log("validating latest hash does not work using test setup")
		t.Fail()
	}

	// Penultimate Previous hash should succeed, as mock latest time > 7 days
	isValid, err = ValidateHash(controller, mock.MockHashes[2])
	if err != nil {
		t.Logf("execution failed with error %v", err)
		t.Fail()
	} else if !isValid {
		t.Log("validating previous hash does not work using test setup")
		t.Fail()
	}
}

/* Checks that the hash validation function fails when expected */
func TestValidateHashFailure(t *testing.T) {
	controller := &mock.MockController{}

	// Checks that an invalid hash returns false with correct error
	isValid, err := ValidateHash(controller, mock.InvalidMockHash)
	errMsg := err.Error()
	if isValid || string(errMsg[0]) != "a" {
		t.Log("validating invalid hash succeeded")
		t.Fail()
	}

	// Checks that a valid hash with invalid time returns false and correct msg
	for i := 0; i < 2; i++ {
		isValid, err = ValidateHash(controller, mock.MockHashes[i+3])
		errMsg = err.Error()
		if isValid || string(errMsg[0]) != "b" {
			t.Logf("validating hash number %d with wrong time succeeded", i)
			t.Fail()
		}
	}
}

/* Checks that storing and retrieving hashes behaves properly */
func TestSetHashAndRetrieveCipher(t *testing.T) {
	controller := &mock.MockController{}
	ciphercode, err := SetHashAndRetrieveCipher(controller)
	if err != nil {
		t.Logf("set hash function has thrown an error: %s", err)
		t.Fail()
	} else if len(ciphercode) != 32 {
		t.Logf("incorrect length cipher: %s", ciphercode)
		t.Fail()
	}
}
