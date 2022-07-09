package app

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

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

const masyunyaStickersName = "masyunya_vk"

func (a *App) masyunyaHandler() tele.HandlerFunc {
	set, err := a.bot.StickerSet(masyunyaStickersName)
	if err != nil {
		log.Println("masyunyaHandler unavailable: ", err)
		return func(c tele.Context) error {
			return nil
		}
	}
	return func(c tele.Context) error {
		return c.Send(&set.Stickers[rand.Intn(len(set.Stickers))])
	}
}

var poppyStickersNames = []string{"pappy2_vk", "poppy_vk"}

func (a *App) poppyHandler() tele.HandlerFunc {
	var stickers []tele.Sticker
	for _, sn := range poppyStickersNames {
		set, err := a.bot.StickerSet(sn)
		if err != nil {
			log.Println("poppyHandler unavailable: ", err)
			return func(c tele.Context) error {
				return nil
			}
		}
		stickers = append(stickers, set.Stickers...)
	}
	return func(c tele.Context) error {
		return c.Send(&stickers[rand.Intn(len(stickers))])
	}
}

const helloChance = 0.2

// handleHello sends a hello sticker
func (a *App) handleHello(c tele.Context) error {
	if strings.HasPrefix(getMessage(c).Raw, "!") || rand.Float64() <= helloChance {
		return c.Send(helloSticker())
	}
	return nil
}
