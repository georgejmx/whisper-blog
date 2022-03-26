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
}
