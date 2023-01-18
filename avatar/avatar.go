package avatar

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

var ErrSize = errors.New("avatar too large")

type Storage struct {
	Bot       *tele.Bot
	Dir       string
	MaxWidth  int
	MaxHeight int
}

func (s *Storage) Set(id int64, avatar *tele.Photo) error {
	if avatar.Width > s.MaxWidth || avatar.Height >= s.MaxHeight {
		return ErrSize
	}
	if err := os.MkdirAll(s.Dir, 0777); err != nil {
		return err
	}

	src, err := s.Bot.File(&avatar.File)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(s.path(id))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (s *Storage) Get(id int64) (avatar *tele.Photo, ok bool) {
	_, err := os.Stat(s.Dir)
	if err != nil {
		return nil, false
	}
	f := tele.FromDisk(s.path(id))
	if !f.OnDisk() {
		return nil, false
	}
	return &tele.Photo{File: f}, true
}

func (s *Storage) path(id int64) string {
	return filepath.Join(s.Dir, strconv.FormatInt(id, 10))
}
