package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
	tp "whisper-blog/types"
	u "whisper-blog/utils"
)

/* Function to validate the provided hash against the chain law, determining
whether a lawful post can be made */
// TODO: Encode the proper chain law e.g after 6 days an extra hash is valid
func ValidateHash(dbo tp.ControllerTemplate, hash string) (bool, error) {
	// Grabbing stored hashes and latest timestamp
	storedHashes, err := dbo.SelectCandidateHashes()
	lastPostTime, err2 := dbo.GrabLatestPosttime()
	if err != nil {
		return false, err
	} else if err2 != nil {
		return false, err2
	}

	// Validating the **Chain Law**
	hashIndex := findCandidateIndex(hash, storedHashes)
	if hashIndex == -1 {
		return false, errors.New("hash does not exist in the blog")
	}
	isValTime, err := validateHashTiming(lastPostTime, hashIndex)
	if err != nil {
		return false, errors.New("error when validing hash timing")
	} else if !isValTime {
		return false, errors.New("hash failed validation timing")
	}

	// We have a valid and correctly timed hash
	return true, nil
}

func validateHashTiming(lastPostTime time.Time, hashIndex int) (bool, error) {
	if hashIndex == 0 {
		return true, nil
	}
	return false, nil
}

/* Sets the new randomly generated hash by inserting into the database. Returns
A string which is the new raw text symmetrically encrypted */
func SetHashAndRetrieveCipher(dbo tp.ControllerTemplate) (string, error) {
	oldHash, err := dbo.SelectHash()
	if err != nil {
		return "fail", err
	}

	// Generating passcode and hash
	rawPasscode := u.GenerateRawPasscode()
	hashBytes := sha256.Sum256([]byte(rawPasscode))
	dbo.InsertHash(hex.EncodeToString(hashBytes[:]))

	// Initialising cipher with the old hash
	bPlaintext := pkcs5Padding([]byte(rawPasscode), aes.BlockSize, 12)
	block, err := aes.NewCipher([]byte(oldHash[28:60]))
	if err != nil {
		return "fail", err
	}

	// Encrypting the raw passcode for response to client
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, []byte("snooping6is9bad0"))
	mode.CryptBlocks(ciphertext, bPlaintext)
	return hex.EncodeToString(ciphertext), nil
}

/* Boilerplate padding function */
func pkcs5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

/* Finds if the provided hash is a candidate, and if so the index */
func findCandidateIndex(providedHash string, candidates [5]string) int {
	for ind, value := range candidates {
		if value == providedHash {
			return ind
		}
	}
	return -1
}
