package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	tele "gopkg.in/telebot.v3"
)

const (
	tokenEnv           = "TOKEN"            // bot token
	dsnEnv             = "DSN"              // data source name
	whitelistEnv       = "WHITELIST"        // comma-separated list of group IDs
	ownersEnv          = "OWNERS"           // comma-separated list of user IDs
	collectStickersEnv = "COLLECT_STICKERS" // set it to 1 to collect stickers
)

type app struct {
	store     *store
	status    *status
	whitelist *whitelist
	owners    *owners
	keyboard  *tele.ReplyMarkup
	bans      *sync.Map
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

	owners, err := getOwners()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Owners: %s\n", owners)

	log.Println("Building a bot...")
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	app := &app{
		store:     store,
		status:    newStatus(),
		whitelist: whitelist,
		owners:    owners,
		keyboard:  newKeyboard(),
		bans:      &sync.Map{},
	}

	bot.Handle(tele.OnText, app.processInput)
	if getCollectStickers() {
		c := newStickersCollector()
		bot.Handle("/write-stickers", c.writeStickers)
		bot.Handle(tele.OnSticker, c.collectSticker)
	}
	log.Println("The bot is running.")
	bot.Start()
}

func getPref() (tele.Settings, error) {
	token := os.Getenv(tokenEnv)
	if token == "" {
		return tele.Settings{}, fmt.Errorf("%s must be set", tokenEnv)
	}
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	return pref, nil
}

func getDSN() (string, error) {
	dsn := os.Getenv(dsnEnv)
	if dsn == "" {
		return "", fmt.Errorf("%s must be set", dsnEnv)
	}
	return dsn, nil
}

func getWhitelist() (*whitelist, error) {
	v := os.Getenv(whitelistEnv)
	if v == "" {
		return nil, fmt.Errorf("%s must be set", whitelistEnv)
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
	v := os.Getenv(collectStickersEnv)
	if v == "1" {
		return true
	}
	return false
}

func getOwners() (*owners, error) {
	v := os.Getenv(ownersEnv)
	if v == "" {
		return nil, fmt.Errorf("%s must be set", ownersEnv)
	}
	var ids []int64
	for _, s := range strings.Split(v, ",") {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	o := newOwners(ids...)
	return o, nil
}
