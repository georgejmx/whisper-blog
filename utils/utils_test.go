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
