package models

import (
	"database/sql"
	"github.com/w3gop2p/letsgoweb/internal/assert"
	"os"

	"testing"
)

func TestUserModelExists(t *testing.T) {
	// Skip the test if the "-short" flag is provided when running the test.
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	// Set up a suite of table-driven tests and expected results.
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the newTestDB() helper function to get a connection pool to
			// our test database. Calling this here -- inside t.Run() -- means
			// that fresh database tables and data will be set up and torn down
			// for each sub-test.
			db := newTestDB(t)
			// Create a new instance of the UserModel.
			m := UserModel{db}
			// Call the UserModel.Exists() method and check that the return
			// value and error match the expected values for the sub-test.
			exists, err := m.Exists(tt.userID)
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}

func newTestDB(t *testing.T) *sql.DB {
	// Establish a sql.DB connection pool for our test database. Because our
	// setup and teardown scripts contains multiple SQL statements, we need
	// to use the "multiStatements=true" parameter in our DSN. This instructs
	// our MySQL database driver to support executing multiple SQL statements
	// in one db.Exec() call.
	db, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}
	// Read the setup SQL script from file and execute the statements.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}
	// Use the t.Cleanup() to register a function *which will automatically be
	// called by Go when the current test (or sub-test) which calls newTestDB()
	// has finished*. In this function we read and execute the teardown script,
	// and close the database connection pool.
	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	})
	// Return the database connection pool.
	return db
}
