package bot

import (
	"encoding/json"
	"os"
	"path/filepath"

	tele "gopkg.in/telebot.v3"
)

var helloStickersPath = filepath.Join(dataPath, "hello-stickers.json")

// helloStickers is a list of stickers saying "Hi".
var helloStickers = func() []*tele.Sticker {
	f, err := os.Open(helloStickersPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var ss []*tele.Sticker
	if err := json.NewDecoder(f).Decode(&ss); err != nil {
		panic(err)
	}
	return ss
}()
