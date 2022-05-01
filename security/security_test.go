package security

import (
	"os"
	"testing"

	mock "github.com/georgejmx/whisper-blog/utils"
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
	isValid, err := ValidateHash(controller, mock.InvalidMockHashes[0])
	errMsg := err.Error()
	if isValid || string(errMsg[0]) != "a" {
		t.Log("validating invalid hash succeeded")
		t.Fail()
	}

	// Checks that an emptyhash returns false with correct error
	isValid, err = ValidateHash(controller, "")
	errMsg = err.Error()
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

/* Checks that hash validation for hashes behaves properly */
func TestValidateReactionHash(t *testing.T) {
	controller := &mock.MockController{}

	// Trying an unused candidate hash
	isValid, gravitas, err := ValidateReactionHash(
		controller, mock.MockHashes[2], 1)
	if err != nil || gravitas != 6 || !isValid {
		t.Log("expected no error and gravitas=6 from unused candidate hash")
		t.Fail()
	}

	// Trying the genesis hash (unused)
	isValid, gravitas, err = ValidateReactionHash(
		controller, mock.MockHashes[4], 1)
	if err != nil || gravitas != 1 || !isValid {
		t.Log("expected no error and gravitas=6 from unused candidate hash")
		t.Fail()
	}

	// Checks that when hash is empty, returns isValid=false but no error
	isValid, gravitas, err = ValidateReactionHash(controller, "", 1)
	if err != nil || gravitas != 2 || isValid {
		t.Logf("gravitas=%v, isValid=%v, err=%s\n", gravitas, isValid, err)
		t.Fail()
	}

}

/* Checks that hash validation for hashes fails when expected */
func TestValidateReactionHashFails(t *testing.T) {
	controller := &mock.MockController{}

	// Checks that attempting to use a hash twice fails
	isValid, _, err := ValidateReactionHash(controller, mock.MockHashes[1], 1)
	if err == nil || isValid {
		t.Log("expected an error for an already used hash")
		t.Fail()
	}

	// Checks that attempting to use a hash twice fails again
	isValid, _, err = ValidateReactionHash(controller, mock.MockHashes[3], 1)
	if err == nil || isValid {
		t.Log("expected an error for an already used hash")
		t.Fail()
	}

	// Checks that attempting react on your own post fails
	isValid, _, err = ValidateReactionHash(controller, mock.MockHashes[0], 1)
	if err == nil || isValid {
		t.Log("expected an error for an already used hash")
		t.Fail()
	}
}

/* Checks that storing and retrieving hashes behaves properly for both genesis
hash and also future posts*/
func TestSetHashAndRetrieveCipher(t *testing.T) {
	// Ensuring that required environment variables are set for tests
	os.Setenv("AES_SPLICE_INDEX", "28")
	os.Setenv("AES_IV", "snooping6is9bad0")

	controller := &mock.MockController{}
	ciphercode, err := SetHashAndRetrieveCipher(
		controller, false, mock.MockHashes[0])
	if err != nil {
		t.Logf("set hash function has thrown an error: %s", err)
		t.Fail()
	} else if len(ciphercode) != 32 || len(ciphercode) == 0 {
		t.Logf("incorrect length cipher: %s", ciphercode)
		t.Fail()
	}

	passcode, err := DecryptCipher(mock.MockHashes[0], ciphercode)
	if err != nil {
		t.Logf("decrypting cipher threw an error: %s", err)
		t.Fail()
	} else if len(passcode) != 12 || len(passcode) == 0 {
		t.Logf("incorrect passcode: %s", passcode)
		t.Fail()
	}

	// Genesis post case
	ciphercode, err = SetHashAndRetrieveCipher(controller, true, "")
	if err != nil {
		t.Logf("set hash function has thrown an error at genesis: %s", err)
		t.Fail()
	} else if len(ciphercode) != 32 || ciphercode == "" {
		t.Logf("incorrect length cipher: %s", ciphercode)
		t.Fail()
	}

	passcode, err = DecryptCipher(RawToHash("gen6si9"), ciphercode)
	if err != nil {
		t.Logf("decrypting cipher threw an error: %s", err)
		t.Fail()
	} else if len(passcode) != 12 || len(passcode) == 0 {
		t.Logf("incorrect passcode: %s", passcode)
		t.Fail()
	}
}
