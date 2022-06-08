package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	tokenEnv     = "TOKEN"            // bot token
	storeEnv     = "STORE"            // database file path
	whitelistEnv = "WHITELIST"        // comma-separated list of group IDs
	stickersEnv  = "COLLECT_STICKERS" // set it to 1 to collect stickers
)

type app struct {
	store     *store
	status    *status
	whitelist *whitelist
	keyboard  *tele.ReplyMarkup
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	pref, err := getPref()
	if err != nil {
		log.Fatal(err)
	}

	dsn, err := getDSN()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connecting to the database at %s...\n", dsn)
	store, err := newStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	whitelist, err := getWhitelist()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Whitelist: %s\n", whitelist)

	log.Println("Building a bot...")
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	app := &app{
		store:     store,
		status:    newStatus(),
		whitelist: whitelist,
		keyboard:  newKeyboard(),
	}

	bot.Handle(tele.OnText, app.processInput)
	if getCollectStickers() {
		sc := newStickersCollector()
		bot.Handle("/write-stickers", sc.writeStickers)
		bot.Handle(tele.OnSticker, sc.collectStickers)
	}
	log.Println("The bot is running.")
	bot.Start()
}

func getPref() (tele.Settings, error) {
	token := os.Getenv("TOKEN")
	if token == "" {
		return tele.Settings{},
			fmt.Errorf("You must provide a bot token in the %s"+
				" environment variable.", tokenEnv)
	}
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	return pref, nil
}

func getDSN() (string, error) {
	v := os.Getenv("STORE")
	if v == "" {
		return "",
			fmt.Errorf("You must provide a database file name in the %s"+
				" environment variable.", storeEnv)
	}
	dsn := fmt.Sprintf("file:%s", v)
	return dsn, nil
}

func getWhitelist() (*whitelist, error) {
	v := os.Getenv(whitelistEnv)
	if v == "" {
		return nil,
			fmt.Errorf("You must provide a list of IDs in the %s"+
				" environment variable.", whitelistEnv)
	}
	w := newWhitelist()
	allowedGroupIDs := strings.Split(v, ",")
	for _, id := range allowedGroupIDs {
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		w.add(i)
	}
	return w, nil
}

func getCollectStickers() bool {
	v := os.Getenv(stickersEnv)
	if v == "1" {
		return true
	}
	return false
}
