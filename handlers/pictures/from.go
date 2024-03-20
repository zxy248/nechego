package pictures

import (
	"io/fs"
	"math/rand/v2"
	"net/http"
	"path/filepath"
	"regexp"
	"sync"

	tu "github.com/zxy248/nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type FromDir struct {
	Path   string
	Regexp *regexp.Regexp

	files []string
	mu    sync.Mutex
}

func (h *FromDir) Match(c tele.Context) bool {
	return h.Regexp.MatchString(c.Text())
}

func (h *FromDir) Handle(c tele.Context) error {
	if err := h.initFiles(); err != nil {
		return err
	}

	file := h.files[rand.N(len(h.files))]
	return tu.SendFile(c, file)
}

func (h *FromDir) initFiles() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.files) > 0 {
		return nil
	}
	return filepath.WalkDir(h.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			h.files = append(h.files, path)
		}
		return nil
	})
}

type Locator interface {
	URL() string
}

type FromURL struct {
	Locator Locator
	Regexp  *regexp.Regexp
}

func (h *FromURL) Match(c tele.Context) bool {
	return h.Regexp.MatchString(c.Text())
}

func (h *FromURL) Handle(c tele.Context) error {
	resp, err := http.Get(h.Locator.URL())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(resp.Body)})
}

type FromStickerPack struct {
	Source []string
	Regexp *regexp.Regexp
}

func (h *FromStickerPack) Match(c tele.Context) bool {
	return h.Regexp.MatchString(c.Text())
}

func (h *FromStickerPack) Handle(c tele.Context) error {
	var all []tele.Sticker
	for _, src := range h.Source {
		set, err := c.Bot().StickerSet(src)
		if err != nil {
			return err
		}
		all = append(all, set.Stickers...)
	}
	return c.Send(&all[rand.N(len(all))])
}
