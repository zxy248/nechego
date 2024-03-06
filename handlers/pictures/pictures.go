package pictures

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	tele "gopkg.in/zxy248/telebot.v3"
)

func randomFile(dir string) (string, error) {
	es, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	if len(es) == 0 {
		return "", fmt.Errorf("randomFile: empty directory %s", dir)
	}
	e := es[rand.N(len(es))]
	return filepath.Join(dir, e.Name()), nil
}

func randomSubdirFile(root string) (string, error) {
	dir, err := randomFile(root)
	if err != nil {
		return "", err
	}
	name, err := randomFile(dir)
	if err != nil {
		return "", err
	}
	return name, nil
}

func getBytes(addr string) ([]byte, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func randomSticker(c tele.Context, set string) (*tele.Sticker, error) {
	s, err := c.Bot().StickerSet(set)
	if err != nil {
		return nil, err
	}
	return &s.Stickers[rand.N(len(s.Stickers))], nil
}

func sendableFromFile(name string) tele.Sendable {
	f := tele.FromDisk(name)
	ext := strings.ToLower(path.Ext(name))
	if ext == ".mp4" || ext == ".mov" {
		return &tele.Video{File: f}
	}
	return &tele.Photo{File: f}
}
