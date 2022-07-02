package sqlite

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func testingDB() *DB {
	name := randString(32)
	sqlDB, err := sql.Open("sqlite3", fmt.Sprintf("testdata/db/%s.db", name))
	if err != nil {
		panic(err)
	}
	wrappedDB := DB{sqlDB}
	if err := wrappedDB.CreateTables(); err != nil {
		panic(err)
	}
	return &wrappedDB
}

func randString(n int) string {
	alphabet := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	r := make([]rune, 0, n)
	for i := 0; i < n; i++ {
		r = append(r, alphabet[rand.Intn(len(alphabet))])
	}
	return string(r)
}
