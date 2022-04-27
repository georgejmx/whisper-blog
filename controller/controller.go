package controller

import (
	"database/sql"
	"time"
	tp "whisper-blog/types"

	_ "github.com/mattn/go-sqlite3"
)

// Database object is where all queries are executed on
type DbController struct {
	db *sql.DB
}

/* Establishes database connection and sets up tables if needed
Returns an error if applicable */
func (dbo *DbController) Init() error {
	var err error = nil

	// Open database connection
	dbo.db, err = sql.Open("sqlite3", "./data/blog_test.db")
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
		time datetime default current_timestamp,
		check (tag >= 0 and tag < 8)
	)`
	if _, err = dbo.db.Exec(query); err != nil {
		return err
	}

	query2 := `create table if not exists Passcode (
		id integer primary key autoincrement not null,
		hash varchar(64) not null
	)`
	if _, err = dbo.db.Exec(query2); err != nil {
		return err
	}

	query3 := `create table if not exists Reaction (
		id integer primary key autoincrement not null,
		postId integer not null,
		descriptor varchar(20) not null,
		gravitas integer not null,
		foreign key(postId) references Post(id)
	)`
	_, err = dbo.db.Exec(query3)
	return err
}

// TODO: make this into GrabData at a later date, to include reactions
func (dbo *DbController) GrabPosts() ([]tp.Post, error) {
	var posts []tp.Post
	tx, _ := dbo.db.Begin()

	// Getting rows from query
	rows, err := tx.Query(`select * from Post order by id desc`)
	if err != nil {
		tx.Rollback()
		return posts, err
	}

	// Adding post rows from database table to the posts variable, unless error
	for rows.Next() {
		var post tp.Post
		if err = rows.Scan(&post.Id, &post.Title, &post.Author, &post.Contents,
			&post.Tag, &post.Descriptors, &post.Time); err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	rows.Close()
	return posts, tx.Commit()
}

func (dbo *DbController) GrabPostReactions(postId int) ([]tp.Reaction, error) {
	var reactions []tp.Reaction
	tx, _ := dbo.db.Begin()

	// Getting rows from query
	rows, err := tx.Query(`select descriptor, sum(gravitas) total_gravitas from 
		Reaction group by descriptor`, postId)
	if err != nil {
		tx.Rollback()
		return reactions, err
	}

	// Adding reaction rows from database table to the
	for rows.Next() {
		var reaction tp.Reaction
		if err = rows.Scan(&reaction.Descriptor, &reaction.Gravitas); err != nil {
			return reactions, err
		}
		reactions = append(reactions, reaction)
	}

	rows.Close()
	return reactions, tx.Commit()
}

/* Gets the timestamp of the latest post */
func (dbo *DbController) GrabLatestTimestamp() (time.Time, error) {
	var timestamp time.Time

	tx, _ := dbo.db.Begin()
	row, err := tx.Query(`select time from Post where id = 
		(select max(id) from Post)`)
	if err != nil {
		tx.Rollback()
		return timestamp, err
	}
	row.Next()
	if err = row.Scan(&timestamp); err != nil {
		return timestamp, err
	}
	return timestamp, tx.Commit()
}

/* Adds a new post to the blog database. Using user data from the frontend
and a generated descriptors and tag
Params: Post with the fields title, author, contents, descriptors, tag,
	codeHash populated
Ensures that all data for a post has been entered */
func (dbo *DbController) AddPost(post tp.Post) error {
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

/* Adds a new reaction to db */
func (dbo *DbController) AddReaction(reaction tp.Reaction) error {
	tx, _ := dbo.db.Begin()
	_, err := tx.Exec(`insert into Reaction (postId, descriptor, gravitas) 
		values (?, ?, ?)`, reaction.PostId, reaction.Descriptor,
		reaction.Gravitas)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

/* Selects the latest passcode hash from the database, for use in validation */
func (dbo *DbController) SelectLatestHash() (string, error) {
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

/* Selects the 5 hashes that can be used for post or reaction validation. This
is an array of the form [latest hash, second latest hash, third latest,
fourth latest, genesis hash] */
func (dbo *DbController) SelectCandidateHashes() ([5]string, error) {
	hashes := [5]string{"", "", "", "", ""}
	tx, _ := dbo.db.Begin()

	// Selecting the most recent 4 hashes with such query, then parsing
	topRows, err := tx.Query(
		`select hash from Passcode order by id desc limit 4`)
	if err != nil {
		tx.Rollback()
		return hashes, err
	}
	i := 0
	for topRows.Next() && i < 4 {
		if err = topRows.Scan(&hashes[i]); err != nil {
			return hashes, err
		}
		i++
	}
	topRows.Close()

	// Selecting the genesis row, then returning the complete array
	genesisRow, err := tx.Query(
		`select hash from Passcode order by id asc limit 1`)
	if err != nil {
		tx.Rollback()
		return hashes, err
	}
	genesisRow.Next()
	if err = genesisRow.Scan(&hashes[4]); err != nil {
		return hashes, err
	}
	return hashes, tx.Commit()
}

/* Adds a new row to the passcode table, with a generated hash */
func (dbo *DbController) InsertHash(hash string) error {
	tx, _ := dbo.db.Begin()
	_, err := tx.Exec(`insert into Passcode (hash) values (?)`, hash)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
