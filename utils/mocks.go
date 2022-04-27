package utils

import (
	"time"
	tp "whisper-blog/types"
)

type MockController struct{}

const TestHash = "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"

/*** Defining the mock database controller methods used for all tests ***/
func (mc *MockController) Init() error {
	return nil
}

func (mc *MockController) GrabPosts() ([]tp.Post, error) {
	testPost := tp.Post{Id: 1, Title: "test", Author: "tester",
		Contents: "test contents", Tag: 0, Descriptors: "t;t;t", Time: time.Now()}
	return []tp.Post{testPost}, nil
}

func (mc *MockController) GrabLatestPosttime() (time.Time, error) {
	return time.Now(), nil
}

func (mc *MockController) AddPost(post tp.Post) error {
	return nil
}

func (mc *MockController) SelectCandidateHashes() ([5]string, error) {
	mockHashes := [5]string{
		"UNUSEDc1884c7659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a09",
		"PREVIOUS8184c7d659a2feaa0c55ad015a3bf4f1b2b0b82215d6c15b0f00a08",
		"PENULTIMATE4c7d659a2feaa0c55ad015a3bf4f1b2b0b82215d6c15b0f00a0a",
		"THIRDccc84c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a0b",
		"GENESISccc87d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a0b"}
	return mockHashes, nil
}

func (mc *MockController) SelectHash() (string, error) {
	return TestHash, nil
}

func (mc *MockController) InsertHash(hash string) error {
	return nil
}
