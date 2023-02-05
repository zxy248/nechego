package main

import (
	"log"
	"math/rand"
	"nechego/fishing"
	"nechego/game"
	"nechego/handlers"
	"os"
	"os/signal"
	"syscall"
	"time"

	tele "gopkg.in/telebot.v3"
)

// refillMarket refills every market in the universe with a new
// product at a specified time interval.
func refillMarket(universe *game.Universe) {
	for range time.NewTicker(time.Minute).C {
		universe.ForEachWorld(func(w *game.World) {
			w.Market.Refill()
		})
	}
}

// restoreEnergy restores the energy levels of all users in the
// universe at a specified time interval.
func restoreEnergy(universe *game.Universe) {
	for range time.NewTicker(time.Minute).C {
		universe.ForEachWorld(func(w *game.World) {
			for _, u := range w.Users {
				e := game.Energy(0.01)
				if u.InventoryFull() {
					e /= 2
				}
				u.Energy.Add(e)
			}
		})
	}
}

// fillNet fills the users' cast fishing nets.
func fillNet(universe *game.Universe) {
	for range time.NewTicker(time.Minute).C {
		universe.ForEachWorld(func(w *game.World) {
			for _, u := range w.Users {
				if rand.Float64() < 0.04 {
					u.FillNet()
				}
			}
		})
	}
}

// stopper gracefully stops the bot after receiving an interrupt
// signal and sends an empty structure on the done channel.
func stopper(bot *tele.Bot, universe *game.Universe) (done chan struct{}) {
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
