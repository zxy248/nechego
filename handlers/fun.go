package handlers

import (
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"

	tele "gopkg.in/telebot.v3"
)

type Mouse struct {
	Path string // path to video file
}

var mouseRe = regexp.MustCompile("^!мыш")

func (h *Mouse) Regexp() *regexp.Regexp {
	return mouseRe
}

func (h *Mouse) Handle(c tele.Context) error {
	return c.Send(&tele.Video{File: tele.FromDisk(h.Path)})
}

type Tiktok struct {
	Path string // path to directory with webms
}

var tiktokRe = regexp.MustCompile("^!тикток")

func (h *Tiktok) Regexp() *regexp.Regexp {
	return tiktokRe
}

func (h *Tiktok) Handle(c tele.Context) error {
	files, err := os.ReadDir(h.Path)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("empty directory")
	}
	f := files[rand.Intn(len(files))]
	return c.Send(&tele.Video{File: tele.FromDisk(filepath.Join(h.Path, f.Name()))})
}

type Game struct{}

var gameRe = regexp.MustCompile("^!игр")

func (h *Game) Regexp() *regexp.Regexp {
	return gameRe
}

func (h *Game) Handle(c tele.Context) error {
	games := [...]*tele.Dice{tele.Dart, tele.Ball, tele.Goal, tele.Slot, tele.Bowl}
	return c.Send(games[rand.Intn(len(games))])
}
