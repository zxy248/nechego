package handlers

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Pass struct {
	Logger *Logger
}

func (h *Pass) Match(c tele.Context) bool {
	return true
}

func (h *Pass) Handle(c tele.Context) error {
	text := strings.TrimSpace(c.Text())
	if !strings.ContainsRune(text, '\n') && text != "" && len(text) < 1024 {
		if err := h.Logger.Log(c.Chat().ID, text); err != nil {
			return err
		}
	}
	return nil
}

type Logger struct {
	Dir string
	mu  sync.Mutex
}

func (l *Logger) Log(id int64, s string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if err := os.MkdirAll(l.Dir, 0777); err != nil {
		return err
	}

	name := filepath.Join(l.Dir, strconv.FormatInt(id, 10))
	flag := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile(name, flag, 0666)
	if err != nil {
		return err
	}

	_, err = io.WriteString(f, s+"\n")
	return err
}
