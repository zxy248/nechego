package main

import (
	"log"
	"nechego/fishing"
	"nechego/game"
	"nechego/handlers"
	"os"
	"os/signal"
	"syscall"

	tele "gopkg.in/telebot.v3"
)

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

// worldInitializer returns a function that creates fishing record
// channels and starts the record announcer.
func worldInitializer(bot *tele.Bot) func(w *game.World) {
	return func(w *game.World) {
		weightRecords := w.History.Records(fishing.Weight)
		lengthRecords := w.History.Records(fishing.Length)
		priceRecords := w.History.Records(fishing.Price)
		tgid := tele.ChatID(w.TGID)
		handlers.RecordAnnouncer(bot, tgid, weightRecords, lengthRecords, priceRecords)
	}
}
