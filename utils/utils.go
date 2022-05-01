package utils

import (
	"bytes"
	"math/rand"
	"strings"
	"time"
)

const DAYS_INT = 86400

/* Generates a new plain-text passcode and hash pair to lead the chain */
func GenerateRawPasscode() string {
	rand.Seed(time.Now().Unix())
	var letters = []rune(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, 12)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

/* Gets the days since last post as int */
func DaysSincePost(lastPostTime time.Time) int {
	prevTime := lastPostTime.Unix()
	curTime := time.Now().Unix()
	return int((curTime - prevTime) / DAYS_INT)
}

/* Validates if the provided hash index has the authority to make a post at
this time */
func ValidateHashTiming(lastPostTime time.Time, hashIndex int) bool {
	daysElapsed := DaysSincePost(lastPostTime)

	// the next person can exclusively make a post for 5 days
	if hashIndex <= 0 && daysElapsed < 5 {
		return true
	}

	// the previous person can also make a post within a week
	if hashIndex <= 1 && daysElapsed >= 5 {
		return true
	}

	// the previous two people can also make a post within 9 days
	if hashIndex <= 2 && daysElapsed >= 7 {
		return true
	}

	// the previous 3 people can make a post within 10 days
	if hashIndex <= 3 && daysElapsed >= 9 {
		return true
	}

	// all candidate hashes, including the genesis hash, can make a post
	// after 10 days have elapsed
	if daysElapsed >= 10 {
		return true
	}

	return false
}

/* Checks if the descriptor is in the descriptors string */
func CheckDescriptor(descriptor, descriptors string) bool {
	if strings.Contains(descriptors, descriptor) &&
		!strings.Contains(descriptor, ";") {
		return true
	}
	return false
}

/* Boilerplate padding function */
func Pkcs5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

/* Boilerplate trimming function */
func Pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
