package main

import (
	"fmt"
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

	dsn := fmt.Sprintf("file:%s", os.Getenv("STORE"))
	if dsn == "file:" {
		log.Fatal("You must provide a database file name in the STORE environment variable.")
	}

	log.Printf("Initializing a database at %s...\n", dsn)
	store, err := newStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating a bot instance...")
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	app := &app{store}
	bot.Handle(tele.OnText, app.handleMessage)

	log.Println("The bot is running.")
	bot.Start()
}
