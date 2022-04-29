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
	testServer        *httptest.Server
	testPostReqBodies = u.TestPosts
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

/* Tests that the AddPost route behaves properly*/
func TestAddPostSuccess(t *testing.T) {
	var respJson tp.Response

	// Add a genesis post
	jsonBody, err := json.Marshal(testPostReqBodies[0])
	if err != nil {
		t.Fatal("unable to marshal Post type into body")
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/data/post", testServer.URL),
		"application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal("unable to make genesis post")
	}

	respData, err := ioutil.ReadAll(resp.Body)
	err2 := json.Unmarshal(respData, &respJson)
	if err != nil || err2 != nil {
		t.Fatal("unable to parse response json into correct type")
	}

	if respJson.Marker != 2 || respJson.Data == "" {
		t.Log("did not get expected response from genesis post")
		t.Logf("marker: %d, data: %s\n", respJson.Marker, respJson.Data)
		t.Fail()
	}

	// Add another post

	// Add 3rd post

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

/* Clearing db then closing server */
func teardownAll() {
	if !r.Clear() {
		log.Fatal("unable to clear testdb")
	}
	testServer.Close()
}
