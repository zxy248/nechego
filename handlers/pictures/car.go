package pictures

import (
	"bytes"
	"encoding/base64"
	"io"
	"nechego/handlers"
	"regexp"

	tele "gopkg.in/telebot.v3"
)

type Car struct{}

func (h *Car) Match(s string) bool {
	return handlers.MatchPrefixes([]string{"!авто", "!машин", "!тачка"}, s)
}

func (h *Car) Handle(c tele.Context) error {
	car, err := randomCar()
	if err != nil {
		return err
	}
	return c.Send(&tele.Photo{File: tele.FromReader(car)})
}

var carImageRegexp = regexp.MustCompile(`<img id = "vehicle" src="data:image/png;base64,(.+)" class="center">`)

func randomCar() (io.Reader, error) {
	page, err := download("https://www.thisautomobiledoesnotexist.com/")
	if err != nil {
		return nil, err
	}
	data := carImageRegexp.FindSubmatch(page)[1]
	return base64.NewDecoder(base64.StdEncoding, bytes.NewReader(data)), nil
}
