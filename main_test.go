package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	r "github.com/georgejmx/whisper-blog/routes"
	x "github.com/georgejmx/whisper-blog/security"
	tp "github.com/georgejmx/whisper-blog/types"
	u "github.com/georgejmx/whisper-blog/utils"
)

type PostResponse struct {
	Message string `json:"message"`
	Marker  int    `json:"marker"`
	Data    string `json:"data"`
}

type GetResponse struct {
	Marker    int       `json:"marker"`
	DaysSince int       `json:"days_since"`
	Chain     []tp.Post `json:"chain"`
}

var (
	respJson                 PostResponse
	testServer               *httptest.Server
	testPostReqBodies        = u.TestPosts
	testInvalidPostReqBodies = u.TestInvalidPosts
	passHashes               = []string{x.RawToHash("gen6si9")}
	invalidHashes            = u.InvalidMockHashes
	hasMaxAnonHash           = false
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

	// adding 4 valid posts in succession
	i := 1
	for i < 5 {
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
		passcode, _ := x.DecryptCipher(
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
		jsonBody, _ := json.Marshal(testPostReqBodies[6])
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

/* Checks that 6 anonymous reactions can be placed but not 7 on latest post.
In the process checks that getting the chain works as expected */
func TestAddAnonReaction(t *testing.T) {
	var chainResp GetResponse
	if len(passHashes) == 1 {
		addGenesisPost(t)
	}

	// Getting the chain, needed to find correct descriptors
	resp, err := http.Get(fmt.Sprintf("%s/data", testServer.URL))
	if err != nil {
		t.Fatal("unable to get chain")
	}
	respData, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respData, &chainResp)
	if chainResp.Marker != 1 || chainResp.DaysSince != 0 {
		t.Fatal("unable to get correct chain format from database")
	}

	// Getting an array of descriptors that could be valid reactions
	latestPost := chainResp.Chain[0]
	lastPostId := latestPost.Id
	descriptors := strings.Split(latestPost.Descriptors, ";")

	// Expecting 6 successes, then a failure. Validates behaviour
	i := 0
	for i < 6 {
		addReaction(true, t, lastPostId, descriptors[i], "")
		i++
	}
	addReaction(false, t, lastPostId, descriptors[9], "")

	// Invalid hash should definitely fail
	addReaction(false, t, lastPostId, descriptors[9], invalidHashes[0])

	hasMaxAnonHash = true
}

/* Checks that the previous 3 hashes on the chain can be used to send a
single reaction independent of the anonymous ones */
func TestAddSignedReaction(t *testing.T) {
	var chainResp GetResponse
	if len(passHashes) < 5 {
		TestAddPostSuccess(t)
	} else if !hasMaxAnonHash {
		TestAddAnonReaction(t)
	}

	// Getting the chain, needed to find correct descriptors
	resp, err := http.Get(fmt.Sprintf("%s/data", testServer.URL))
	if err != nil {
		t.Fatal("unable to get chain")
	}
	respData, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respData, &chainResp)
	if chainResp.Marker != 1 || chainResp.DaysSince != 0 {
		t.Fatal("unable to get correct chain format from database")
	}

	// Getting an array of descriptors that could be valid reactions
	latestPost := chainResp.Chain[0]
	lastPostId := latestPost.Id
	descriptors := strings.Split(latestPost.Descriptors, ";")
	maxInd := len(passHashes) - 1

	// Using the genesis hash should work once
	addReaction(true, t, lastPostId, descriptors[9], passHashes[1])
	addReaction(false, t, lastPostId, descriptors[9], passHashes[1])

	// Previous and penultimate hash should work once
	i := 1
	for i < 4 {
		addReaction(true, t, lastPostId, descriptors[9], passHashes[maxInd-i])
		addReaction(false, t, lastPostId, descriptors[9], passHashes[maxInd-i])
		i++
	}

	// Latest hash should fail
	addReaction(false, t, lastPostId, descriptors[9], passHashes[maxInd])
}

/* Adds a test reaction */
func addReaction(
	isValid bool, t *testing.T, postId int, descriptor, hash string) {
	reaction := tp.Reaction{PostId: postId, Descriptor: descriptor}
	if hash != "" {
		reaction.GravitasHash = hash
	}
	body, _ := json.Marshal(reaction)
	resp, err := http.Post(fmt.Sprintf("%s/data/react", testServer.URL),
		"application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Log("unable to make reaction")
		t.Fail()
	}

	// Parse response into type, checking we got success
	respData, err := ioutil.ReadAll(resp.Body)
	err2 := json.Unmarshal(respData, &respJson)
	if err != nil || err2 != nil {
		t.Fatal("unable to parse response json into correct type")
	}

	// If was valid reaction, fail if response body indicates failure
	if isValid {
		if respJson.Marker != 1 || len(respJson.Data) == 0 {
			t.Log("seemingly valid reaction produced failure response")
			t.Log(respJson.Message)
			t.Fail()
		}
	} else {
		if respJson.Marker != 0 {
			t.Log("seemingly invalid reaction produced success response")
			t.Log(respJson.Message)
			t.Fail()
		}
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
	passcode, err := x.DecryptCipher(passHashes[0], respJson.Data)
	if err != nil || len(passcode) != 12 {
		t.Logf("error decrypting cipher: %v, passcode length: %d",
			err.Error(), len(passcode))
		t.Fail()
	}
	passHashes = append(passHashes, x.RawToHash(passcode))
}

/* Clearing db then closing server */
func teardownAll() {
	if !r.Clear() {
		log.Fatal("unable to clear testdb")
	}
	testServer.Close()
}
