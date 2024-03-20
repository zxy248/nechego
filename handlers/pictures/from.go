package pictures

import (
	"io/fs"
	"math/rand/v2"
	"path/filepath"
	"regexp"
	"sync"

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
	return c.Send(sendableFromFile(file))
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
