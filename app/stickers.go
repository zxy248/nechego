package app

import (
	"encoding/json"
	"math/rand"
	"os"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type Stickers struct {
	Hello, Masyunya, Poppy, Sima StickerPack
}

func (a *App) InitStickers() error {
	a.stickers = &Stickers{}
	var err error
	a.stickers.Hello, err = a.helloStickers()
	if err != nil {
		return err
	}
	a.stickers.Masyunya, err = a.masyunyaStickers()
	if err != nil {
		return err
	}
	a.stickers.Poppy, err = a.poppyStickers()
	if err != nil {
		return err
	}
	a.stickers.Sima, err = a.simaStickers()
	if err != nil {
		return err
	}
	return nil
}

type StickerPack []tele.Sticker

func (s StickerPack) Random() *tele.Sticker {
	return &s[rand.Intn(len(s))]
}

const helloStickersPath = "hello-stickers.json"

func (a *App) helloStickers() (StickerPack, error) {
	f, err := os.Open(a.Locate(helloStickersPath))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var s []tele.Sticker
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}
	return s, nil
}

const masyunyaStickersName = "masyunya_vk"

func (a *App) masyunyaStickers() (StickerPack, error) {
	set, err := a.bot.StickerSet(masyunyaStickersName)
	if err != nil {
		return nil, err
	}
	return set.Stickers, nil
}

var poppyStickersNames = []string{"pappy2_vk", "poppy_vk"}

func (a *App) poppyStickers() (StickerPack, error) {
	var stickers []tele.Sticker
	for _, s := range poppyStickersNames {
		set, err := a.bot.StickerSet(s)
		if err != nil {
			return nil, err
		}
		stickers = append(stickers, set.Stickers...)
	}
	return stickers, nil
}

var simaStickersName = "catsima_vk"

func (a *App) simaStickers() (StickerPack, error) {
	set, err := a.bot.StickerSet(simaStickersName)
	if err != nil {
		return nil, err
	}
	return set.Stickers, nil
}

// !привет
func (a *App) handleHello(c tele.Context) error {
	if strings.HasPrefix(getMessage(c).Raw, "!") || rand.Float64() <= a.pref.HelloChance {
		return c.Send(a.stickers.Hello.Random())
	}
	return nil
}

// !масюня
func (a *App) handleMasyunya(c tele.Context) error {
	return c.Send(a.stickers.Masyunya.Random())
}

// !паппи
func (a *App) handlePoppy(c tele.Context) error {
	return c.Send(a.stickers.Poppy.Random())
}

// !сима
func (a *App) handleSima(c tele.Context) error {
	return c.Send(a.stickers.Sima.Random())
}
