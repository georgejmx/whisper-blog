package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"strconv"

	tp "github.com/georgejmx/whisper-blog/types"
	u "github.com/georgejmx/whisper-blog/utils"
)

/* Function to validate the provided hash against the **Chain Law**, determining
whether a lawful post can be made */
func ValidateHash(dbo tp.ControllerTemplate, hash string) (bool, error) {
	// Grabbing stored hashes and latest timestamp
	storedHashes, err := dbo.SelectCandidateHashes()
	lastPostTime, err2 := dbo.SelectLatestTimestamp()
	if err != nil {
		return false, err
	} else if err2 != nil {
		return false, err2
	}

	// Validating the Chain Law
	hashIndex := findHashIndex(hash, storedHashes)
	if hashIndex == -1 {
		return false, errors.New("a: hash will never have ability to make post")
	}
	if isValTime := u.ValidateHashTiming(lastPostTime, hashIndex); !isValTime {
		return false, errors.New("b: hash failed validation timing")
	}

	// We have a valid and correctly timed hash
	return true, nil
}

/* Function to validate a reaction hash against the **Chain Law**, determining
what level of gravitas the reaction will have */
func ValidateReactionHash(
	dbo tp.ControllerTemplate, hash string, postId int) (bool, int, error) {

	// Checking for a null hash
	if len(hash) < 64 {
		return false, 2, nil
	}

	// Performing db operations
	storedHashes, err := dbo.SelectCandidateHashes()
	postReactionHashes, err2 := dbo.SelectPostReactionHashes(postId)
	if err != nil || err2 != nil {
		return false, 2, err
	}

	// Finding whether this hash has already been validated
	if findHashIndex(hash, postReactionHashes) >= 0 {
		return false, 2, errors.New("hash has already reacted")
	}

	// Finding if this hash is a candidate, and determining gravitas
	candidateHashIndex := findHashIndex(hash, storedHashes)
	if candidateHashIndex == -1 {
		return false, 2, nil
	} else if candidateHashIndex == 0 {
		return false, 0, errors.New(
			"you do not have gravitas to react on your own post")
	} else if candidateHashIndex == 4 {
		return true, 1, nil
	}

	// We have an unused hash of gravitas 6, with no errors
	return true, 6, nil
}

/* Sets the new randomly generated hash by inserting into the database. Returns
A string which is the new raw text symmetrically encrypted */
func SetHashAndRetrieveCipher(dbo tp.ControllerTemplate) (string, error) {
	spliceInd, _ := strconv.ParseInt(os.Getenv("AES_SPLICE_INDEX"), 10, 64)
	oldHash, err := dbo.SelectLatestHash()
	if err != nil {
		return "fail", err
	}

	// Generating passcode and hash
	rawPasscode := u.GenerateRawPasscode()
	hashBytes := sha256.Sum256([]byte(rawPasscode))
	dbo.InsertHash(hex.EncodeToString(hashBytes[:]))

	// Initialising cipher with the old hash
	bPlaintext := pkcs5Padding([]byte(rawPasscode), aes.BlockSize, 12)
	block, err := aes.NewCipher([]byte(oldHash[spliceInd : spliceInd+32]))
	if err != nil {
		return "fail", err
	}

	// Encrypting the raw passcode for response to client
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, []byte(os.Getenv("AES_IV")))
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
func findHashIndex(providedHash string, candidates [5]string) int {
	for ind, value := range candidates {
		if value == providedHash {
			return ind
		}
	}
	return -1
}
