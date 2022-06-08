package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	tele "gopkg.in/telebot.v3"
)

const helloStickersPath = "hello-stickers.json"

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

type stickersCollector struct {
	stickers []*tele.Sticker
	mu       *sync.Mutex
}

func newStickersCollector() *stickersCollector {
	return &stickersCollector{
		[]*tele.Sticker{},
		&sync.Mutex{},
	}
}

func (s *stickersCollector) collectStickers(c tele.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	v := c.Message().Sticker
	s.stickers = append(s.stickers, v)

	fmt.Println("Sticker has been collected.")
	return nil
}

func (s *stickersCollector) writeStickers(c tele.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Create(fmt.Sprintf("stickers-%v.json", time.Now().Unix()))
	if err != nil {
		return err
	}
	defer f.Close()
	if err = json.NewEncoder(f).Encode(s.stickers); err != nil {
		return err
	}
	fmt.Println("Stickers have been written.")
	return nil
}
