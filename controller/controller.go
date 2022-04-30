package controller

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	tp "github.com/georgejmx/whisper-blog/types"

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
	dbo.db, err = sql.Open("sqlite3", os.Getenv("DB_FILEPATH"))
	if err != nil {
		return err
	}
	dbo.db.SetConnMaxLifetime(time.Minute * 2)
	dbo.db.SetMaxOpenConns(10)
	dbo.db.SetMaxIdleConns(10)

	// Define database tables
	queries := [3]string{
		`create table if not exists Post (
			id integer primary key autoincrement not null,
			title varchar(40) not null unique,
			author varchar(10),
			contents varchar(1500) not null,
			tag integer not null,
			descriptors varchar(210),
			time datetime default current_timestamp,
			check (tag >= 0 and tag < 8)
		)`,
		`create table if not exists Passcode (
			id integer primary key autoincrement not null,
			hash varchar(64) not null
		)`,
		`create table if not exists Reaction (
			id integer primary key autoincrement not null,
			postId integer not null,
			descriptor varchar(20) not null,
			gravitas integer not null,
			gravitasHash varchar(64),
			foreign key(postId) references Post(id),
			check (gravitas <= 6)
		)`,
	}

	// Execute all table creation on database
	tx, _ := dbo.db.Begin()
	for _, query := range queries {
		_, err = tx.Exec(query)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

/* Gets all Post tuples from sqlite */
func (dbo *DbController) SelectPosts() ([]tp.Post, error) {
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

/* Gets Reaction tuples from sqlite grouped by each descriptor. Returns a slice
with an ascending list of such tuples ordered by their total gravitas */
func (dbo *DbController) SelectPostReactions(postId int) ([]tp.Reaction, error) {
	var reactions []tp.Reaction
	tx, _ := dbo.db.Begin()

	// Getting rows from query
	rows, err := tx.Query(`select descriptor, sum(gravitas) total_gravitas from 
		Reaction where postId = ? group by descriptor`, postId)
	if err != nil {
		tx.Rollback()
		return reactions, err
	}

	// Adding reaction rows from database table to the
	for rows.Next() {
		var reaction tp.Reaction
		err = rows.Scan(&reaction.Descriptor, &reaction.Gravitas)
		if err != nil {
			return reactions, err
		}
		reactions = append(reactions, reaction)
	}

	rows.Close()
	return reactions, tx.Commit()
}

/* Gets the timestamp of the latest post */
func (dbo *DbController) SelectLatestTimestamp() (time.Time, error) {
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
and a generated descriptors and tag. // Params: Post with the fields title,
author, contents, descriptors, tag, codeHash populated //
Ensures that all data for a post has been entered */
func (dbo *DbController) InsertPost(post tp.Post) error {
	tx, _ := dbo.db.Begin()
	_, err := tx.Exec(`insert into Post (title, author, contents, descriptors, 
		tag) values (?, ?, ?, ?, ?)`,
		post.Title, post.Author, post.Contents, post.Descriptors, post.Tag)
	fmt.Printf("author: %s\n", post.Author)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

/* Adds a new reaction to db */
func (dbo *DbController) InsertReaction(reaction tp.Reaction) error {
	tx, _ := dbo.db.Begin()
	_, err := tx.Exec(`insert into Reaction (postId, descriptor, gravitas
		gravitasHash) values (?, ?, ?, ?)`, reaction.PostId,
		reaction.Descriptor, reaction.Gravitas, reaction.GravitasHash)
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

/* Selects the reaction hashes associated with a given post */
func (dbo *DbController) SelectPostReactionHashes(postId int) ([5]string, error) {
	reactionHashes := [5]string{"", "", "", "", ""}
	tx, _ := dbo.db.Begin()

	// Selecting all such hashes
	rows, err := tx.Query(`select distinct gravitasHash from Reaction where 
		postId = ?`, postId)
	if err != nil {
		tx.Rollback()
		return reactionHashes, err
	}
	i := 0
	for rows.Next() && i < 4 {
		var reactionHash string
		if err := rows.Scan(&reactionHash); err != nil {
			return reactionHashes, err
		}
		reactionHashes[i] = reactionHash
		i++
	}
	rows.Close()
	return reactionHashes, tx.Commit()
}

/* Selects the descriptors string from the post with id *postId* */
func (dbo *DbController) SelectDescriptors(postId int) (string, error) {
	var descriptors string

	tx, _ := dbo.db.Begin()
	row, err := tx.Query(`select descriptors from Post where id = ?`, postId)
	if err != nil {
		tx.Rollback()
		return descriptors, err
	}
	row.Next()
	if err = row.Scan(&descriptors); err != nil {
		return descriptors, err
	}
	return descriptors, tx.Commit()
}

/* Selects the number of anonymous reactions (those with gravitas 2) made on
the specified post */
func (dbo *DbController) SelectAnonReactionCount(postId int) (int, error) {
	var count int

	tx, _ := dbo.db.Begin()
	row, err := tx.Query(`select count(*) from Reaction where gravitas = 2 
		and postId = ?`, postId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	row.Next()
	if err = row.Scan(&count); err != nil {
		return count, err
	}
	return count, tx.Commit()
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

/* Clears db, for use in integration tests */
func (dbo *DbController) Clear() bool {
	queries := [3]string{`drop table Passcode`, `drop table Reaction`,
		`drop table Post`}

	// Execute all table creation on database
	tx, _ := dbo.db.Begin()
	for _, query := range queries {
		_, err := tx.Exec(query)
		if err != nil {
			tx.Rollback()
			return false
		}
	}
	dbo.db.Close()
	return tx.Commit() == nil
}
