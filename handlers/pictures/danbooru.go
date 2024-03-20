package pictures

import (
	"math/rand/v2"
	"time"

	"github.com/zxy248/nechego/handlers"
	"github.com/zxy248/nechego/handlers/pictures/danbooru"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Danbooru struct{}

func (h *Danbooru) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!данбору")
}

func (h *Danbooru) Handle(c tele.Context) error {
	pic := <-danbooruPictures
	photo := &tele.Photo{File: tele.FromURL(pic.URL)}
	if pic.Rating == danbooru.Explicit {
		photo.Caption = warningNSFW()
		photo.HasSpoiler = true
	}
	return c.Send(photo, tele.ModeHTML)
}

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

func warningNSFW() string {
	s := [...]string{
		"🔞 Осторожно! Только для взрослых.",
		"<i>Содержимое предназначено для просмотра лицами старше 18 лет.</i>",
		"<b>ВНИМАНИЕ!</b> Вы увидите фотографии взрослых голых женщин. Будьте сдержанны.",
	}
	return s[rand.N(len(s))]
}

var danbooruPictures = func() chan *danbooru.Picture {
	const workers = 4
	const size = 16

	pics := make(chan *danbooru.Picture, size)
	for range workers {
		go func() {
			for {
				pics <- danbooruPicture()
			}
		}()
	}
	return pics
}()

func danbooruPicture() *danbooru.Picture {
	const timeout = 2 * time.Second
	const score = 50

	pic, err := danbooru.Get()
	if err != nil {
		time.Sleep(timeout)
		return danbooruPicture()
	}
	if pic.Score < score {
		return danbooruPicture()
	}
	return pic
}
