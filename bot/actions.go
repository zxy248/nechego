package main

import (
	"log"
	"nechego/game"
	"os"
	"os/signal"
	"syscall"
	"time"

	tele "gopkg.in/telebot.v3"
)

func refillMarket(universe *game.Universe) {
	for range time.NewTicker(time.Minute).C {
		universe.ForEachWorld(func(w *game.World) {
			w.Market.Refill()
		})
	}
}

func restoreEnergy(universe *game.Universe) {
	const step = 4
	for range time.NewTicker(step * time.Minute).C {
		universe.ForEachWorld(func(w *game.World) {
			for _, u := range w.Users {
				if u.Inventory.Count() > game.InventorySize {
					u.RestoreEnergy(step / 2)
				} else {
					u.RestoreEnergy(step)
				}
			}
		})
	}
}

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
