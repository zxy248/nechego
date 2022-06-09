package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	tele "gopkg.in/telebot.v3"
)

const helloStickersPath = "data/hello-stickers.json"

// helloStickers is a list of stickers saying "Hi".
var helloStickers = func() []*tele.Sticker {
	f, err := os.Open(helloStickersPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var ss []*tele.Sticker
	if err := json.NewDecoder(f).Decode(&ss); err != nil {
		panic(err)
	}
	return ss
}()

// stickerCollector is a list of stickers to write to the file later.
type stickerCollector struct {
	stickers []*tele.Sticker
	mu       *sync.Mutex
}

// newStickersCollector returns a new sticker collector.
func newStickersCollector() *stickerCollector {
	return &stickerCollector{
		[]*tele.Sticker{},
		&sync.Mutex{},
	}
}

// collectSticker adds the sticker to the sticker list.
func (s *stickerCollector) collectSticker(c tele.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	v := c.Message().Sticker
	s.stickers = append(s.stickers, v)

	log.Println("Sticker has been collected.")
	return nil
}

const stickersFileFormat = "stickers-%v.json"

// writeStickers writes the sticker list to the file.
func (s *stickerCollector) writeStickers(c tele.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Create(fmt.Sprintf(stickersFileFormat, time.Now().Unix()))
	if err != nil {
		return err
	}
	defer f.Close()
	if err = json.NewEncoder(f).Encode(s.stickers); err != nil {
		return err
	}
	log.Println("Stickers have been written.")
	return nil
}
