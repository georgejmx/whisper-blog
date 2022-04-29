package utils

import (
	"time"

	tp "github.com/georgejmx/whisper-blog/types"
)

// Mock database controller, used by all unit tests
type MockController struct{}

// Data to populate this mock controller
var (
	MockPost = tp.Post{
		Id:          1,
		Title:       "test",
		Author:      "tester",
		Contents:    "test contents",
		Tag:         0,
		Descriptors: "test;experimental;mocking",
		Time:        generateMockTime(),
	}
	MockReaction = tp.Reaction{
		Id:         1,
		PostId:     1,
		Descriptor: "experimental",
		Gravitas:   6,
	}
	MockReaction2 = tp.Reaction{
		Id:         2,
		PostId:     1,
		Descriptor: "mocking",
		Gravitas:   4,
	}
	MockHashes = [5]string{
		"UNUSEDc1884c7659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a09a",
		"PREVIOUS8184c7d659a2feaa0c55ad015a3bf4f1b2b0b82215d6c15b0f00a08a",
		"PENULTIMATE4c7d659a2feaa0c55ad015a3bf4f1b2b0b82215d6c15b0f00a0aa",
		"THIRDccc84c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a0ba",
		"GENESISccc87d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a0ba"}
	InvalidMockHash = `INVALIDc884c7659a2feaa0c55ad015a3bf4f1b2b0b822cd1a
		5d6c15b0f00a09`
)

// Mock method implementation
func (mc *MockController) Init() error {
	return nil
}

// Mock method implementation
func (mc *MockController) InsertPost(post tp.Post) error {
	return nil
}

// Mock method implementation
func (mc *MockController) InsertReaction(reaction tp.Reaction) error {
	return nil
}

// Mock method implementation
func (mc *MockController) SelectPosts() ([]tp.Post, error) {
	return []tp.Post{MockPost}, nil
}

// Mock method implementation
func (mc *MockController) SelectPostReactions(postId int) ([]tp.Reaction, error) {
	return []tp.Reaction{MockReaction, MockReaction2}, nil
}

// Mock method implementation
func (mc *MockController) InsertHash(hash string) error {
	return nil
}

// Mock method implementation
func (mc *MockController) SelectLatestTimestamp() (time.Time, error) {
	return generateMockTime(), nil
}

// Mock method implementation
func (mc *MockController) SelectLatestHash() (string, error) {
	return MockHashes[0], nil
}

// Mock method implementation
func (mc *MockController) SelectCandidateHashes() ([5]string, error) {
	return MockHashes, nil
}

// Mock method implementation
func (mc *MockController) SelectPostReactionHashes(
	postId int) ([5]string, error) {
	return [5]string{MockHashes[1], MockHashes[3], "", "", ""}, nil
}

// Mock method implementation
func (mc *MockController) SelectDescriptors(postId int) (string, error) {
	if postId == 1 {
		return MockPost.Descriptors, nil
	} else {
		return "", nil
	}
}

// Mock method implementation
func (mc *MockController) SelectAnonReactionCount(postId int) (int, error) {
	if postId == 1 {
		return 2, nil
	} else {
		return 0, nil
	}
}

/* Generates a mock time, which is just over 1 week ago */
func generateMockTime() time.Time {
	hoursAgo := (-7 * 24) - 1
	return time.Now().Add(time.Duration(hoursAgo) * time.Hour)
}
