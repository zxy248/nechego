package pictures

import (
	"github.com/zxy248/nechego/handlers"
	"github.com/zxy248/nechego/handlers/pictures/danbooru"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Fap struct{}

var fapRe = handlers.NewRegexp("^!(др[ао]ч|фап|эро|порн)")

func (h *Fap) Match(c tele.Context) bool {
	return fapRe.MatchString(c.Text())
}

func (h *Fap) Handle(c tele.Context) error {
	var pic *danbooru.Picture
	for {
		pic = <-danbooruPictures
		if pic.Rating.NSFW() {
			break
		}
	}
	emoji := map[danbooru.Rating]string{
		danbooru.Explicit:     "🔞",
		danbooru.Questionable: "❓",
	}
	photo := &tele.Photo{
		File:       tele.FromURL(pic.URL),
		Caption:    emoji[pic.Rating],
		HasSpoiler: true,
	}
	return c.Send(photo, tele.ModeHTML)
}
