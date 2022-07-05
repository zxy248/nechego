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
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

const (
	tokenEnv = "NECHEGO_TOKEN" // bot token
	dsnEnv   = "NECHEGO_DSN"   // data source name
	ownerEnv = "NECHEGO_OWNER" // bot owner's id
	debugEnv = "NECHEGO_DEBUG" // set to 1 for debug output
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	bot, err := tele.NewBot(tele.Settings{
		Token:  token(),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	model, err := sqlite.NewModel(db())
	if err != nil {
		log.Fatal(err)
	}

	var logger *zap.Logger
	if debug() {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	a := app.NewApp(bot, model, logger)
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

func debug() bool {
	v := os.Getenv(debugEnv)
	return v == "1"
}
