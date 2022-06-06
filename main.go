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
	store  *store
	status *status
}

func init() {
	log.Println("Initializing the random number generator...")
	rand.Seed(time.Now().UnixNano())
}

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	dsn := fmt.Sprintf("file:%s", os.Getenv("STORE"))
	if dsn == "file:" {
		log.Fatal("You must provide a database file name in the STORE environment variable.")
	}

	log.Printf("Connecting to the database at %s...\n", dsn)
	store, err := newStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Building a bot...")
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	app := &app{
		store:  store,
		status: newStatus(),
	}

	bot.Handle(tele.OnText, app.processInput)

	log.Println("The bot is running.")
	bot.Start()
}
