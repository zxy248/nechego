package pictures

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"nechego/handlers"
	"net/http"

	tele "gopkg.in/telebot.v3"
)

type Cat struct{}

func (h *Cat) Match(s string) bool {
	return handlers.HasPrefix(s, "!кот", "!кош")
}

func (h *Cat) Handle(c tele.Context) error {
	pic, err := randomCatPicture()
	if err != nil {
		return err
	}
	return c.Send(&tele.Photo{File: tele.FromReader(pic)})
}

type catPicture struct {
	bytes.Buffer
}

func catPictureURL() string {
	const base = "https://d2ph5fj80uercy.cloudfront.net"
	page := 1 + rand.Intn(6)
	n := rand.Intn(5000)
	return fmt.Sprintf("%s/%02d/cat%d.jpg", base, page, n)
}

func randomCatPicture() (*catPicture, error) {
	r, err := http.Get(catPictureURL())
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	pic := &catPicture{}
	io.Copy(pic, r.Body)
	return pic, nil
}
