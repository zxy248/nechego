package pictures

import (
	"bytes"
	"math/rand"
	"nechego/danbooru"
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Danbooru struct {
	API *danbooru.Danbooru
}

func (h *Danbooru) Match(s string) bool {
	return handlers.HasPrefix(s, "!данбору")
}

func (h *Danbooru) Handle(c tele.Context) error {
	pic, err := h.API.Get(danbooru.All)
	if err != nil {
		return err
	}
	photo := &tele.Photo{File: tele.FromReader(bytes.NewReader(pic.Data))}
	if pic.Rating == danbooru.Explicit {
		photo.Caption = randomWarningNSFW()
		photo.HasSpoiler = true
	}
	return c.Send(photo, tele.ModeHTML)
}

func randomWarningNSFW() string {
	caps := [...]string{
		"🔞 Осторожно! Только для взрослых.",
		"<i>Содержимое предназначено для просмотра лицами старше 18 лет.</i>",
		"<b>ВНИМАНИЕ!</b> Вы увидите фотографии взрослых голых женщин. Будьте сдержанны.",
	}
	return caps[rand.Intn(len(caps))]
}
