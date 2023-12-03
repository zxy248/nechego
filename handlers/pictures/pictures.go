package pictures

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	tele "gopkg.in/telebot.v3"
)

func randomFile(dir string) (string, error) {
	es, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	if len(es) == 0 {
		return "", fmt.Errorf("randomFile: empty directory %s", dir)
	}
	e := es[rand.Intn(len(es))]
	return filepath.Join(dir, e.Name()), nil
}

func randomSubdirFile(root string) (string, error) {
	dir, err := randomFile(root)
	if err != nil {
		return "", err
	}
	f, err := randomFile(dir)
	if err != nil {
		return "", err
	}
	return f, nil
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
	return &s.Stickers[rand.Intn(len(s.Stickers))], nil
}
