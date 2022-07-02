package app

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	tele "gopkg.in/telebot.v3"
)

var helloStickersPath = filepath.Join(dataPath, "hello-stickers.json")

// helloStickers is a list of stickers saying "Hi".
var helloStickers = func() []*tele.Sticker {
	f, err := os.Open(helloStickersPath)
	if err != nil {
		log.Printf("helloStickers: %v", err)
		return nil
	}
	defer f.Close()

	var ss []*tele.Sticker
	if err := json.NewDecoder(f).Decode(&ss); err != nil {
		log.Printf("helloStickers: %v", err)
		return nil
	}
	return ss
}()

func helloSticker() *tele.Sticker {
	return helloStickers[rand.Intn(len(helloStickers))]
}
