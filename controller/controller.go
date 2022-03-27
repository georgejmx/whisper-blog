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
 * Returns an error if applicable
**/
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
	_, err = dbo.db.Exec(query)
	return err
}

func (dbo *DatabaseObject) GrabPosts() ([]tp.Post, error) {
	var posts []tp.Post

	// Getting rows from query
	query := `select id, title, author, contents, tag, descriptors from Post`
	rows, err := dbo.db.Query(query)
	if err != nil {
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
	return posts, nil
}

func (dbo *DatabaseObject) AddPost(post tp.Post) error {
	tx, _ := dbo.db.Begin()
	_, err := tx.Exec(`insert into Post (title, author, contents, descriptors, 
		tag) values (?, ?, ?, ?, ?)`,
		post.Title, post.Author, post.Contents, post.Descriptors, post.Tag)
	if err != nil {
		return err
	}
	return tx.Commit()
}
