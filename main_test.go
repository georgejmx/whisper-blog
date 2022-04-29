package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	r "github.com/georgejmx/whisper-blog/routes"
)

var testServer *httptest.Server

/* Integration tests entry point */
func TestMain(m *testing.M) {
	testServer = httptest.NewServer(setup(false))
	code := m.Run()
	teardownAll()
	os.Exit(code)
}

// Tests core server process
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

/* Clearing db then closing server */
func teardownAll() {
	if !r.Clear() {
		log.Fatal("unable to clear testdb")
	}
	testServer.Close()
}
