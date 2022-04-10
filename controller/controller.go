package controller

import (
	"database/sql"
	tp "whisper-blog/types"

	_ "github.com/mattn/go-sqlite3"
)

// Database object is where all queries are executed on
type DatabaseObject struct {
	db *sql.DB
}

/* Establishes database connection and sets up tables if needed
Returns an error if applicable */
func (dbo *DatabaseObject) Init() error {
	var err error = nil

	// Open database connection
	dbo.db, err = sql.Open("sqlite3", "./data/blog.db")
	if err != nil {
		return err
	}

	// Configure database tables correctly
	query := `create table if not exists Post (
		id integer primary key autoincrement not null,
		title varchar(40) not null unique,
		author varchar(10),
		contents varchar(1500) not null,
		tag integer not null,
		descriptors varchar(210),
		time datetime default CURRENT_TIMESTAMP,
		check (tag >= 0 and tag < 8)
	)`
	if _, err = dbo.db.Exec(query); err != nil {
		return err
	}

	query2 := `create table if not exists Passcode (
		id integer primary key autoincrement not null,
		hash varchar(64) not null
	)`
	_, err = dbo.db.Exec(query2)
	return err
}

// TODO: make this into GrabData at a later date, to include reactions
func (dbo *DatabaseObject) GrabPosts() ([]tp.Post, error) {
	var posts []tp.Post
	tx, _ := dbo.db.Begin()

	// Getting rows from query
	rows, err := tx.Query(`select id, title, author, contents, tag, descriptors 
		from Post`)
	if err != nil {
		tx.Rollback()
		return posts, err
	}

	// Adding post rows from database table to the posts variable, unless error
	for rows.Next() {
		var post tp.Post
		if err2 := rows.Scan(&post.Id, &post.Title, &post.Author,
			&post.Contents, &post.Tag, &post.Descriptors); err2 != nil {
			return posts, err2
		}
		posts = append(posts, post)
	}

	rows.Close()
	return posts, tx.Commit()
}

/* Adds a new post to the blog database. Using user data from the frontend
and a generated descriptors and tag
Params: Post with the fields title, author, contents, descriptors, tag,
	codeHash populated
Ensures that all data for a post has been entered */
func (dbo *DatabaseObject) AddPost(post tp.Post) error {
	tx, _ := dbo.db.Begin()
	_, err := tx.Exec(`insert into Post (title, author, contents, descriptors, 
		tag) values (?, ?, ?, ?, ?)`,
		post.Title, post.Author, post.Contents, post.Descriptors, post.Tag)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

/* Selects the latest passcode hash from the database, for use in validation */
// TODO: Paramerise to grab the most recent 3 hashes and genesis hash
func (dbo *DatabaseObject) SelectHash() (string, error) {
	var hash string

	tx, _ := dbo.db.Begin()
	row, err := tx.Query(`select hash from Passcode where id = 
		(select max(id) from Passcode)`)
	if err != nil {
		tx.Rollback()
		return hash, err
	}
	row.Next()
	if err = row.Scan(&hash); err != nil {
		return hash, err
	}
	return hash, tx.Commit()
}

func (dbo *DatabaseObject) InsertHash(hash string) error {
	tx, _ := dbo.db.Begin()
	_, err := tx.Exec(`insert into Passcode (hash) values (?)`, hash)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
