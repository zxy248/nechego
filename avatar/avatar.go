package avatar

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

type Storage struct {
	Dir string
}

func (s *Storage) Set(id int64, p io.Reader) error {
	if err := os.MkdirAll(s.Dir, 0777); err != nil {
		return err
	}
	f, err := os.Create(s.path(id))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, p)
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
