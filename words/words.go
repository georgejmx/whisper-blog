package words

import (
	"bufio"
	"bytes"
	_ "embed"
	"math/rand"
	"strings"
	"time"
)

const WORD_COUNT = 1350 // number of words in *adjectives.txt*

//go:embed adjectives.txt
var ab []byte

/* Generates a descriptors string of the format 'word;word;word;' of length 10,
to be bound to a post in the database. Gets the length then randomly picks
words from *adjectives.txt* to add to the string */
func GenerateDescriptors() (string, error) {
	var descriptors [10]string

	// Generating 10 random words using loop
	i := 0
	for i < 10 {
		descriptors[i] = GenerateDescriptor(WORD_COUNT)
		i++
	}
	return strings.Join(descriptors[:], ";"), nil
}

/* Generates a single descriptor */
func GenerateDescriptor(wordCount int) string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(wordCount)
	descriptor, err := parseWord(index)
	if err != nil {
		return "error"
	}
	return descriptor
}

/* Parses a word at the specificed random index from *adjectives.txt*.
Does so by scanning lines until the index is reached
Returns: string; the parsed word */
func parseWord(index int) (string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(ab))
	scanner.Split(bufio.ScanLines)

	scanCount := 0
	for scanner.Scan() && scanCount < index {
		scanner.Bytes()
		scanCount++
	}
	return scanner.Text(), nil
}
