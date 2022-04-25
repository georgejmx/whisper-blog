package types

// Represents a post on the UI
type Post struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Contents    string `json:"contents"`
	Tag         int    `json:"tag"`
	Descriptors string `json:"descriptors"`
	Time        string `json:"time"`
	Hash        string `json:"hash,omitempty"`
}

// A template for an object that performs database interactions
type ControllerTemplate interface {
	Init() error
	GrabPosts() ([]Post, error)
	AddPost(post Post) error
	SelectHash() (string, error)
	InsertHash(hash string) error
}
