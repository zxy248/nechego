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

func (h *Danbooru) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!данбору")
}

func (h *Danbooru) Handle(c tele.Context) error {
	pic, err := h.API.Get(danbooru.All)
	if err != nil {
		return err
	}
	r := bytes.NewReader(pic.Data)
	p := &tele.Photo{File: tele.FromReader(r)}
	if pic.Rating == danbooru.Explicit {
		p.Caption = warningNSFW()
		p.HasSpoiler = true
	}
	return c.Send(p, tele.ModeHTML)
}

func warningNSFW() string {
	s := [...]string{
		"🔞 Осторожно! Только для взрослых.",
		"<i>Содержимое предназначено для просмотра лицами старше 18 лет.</i>",
		"<b>ВНИМАНИЕ!</b> Вы увидите фотографии взрослых голых женщин. Будьте сдержанны.",
	}
	return s[rand.Intn(len(s))]
}
