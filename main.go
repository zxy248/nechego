package main

import (
	"fmt"
	"math/rand"
	"nechego/app"
	"nechego/model"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
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
		Token:  botToken(),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	fail(err)

	mod := model.NewModel(database())

	var logger *zap.Logger
	if debugFlag() {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	fail(err)
	defer logger.Sync()

	a := app.NewApp(bot, mod, logger)
	a.Start()
}

func botToken() string {
	v := os.Getenv(tokenEnv)
	if v == "" {
		fail(fmt.Errorf("%s not set", tokenEnv))
	}
	return v
}

func database() *sqlx.DB {
	db := sqlx.MustConnect("sqlite3", databaseDSN())
	return db
}

func databaseDSN() string {
	v := os.Getenv(dsnEnv)
	if v == "" {
		fail(fmt.Errorf("%s not set", dsnEnv))
	}
	return v
}

func ownerUID() int64 {
	owner := os.Getenv(ownerEnv)
	if owner == "" {
		fail(fmt.Errorf("%s not set", ownerEnv))
	}
	id, err := strconv.ParseInt(owner, 10, 64)
	if err != nil {
		fail(err)
	}
	return id
}

func debugFlag() bool {
	v := os.Getenv(debugEnv)
	return v == "1"
}

func fail(err error) {
	if err != nil {
		panic(err)
	}
}
