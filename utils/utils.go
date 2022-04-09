package utils

import (
	"bufio"
	"bytes"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

const adjectivesPath string = "./data/adjectives.txt"

/* Generates a descriptors string of the format 'word;word;word;' of length 10,
 * to be bound to a post in the database. Gets the length then randomly picks
 * words from *adjectives.txt* to add to the string
 */
func GenerateDescriptors() (string, error) {
	var descriptors [10]string

	// Getting total number of adjectives, for use by random number generator
	wordCount, err := parseWordsCount()
	if err != nil {
		return "fail", err
	}

	// Generating 10 random words using loop
	i := 0
	rand.Seed(time.Now().UnixNano())
	for i < 10 {
		index := rand.Intn(wordCount)
		decriptor, err := parseWord(index)
		if err != nil {
			return "fail", err
		}
		descriptors[i] = decriptor
		i++
	}
	return strings.Join(descriptors[:], ";"), nil
}

/* Find the number of lines present in the specified file. Used for
 * getting the total number of possible adjectives to search through
 * that are present in *adjectives.txt*
 */
func parseWordsCount() (int, error) {
	file, err := os.Open(adjectivesPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	// Keeps reading buffers, counting line separations
	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		if err != nil && err != io.EOF {
			return 0, err
		} else if err == io.EOF {
			break
		}
	}
	return count + 1, nil
}

/* Parses a word at the specificed random index from *adjectives.txt*.
 * Does so by scanning lines until the index is reached
 * Returns: string; the parsed word
 */
func parseWord(index int) (string, error) {
	file, err := os.Open(adjectivesPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanCount := 0
	for scanner.Scan() && scanCount < index {
		scanner.Bytes()
		scanCount++
	}
	return scanner.Text(), nil
}
