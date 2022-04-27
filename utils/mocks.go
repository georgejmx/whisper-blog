package utils

import (
	"time"
	tp "whisper-blog/types"
)

// Mock database controller, used by all unit tests
type MockController struct{}

var (
	MockHashes = [5]string{
		"UNUSEDc1884c7659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a09",
		"PREVIOUS8184c7d659a2feaa0c55ad015a3bf4f1b2b0b82215d6c15b0f00a08",
		"PENULTIMATE4c7d659a2feaa0c55ad015a3bf4f1b2b0b82215d6c15b0f00a0a",
		"THIRDccc84c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a0b",
		"GENESISccc87d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a0b"}
	InvalidMockHash = `INVALIDc884c7659a2feaa0c55ad015a3bf4f1b2b0b822cd1
		5d6c15b0f00a09`
)

// Mock method implementation
func (mc *MockController) Init() error {
	return nil
}

// Mock method implementation
func (mc *MockController) GrabPosts() ([]tp.Post, error) {
	testPost := tp.Post{Id: 1, Title: "test", Author: "tester",
		Contents: "test contents", Tag: 0, Descriptors: "t;t;t",
		Time: generateMockTime()}
	return []tp.Post{testPost}, nil
}

// Mock method implementation
func (mc *MockController) GrabLatestTimestamp() (time.Time, error) {
	return generateMockTime(), nil
}

// Mock method implementation
func (mc *MockController) AddPost(post tp.Post) error {
	return nil
}

// Mock method implementation
func (mc *MockController) SelectCandidateHashes() ([5]string, error) {
	return MockHashes, nil
}

// Mock method implementation
func (mc *MockController) SelectLatestHash() (string, error) {
	return MockHashes[0], nil
}

// Mock method implementation
func (mc *MockController) InsertHash(hash string) error {
	return nil
}

/* Generates a mock time, which is just over 1 week ago */
func generateMockTime() time.Time {
	hoursAgo := (-7 * 24) - 1
	return time.Now().Add(time.Duration(hoursAgo) * time.Hour)
}
