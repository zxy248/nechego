package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zxy248/nechego/game"

	tele "gopkg.in/zxy248/telebot.v3"
)

var (
	botToken         = getenv("NECHEGO_TOKEN")
	assetsDirectory  = getenv("NECHEGO_ASSETS")
	storageDirectory = getenv("NECHEGO_STORAGE")
)

func main() {
	bot, err := tele.NewBot(tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal("cannot build bot: ", err)
	}
	app := &App{Universe: game.NewUniverse(storageDirectory)}
	srv := &Server{
		Bot:      bot,
		Handlers: app.Handlers(),
	}
	go srv.Start()
	<-interrupt()
	srv.Stop()
	if err := app.Shutdown(); err != nil {
		log.Fatal("cannot shutdown: ", err)
	}
	log.Println("successful shutdown")
}

func getenv(s string) string {
	e := os.Getenv(s)
	if e == "" {
		panic(fmt.Sprintf("%s not set", s))
	}
	return e
}

func interrupt() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}
