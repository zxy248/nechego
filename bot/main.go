package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zxy248/nechego/data"

	tele "gopkg.in/zxy248/telebot.v3"
)

var (
	botToken        = getenv("NECHEGO_TOKEN")
	assetsDirectory = getenv("NECHEGO_ASSETS")
)

func main() {
	bot, err := tele.NewBot(tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal("cannot build bot: ", err)
	}

	pool, err := pgxpool.New(context.Background(), getenv("NECHEGO_DATABASE"))
	if err != nil {
		log.Fatal("cannot create connection pool: ", err)
	}
	defer pool.Close()

	app := &App{Queries: data.New(pool)}
	srv := &Server{
		Bot:      bot,
		Handlers: app.Handlers(),
	}
	go srv.Start()
	<-interrupt()
	srv.Stop()
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
