package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

type app struct {
	store *store
}

func main() {
	rand.Seed(time.Now().Unix())

	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	dsn := "file:store.db"
	store, err := newStore(dsn)
	if err != nil {
		log.Fatal(err)
	}
	app := &app{store}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle(tele.OnText, app.handleMessage)
	bot.Start()
}
