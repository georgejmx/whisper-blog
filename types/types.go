package types

import "time"

// Represents a post convertible to pretty JSON
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

// Represents the HTML data of a post on the UI
type PostHtmlContent struct {
	Colour      string
	Timestring  string
	IsSuccessor bool
	Title       string
	Contents    string
	Author      string
	Reactions   []Reaction
}

// Contains above data needed for HTML content structure
type HtmlContainer struct {
	HtmlPosts []PostHtmlContent
}

// Represents a reaction in JSON
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
	SelectCandidateHashes() ([5]string, error)
	SelectPostReactionHashes(postId int) ([5]string, error)
	SelectDescriptors(postId int) (string, error)
	SelectAnonReactionCount(postId int) (int, error)
	InsertPost(post Post) error
	InsertReaction(reaction Reaction) error
	InsertHash(hash string) error
	Clear() bool
}
