package utils

import (
	"strings"
	"testing"
)

/* Tests that the correct descriptors for a post are generated */
func TestGenerateDescriptors(t *testing.T) {
	adjectivesPath = "../data/adjectives.txt"
	output, err := GenerateDescriptors()
	if err != nil {
		t.Log("error generating descriptors", err)
		t.Fail()
	}

	outputs := strings.Split(output, ";")
	if len(outputs) != 10 {
		t.Log("descriptors is not 10 words conactenated by ';'")
		t.Fail()
	}
}

/* Tests that parseWord does what it needs to; the exact functionality works
for all edge cases */
func TestParseWord(t *testing.T) {
	adjectivesPath = "../data/adjectives.txt"
	output, err := parseWord(4)
	if err != nil {
		t.Log("error when opening the file: ", err)
		t.Fail()
	}

	if output != "adventurous" {
		t.Logf("expecting 'adventurous'. found %s", output)
		t.Fail()
	}

	outputZero, _ := parseWord(0)
	if outputZero != "abandoned" {
		t.Logf("expecting 'abandoned'. found %s", output)
		t.Fail()
	}

	outputEnd, _ := parseWord(1346)
	if outputEnd != "zigzag" {
		t.Logf("expecting 'zigzag'. found %s", output)
		t.Fail()
	}
}
