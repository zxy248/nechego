package pictures

import (
	"bytes"
	"encoding/base64"
	"nechego/handlers"
	"regexp"

	tele "gopkg.in/telebot.v3"
)

type Car struct{}

func (h *Car) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!авто", "!машин", "!тачка")
}

func (h *Car) Handle(c tele.Context) error {
	const url = "https://www.thisautomobiledoesnotexist.com/"
	body, err := getBytes(url)
	if err != nil {
		return err
	}
	b64 := carImageRe.FindSubmatch(body)[1]
	data := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(b64))
	return c.Send(&tele.Photo{File: tele.FromReader(data)})
}

var carImageRe = regexp.MustCompile(`<img id = "vehicle" src="data:image/png;base64,(.+)" class="center">`)
