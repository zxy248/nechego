package main

import (
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	tele "gopkg.in/telebot.v3"
)

type pairOfTheDay struct {
	set  bool
	prev time.Time
	x    int64
	y    int64
	mu   *sync.Mutex
}

type app struct {
	store *store
	pair  *pairOfTheDay
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

	pair := &pairOfTheDay{
		set: false,
		mu:  new(sync.Mutex),
	}
	app := &app{store, pair}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle(tele.OnText, app.handleMessage)
	bot.Start()
}
