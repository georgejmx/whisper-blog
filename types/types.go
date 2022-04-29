package types

import "time"

// Represents a post on the UI
type Post struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	Contents    string     `json:"contents"`
	Tag         int        `json:"tag"`
	Descriptors string     `json:"descriptors"`
	Time        time.Time  `json:"time"`
	Hash        string     `json:"hash,omitempty"`
	Reactions   []Reaction `json:"reactions,omitempty"`
}

// Represents a reaction on the UI
type Reaction struct {
	Id           int    `json:"id,omitempty"`
	PostId       int    `json:"postId,omitempty"`
	Descriptor   string `json:"descriptor"`
	Gravitas     int    `json:"gravitas"`
	GravitasHash string `json:"hash,omitempty"`
}

// A template for an object that performs database interactions
type ControllerTemplate interface {
	Init() error
	SelectPosts() ([]Post, error)
	SelectPostReactions(postId int) ([]Reaction, error)
	SelectLatestTimestamp() (time.Time, error)
	SelectLatestHash() (string, error)
	SelectCandidateHashes() ([5]string, error)
	SelectPostReactionHashes(postId int) ([5]string, error)
	SelectDescriptors(postId int) (string, error)
	SelectAnonReactionCount(postId int) (int, error)
	InsertPost(post Post) error
	InsertReaction(reaction Reaction) error
	InsertHash(hash string) error
}
