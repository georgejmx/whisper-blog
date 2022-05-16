package utils

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const DAYS_INT = 86400

/* Generates a new plain-text to lead the chain, that is moderately secure
in plaintext */
func GenerateRawPasscode() string {
	var letters = []rune(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, 12)
	for i := range s {
		bigrand, _ := rand.Int(rand.Reader, big.NewInt(9999999999))
		smallrand, _ := strconv.ParseInt(bigrand.Text(10), 10, 0)
		smallrand += time.Now().UnixMilli()
		s[i] = letters[smallrand%62]
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

/* Gets the tailwind CSS colour of a tag */
func GetTagColour(tag int) string {
	switch tag {
	case 0:
		return "bg-slate-300"
	case 1:
		return "bg-orange-300"
	case 2:
		return "bg-indigo-300"
	case 3:
		return "bg-pink-300"
	case 4:
		return "bg-cyan-300"
	case 5:
		return "bg-lime-300"
	case 6:
		return "bg-orange-100"
	case 7:
		return "bg-purple-300"
	default:
		return "bg-red-500"
	}
}

/* Gets a UI suitable time for the post */
func GetTimestring(moment time.Time) string {
	rfc := moment.Format(time.RFC1123)
	return rfc[0:16]
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
