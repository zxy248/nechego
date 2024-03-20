package pictures

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"regexp"

	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Car struct{}

func (h *Car) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!авто", "!машин", "!тачка")
}

func (h *Car) Handle(c tele.Context) error {
	const url = "https://www.thisautomobiledoesnotexist.com/"

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	b64 := carImageRe.FindSubmatch(body)[1]
	data := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(b64))
	return c.Send(&tele.Photo{File: tele.FromReader(data)})
}

var carImageRe = regexp.MustCompile(`<img id = "vehicle" src="data:image/png;base64,(.+)" class="center">`)
