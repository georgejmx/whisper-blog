package utils

import "testing"

/* Tests the validate hash function, with edge cases around mock latest time */
func TestValidateHashTiming(t *testing.T) {
	t.Log(generateMockTime())

	if outcome := ValidateHashTiming(generateMockTime(), 2); !outcome {
		t.Log("penultimate previous post index failed!?!")
		t.Fail()
	}

	if outcome := ValidateHashTiming(generateMockTime(), 3); outcome {
		t.Log("third previous post index succeeded!?!")
		t.Fail()
	}
}

/* Tests that a correct passcode is generated, that are not easily recreated
using the current exact time */
func TestGenerateRawPasscode(t *testing.T) {
	passcodes := []string{"aaaaaaaaaaa!"}
	i := 0
	for i < 5 {
		passcode := GenerateRawPasscode()
		if len(passcode) != 12 {
			t.Logf("incorrect passcode parsed: %s\n", passcode)
			t.Fail()
		} else if passcodes[i] == passcode {
			t.Logf("duplicate passcodes found: %s\n", passcode)
			t.Fail()
		}
		passcodes = append(passcodes, passcode)
		i++
	}
}
