package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	r "github.com/georgejmx/whisper-blog/routes"
	x "github.com/georgejmx/whisper-blog/security"
	tp "github.com/georgejmx/whisper-blog/types"
	u "github.com/georgejmx/whisper-blog/utils"
)

var (
	testServer               *httptest.Server
	testPostReqBodies        = u.TestPosts
	testInvalidPostReqBodies = u.TestInvalidPosts
	respJson                 tp.Response
	passHashes               = []string{x.RawToHash("gen6si9")}
	invalidHashes            = u.InvalidMockHashes
)

/* Integration tests entry point */
func TestMain(m *testing.M) {
	testServer = httptest.NewServer(setup(false))
	code := m.Run()
	teardownAll()
	os.Exit(code)
}

/* Tests core server process */
func TestProcess(t *testing.T) {
	// Check that request can be made
	resp, err := http.Get(fmt.Sprintf("%s/", testServer.URL))
	if err != nil {
		t.Fatal("unable to make request")
	}

	// Check for incorrect response header or format
	_, ok := resp.Header["Content-Type"]
	if resp.StatusCode != 200 || !ok {
		t.Fatal("base request not processed correctly")
	}
}

/* Tests that the AddPost route behaves properly */
func TestAddPostSuccess(t *testing.T) {
	// add genesis post if not done already
	if len(passHashes) == 1 {
		addGenesisPost(t)
	}

	// adding 3 valid posts in succession
	i := 1
	for i < 4 {
		// Making a post with a valid request body
		testPostReqBodies[i].Hash = passHashes[len(passHashes)-1]
		jsonBody, _ := json.Marshal(testPostReqBodies[i])
		resp, err := http.Post(
			fmt.Sprintf("%s/data/post", testServer.URL),
			"application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("unable to make post with index %d", i)
		}

		// Checking for a valid response
		respData, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(respData, &respJson)
		if respJson.Marker != 1 || len(respJson.Data) != 32 {
			t.Logf("did not get expected response from %dth post", i)
			t.Logf("marker: %d, data: %s\n", respJson.Marker, respJson.Data)
			t.Log(respJson.Message)
			t.Fatal()
		}

		// Retreiving passcode, adding its hash to our hash list
		passcode, _ := decryptCipher(
			passHashes[len(passHashes)-1], respJson.Data)
		if len(passcode) != 12 {
			t.Logf("error: %v, passcode: %s", err.Error(), passcode)
			t.Fail()
		}
		passHashes = append(passHashes, x.RawToHash(passcode))
		i++
	}
}

/* Test that AddPost fails to add a post when criteria are not met */
func TestAddPostFailure(t *testing.T) {
	// add genesis post if not done already
	if len(passHashes) == 1 {
		addGenesisPost(t)
	}

	// Adding 3 invalid posts in succession
	i := 0
	for i < 3 {
		// Making post with invalid request body
		testInvalidPostReqBodies[i].Hash = passHashes[len(passHashes)-1]
		jsonBody, _ := json.Marshal(testInvalidPostReqBodies[i])
		resp, err := http.Post(
			fmt.Sprintf("%s/data/post", testServer.URL),
			"application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("unable to make invalid post with index %d", i)
		}

		// Checking for invalid response
		respData, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(respData, &respJson)
		if respJson.Marker != 0 {
			t.Logf("did not get expected response from %dth invalid post", i)
			t.Log(respJson.Message)
			t.Fatal()
		}
		i++
	}

	// Attempting to add 3 valid posts, but with empty or invalid hash
	i = 0
	for i < 3 {
		testPostReqBodies[5].Hash = invalidHashes[i]
		jsonBody, _ := json.Marshal(testPostReqBodies[5])
		resp, err := http.Post(
			fmt.Sprintf("%s/data/post", testServer.URL),
			"application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("unable to make invalid post with index %d", i+3)
		}

		// Checking for invalid response
		respData, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(respData, &respJson)
		if respJson.Marker != 0 {
			t.Logf("did not get expected response from %dth invalid post", i+3)
			t.Log(respJson.Message)
			t.Fatal()
		}
		i++
	}
}

/* Tests this utility function we need */
func TestDecryptCipher(t *testing.T) {
	prevHash := x.RawToHash("gen6si9")
	responseData := "b9f4b247240b9bc78756d4d83150a99e"
	passcode, err := decryptCipher(prevHash, responseData)
	if err != nil || len(passcode) != 12 {
		t.Logf("err: %v, passcode we got: %v\n", err, passcode)
		t.Fail()
	}
}

/* Adds a genesis post to chain. Is needed for all major tests */
func addGenesisPost(t *testing.T) {
	// Create json request body
	jsonBody, err := json.Marshal(testPostReqBodies[0])
	if err != nil {
		t.Fatal("unable to marshal Post type into body")
	}

	// Send the request
	resp, err := http.Post(
		fmt.Sprintf("%s/data/post", testServer.URL),
		"application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal("unable to make genesis post")
	}

	// Parse response into type, checking we got a success response
	respData, err := ioutil.ReadAll(resp.Body)
	err2 := json.Unmarshal(respData, &respJson)
	if err != nil || err2 != nil {
		t.Fatal("unable to parse response json into correct type")
	} else if respJson.Marker != 2 || len(respJson.Data) != 32 {
		t.Log("did not get expected response from genesis post")
		t.Logf("marker: %d, data: %s\n", respJson.Marker, respJson.Data)
		t.Fail()
	}

	// Parsing a raw passcode from the response, storign this hash
	passcode, err := decryptCipher(passHashes[0], respJson.Data)
	if err != nil || len(passcode) != 12 {
		t.Logf("error decrypting cipher: %v, passcode length: %d",
			err.Error(), len(passcode))
		t.Fail()
	}
	passHashes = append(passHashes, x.RawToHash(passcode))
}

/* For use in integration tests, also a reference for the frontend js
implementation */
func decryptCipher(prevHash, cipherStr string) (string, error) {
	cipherBytes, _ := hex.DecodeString(cipherStr)
	block, err := aes.NewCipher([]byte(prevHash[28:60]))
	output := make([]byte, len(cipherBytes))
	mode := cipher.NewCBCDecrypter(block, []byte("snooping6is9bad0"))
	mode.CryptBlocks(output, cipherBytes)
	output = u.Pkcs5Trimming(output)
	return string(output), err
}

/* Clearing db then closing server */
func teardownAll() {
	if !r.Clear() {
		log.Fatal("unable to clear testdb")
	}
	testServer.Close()
}
