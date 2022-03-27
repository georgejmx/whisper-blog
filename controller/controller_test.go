package controller

import (
	"fmt"
	"testing"
	tp "whisper-blog/types"

	"github.com/DATA-DOG/go-sqlmock"
)

var (
	testDbo DatabaseObject
	mock    sqlmock.Sqlmock
)

func TestAddPostSuccess(t *testing.T) {
	// Opening stub database
	var err error
	testDbo.db, mock, err = sqlmock.New()
	if err != nil {
		t.Log("error when opening stub database")
		t.Fail()
	}

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
	if err = testDbo.AddPost(testPost); err != nil {
		t.Log(fmt.Sprintf("error not expected when adding post: %s", err))
		t.Fail()
	}

	// Ensuring all expectations met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Log("there were unfulfilled expectations")
		t.Fail()
	}

	testDbo.db.Close()
}

func TestAddPostFailure(t *testing.T) {
	// Opening stub database
	var err error
	testDbo.db, mock, err = sqlmock.New()
	if err != nil {
		t.Log("error when opening stub database")
		t.Fail()
	}

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
	if err = testDbo.AddPost(testPost); err == nil {
		t.Log(fmt.Sprintf("was expecting error when adding post: %s", err))
	}

	// Ensuring all expectations met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Log("there were unfulfilled expectations")
		t.Fail()
	}
}
