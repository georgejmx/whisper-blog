package controller

import (
	"database/sql/driver"
	"fmt"
	"testing"
	"time"
	tp "whisper-blog/types"

	"github.com/DATA-DOG/go-sqlmock"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

var (
	testDbo DbController
	mock    sqlmock.Sqlmock
	err     error
)

/* Called at the begining of each test; sets up stub db */
func setupTest(t *testing.T) {
	testDbo.db, mock, err = sqlmock.New()
	if err != nil {
		t.Log("error when opening stub database")
		t.Fail()
	}
}

/* Tests that the GrabPost controller behaves as expected. Ensures that the SQL
driver correctly processes grabbing posts */
func TestSelectPostsSuccess(t *testing.T) {
	setupTest(t)

	// Needed for correct query matching of nested query with time type
	testDbo.db, mock, err = sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	// Mocking db operations by populating this mock database
	headers := []string{"id", "title", "author", "contents", "tag",
		"descriptors", "time"}
	rows := sqlmock.NewRows(headers).
		AddRow(1, "test title", "tester", "testing is so cool", 3,
			"t;t;t;t", time.Now()).
		AddRow(2, "test title", "tester 2", "bruh", 4, "t;t;t;t", time.Now())

	mock.ExpectBegin()
	mock.ExpectQuery(`select * from Post order by id desc`).WillReturnRows(rows)
	mock.ExpectCommit()

	// Running the real function with above parameters
	if _, err = testDbo.SelectPosts(); err != nil {
		t.Logf("error not expected when grabbing posts: %s", err)
		t.Fail()
	}
	teardownTest(t)
}

/* Tests that the GrabPost controller behaves as expected. Ensures that the SQL
driver correctly behaves in the case of an error */
func TestSelectPostsFailure(t *testing.T) {
	setupTest(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`select * from Post order by id desc`).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	// Running the real function with above parameters
	if _, err = testDbo.SelectPosts(); err == nil {
		t.Logf("expecting error when grabbing posts: %s", err)
		t.Fail()
	}
	testDbo.db.Close()
}

/* Tests that selecting a hash is a valid query and does whats expected */
func TestSelectLatestHash(t *testing.T) {
	setupTest(t)

	// Needed for correct query matching of nested queries
	testDbo.db, mock, err = sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	latestHash := `9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c
		15b0f00a08`
	rows := sqlmock.NewRows([]string{"hash"}).AddRow(latestHash)

	mock.ExpectBegin()
	mock.ExpectQuery(`select hash from Passcode where id = 
	(select max(id) from Passcode)`).WillReturnRows(rows)
	mock.ExpectCommit()

	if _, err = testDbo.SelectLatestHash(); err != nil {
		t.Logf("error not expected when selecting hash: %s", err)
		t.Fail()
	}
	teardownTest(t)
}

/* Tests that selecting candidate hashes works as expected*/
func TestSelectCandidateHashes(t *testing.T) {
	setupTest(t)

	testDbo.db, mock, err = sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mock.ExpectBegin()

	// Testing getting recent hashes
	sampleHashes := [4]string{
		"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a09",
		"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a0a",
		"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a0b",
	}
	rows := sqlmock.NewRows([]string{"hash"}).AddRow(sampleHashes[0]).
		AddRow(sampleHashes[1]).AddRow(sampleHashes[2]).
		AddRow(sampleHashes[3])
	mock.ExpectQuery(`select hash from Passcode order by id desc limit 4`).
		WillReturnRows(rows)

	// Testing getting genesis hash
	rows2 := sqlmock.NewRows([]string{"hash"}).
		AddRow("9f86d081884c7d659a2feaa055ad015a3bf4f1b2b0b822cd15d6c15b0f00a0bc")
	mock.ExpectQuery(`select hash from Passcode order by id asc limit 1`).
		WillReturnRows(rows2)

	// Tests that these hashes are correctly sandwiched together
	mock.ExpectCommit()
	testHashes, err := testDbo.SelectCandidateHashes()
	if err != nil {
		t.Logf("error not expected when selecting hashes: %s", err)
		t.Fail()
	}
	t.Logf("hashes fetched: %v\n", testHashes)

	teardownTest(t)
}

/* Tests that the InsertPost controller behaves as expected. Ensures that the SQL
driver correctly processes a successful entry */
func TestInsertPostSuccess(t *testing.T) {
	setupTest(t)

	// Mocking db operations with test post
	testPost := tp.Post{
		Title:       "test title",
		Author:      "tester",
		Contents:    "im a test",
		Descriptors: "test;test;test;test;test;test;test;test;test;test",
		Tag:         2,
	}
	mock.ExpectBegin()
	mock.ExpectExec("insert into Post").
		WithArgs(testPost.Title, testPost.Author, testPost.Contents,
			testPost.Descriptors, testPost.Tag).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Adding test post to mock database
	if err = testDbo.InsertPost(testPost); err != nil {
		t.Logf("error not expected when adding post: %s", err)
		t.Fail()
	}
	teardownTest(t)
}

/* Tests that the InsertPost controller behaves as expected. Ensures that the SQL
driver correctly processes a failed entry */
func TestInsertPostFailure(t *testing.T) {
	setupTest(t)

	// Mocking db operations with test post
	testPost := tp.Post{
		Title:       "test title",
		Author:      "tester",
		Contents:    "im a failed test",
		Descriptors: "test;test;test;test;test;test;test;test;test;test",
		Tag:         2,
	}
	mock.ExpectBegin()
	mock.ExpectExec("insert into Post").
		WithArgs(testPost.Title, testPost.Author, testPost.Contents,
			testPost.Descriptors, testPost.Tag).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	// Adding test post to mock database
	if err = testDbo.InsertPost(testPost); err == nil {
		t.Logf("was expecting error when adding post: %s", err)
	}
	teardownTest(t)
}

/* Called at the end of every test; ensuring all expectations met and database
is cleared */
func teardownTest(t *testing.T) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Log("there were unfulfilled expectations")
		t.Fail()
	}
	testDbo.db.Close()
}
