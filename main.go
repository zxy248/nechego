package main

import (
	"database/sql"
	"log"
	"math/rand"
	"nechego/bot"
	"nechego/model"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
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
	cfg := &bot.Config{
		Token: token(),
		DB:    db(),
		Owner: owner(),
	}
	if err := cfg.DB.Setup(); err != nil {
		log.Fatal(err)
	}
	bot, err := bot.NewBot(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("The bot has started.")
	bot.Start()
}

func token() string {
	t := os.Getenv(tokenEnv)
	if t == "" {
		log.Fatalf("%v not set", tokenEnv)
	}
	return t
}

func db() *model.DB {
	db, err := sql.Open("sqlite3", dsn())
	if err != nil {
		log.Fatal(err)
	}
	return &model.DB{DB: db}
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
