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
	dirs, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	if len(dirs) == 0 {
		return "", fmt.Errorf("randomFile: empty directory %s", dir)
	}
	return filepath.Join(dir, dirs[rand.Intn(len(dirs))].Name()), nil
}

func randomFileFromSubdir(dir string) (string, error) {
	f1, err := randomFile(dir)
	if err != nil {
		return "", err
	}
	f2, err := randomFile(f1)
	if err != nil {
		return "", err
	}
	return f2, nil
}

func download(addr string) ([]byte, error) {
	r, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func randomSticker(c tele.Context, setName string) (*tele.Sticker, error) {
	set, err := c.Bot().StickerSet(setName)
	if err != nil {
		return nil, err
	}
	return &set.Stickers[rand.Intn(len(set.Stickers))], nil
}
