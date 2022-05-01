# Archive

This is where code that has been removed is placed. This is code that is no
longer used in production or integration tests, however may be useful later

## Selecting Latest Hash

No longer needed because of `SelectCandidateHashes` superceeds it

`controller.go`
---
```
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
```

`controller_test.go`
---
```
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
```