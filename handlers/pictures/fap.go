package pictures

import (
	"bytes"
	"nechego/danbooru"
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Fap struct {
	API *danbooru.Danbooru
}

var fapRe = handlers.Regexp("^!(др[ао]ч|фап|эро|порн)")

func (h *Fap) Match(c tele.Context) bool {
	return fapRe.MatchString(c.Text())
}

func (h *Fap) Handle(c tele.Context) error {
	pic, err := h.API.Get(danbooru.NSFW)
	if err != nil {
		return err
	}
	photo := &tele.Photo{
		File:       tele.FromReader(bytes.NewReader(pic.Data)),
		Caption:    ratingEmoji(pic.Rating),
		HasSpoiler: true,
	}
	return c.Send(photo, tele.ModeHTML)
}

func ratingEmoji(r danbooru.Rating) string {
	switch r {
	case danbooru.Explicit:
		return "🔞"
	case danbooru.Questionable:
		return "❓"
	}
	return ""
}
