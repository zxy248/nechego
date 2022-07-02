package main

import (
	"database/sql"
	"log"
	"math/rand"
	"nechego/app"
	"nechego/model/sqlite"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	tele "gopkg.in/telebot.v3"
)

const (
	tokenEnv = "NECHEGO_TOKEN" // bot token
	dsnEnv   = "NECHEGO_DSN"   // data source name
	ownerEnv = "NECHEGO_OWNER" // bot owner's id
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	b, err := tele.NewBot(tele.Settings{
		Token:  token(),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}
	m, err := sqlite.NewModel(db())
	if err != nil {
		log.Fatal(err)
	}
	a := app.NewApp(b, m)
	log.Println("The bot has started.")
	a.Start()
}

func token() string {
	t := os.Getenv(tokenEnv)
	if t == "" {
		log.Fatalf("%v not set", tokenEnv)
	}
	return t
}

func db() *sql.DB {
	db, err := sql.Open("sqlite3", dsn())
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func dsn() string {
	d := os.Getenv(dsnEnv)
	if d == "" {
		log.Fatalf("%v not set", dsnEnv)
	}
	return d
}

func owner() int64 {
	o := os.Getenv(ownerEnv)
	if o == "" {
		log.Fatalf("%v not set", ownerEnv)
	}
	i, err := strconv.ParseInt(o, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return i
}
