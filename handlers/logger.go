package handlers

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

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

	flag := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile(l.path(id), flag, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.WriteString(f, s+"\n")
	return err
}

func (l *Logger) Messages(id int64) ([]string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.Open(l.path(id))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var messages []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		messages = append(messages, scanner.Text())
	}
	return messages, nil
}

func (l *Logger) path(id int64) string {
	return filepath.Join(l.Dir, strconv.FormatInt(id, 10))
}
