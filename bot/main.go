package main

import (
	"fmt"
	"log"
	"math/rand"
	"nechego/avatar"
	"nechego/game"
	"nechego/handlers"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	universeDirectory = "universe"
	avatarDirectory   = "avatar"
)

var (
	botToken        = getEnv("NECHEGO_TOKEN")
	assetsDirectory = getEnv("NECHEGO_ASSETS")
)

func main() {
	rand.Seed(time.Now().UnixNano())

	pref := tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	as := &avatar.Storage{
		Dir:       avatarDirectory,
		MaxWidth:  1500,
		MaxHeight: 1500,
		Bot:       bot,
	}
	u := game.NewUniverse(universeDirectory, worldInitializer(bot))

	// Register handlers.
	handlers := []ContextService{}
	for _, h := range RemainingHandlers(u, as) {
		handlers = append(handlers, &TextHandler{h})
	}
	handlers = append(handlers, PictureHandlers()...)
	handlers = append(handlers, CasinoHandlers(u)...)
	handlers = append(handlers, CallbackHandlers(u)...)

	// Apply middleware.
	mw := globalMiddleware(u, as)
	for i, h := range handlers {
		handlers[i] = Wrap(h, mw...)
	}

	// Start the bot.
	rt := &Router{handlers}
	endpoints := []string{
		tele.OnText,
		tele.OnPhoto,
		tele.OnCallback,
		tele.OnDice,
	}
	for _, e := range endpoints {
		bot.Handle(e, rt.Dispatch)
	}
	serve(bot, u)
}

func serve(bot *tele.Bot, u *game.Universe) {
	done := notifyStop(bot, u)
	bot.Start()
	<-done
	log.Println("Successful shutdown.")
}

// worldInitializer returns a function that registers the fishing
// record announcer.
func worldInitializer(bot *tele.Bot) func(*game.World) {
	return func(w *game.World) {
		w.History.Announce(handlers.RecordAnnouncer(bot, tele.ChatID(w.TGID)))
	}
}

// notifyStop gracefully stops the bot after receiving an interrupt
// signal and sends an empty structure on the done channel.
func notifyStop(bot *tele.Bot, universe *game.Universe) (done chan struct{}) {
	done = make(chan struct{})
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-interrupt

		log.Println("Stopping the bot...")
		bot.Stop()

		log.Println("Saving the universe...")
		if err := universe.SaveAll(); err != nil {
			log.Fatal(err)
		}
		done <- struct{}{}
	}()
	return done
}

func assetPath(s string) string {
	return filepath.Join(assetsDirectory, s)
}

func getEnv(s string) string {
	e := os.Getenv(s)
	if e == "" {
		panic(fmt.Sprintf("%s not set", s))
	}
	return e
}
