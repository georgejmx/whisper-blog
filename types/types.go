package types

import "time"

// Represents a post on the UI
type Post struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Contents    string    `json:"contents"`
	Tag         int       `json:"tag"`
	Descriptors string    `json:"descriptors"`
	Time        time.Time `json:"time"`
	Hash        string    `json:"hash,omitempty"`
}

// Represents a reaction on the UI
type Reaction struct {
	Id           int    `json:"id"`
	PostId       int    `json:"postId"`
	Descriptor   string `json:"descriptor"`
	Gravitas     int    `json:"gravitas"`
	GravitasHash string `json:"hash,omitempty"`
}

// A template for an object that performs database interactions
type ControllerTemplate interface {
	Init() error
	GrabPosts() ([]Post, error)
	GrabLatestTimestamp() (time.Time, error)
	SelectLatestHash() (string, error)
	SelectCandidateHashes() ([5]string, error)
	AddPost(post Post) error
	InsertHash(hash string) error
}
