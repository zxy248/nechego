package profile

import (
	"nechego/avatar"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Avatar struct {
	Avatars   *avatar.Storage
	MaxWidth  int
	MaxHeight int
}

var avatarRe = handlers.Regexp("^!ава")

func (h *Avatar) Match(c tele.Context) bool {
	return avatarRe.MatchString(c.Text())
}

func (h *Avatar) Handle(c tele.Context) error {
	id := c.Sender().ID
	if r, ok := tu.Reply(c); ok {
		id = r.ID
	}
	if p := c.Message().Photo; p != nil {
		if !h.validSize(p) {
			const f = "📷 Максимальный размер аватара %dx%d пикселей."
			return c.Send(f, h.MaxWidth, h.MaxHeight)
		}
		if err := h.setPhoto(c, id, p); err != nil {
			return err
		}
		return c.Send("📸 Аватар установлен.")
	}
	if p, ok := h.Avatars.Get(id); ok {
		return c.Send(p)
	}
	return c.Send("📷 Прикрепите изображение.")
}

func (h *Avatar) validSize(p *tele.Photo) bool {
	return p.Width <= h.MaxWidth && p.Height <= h.MaxHeight
}

func (h *Avatar) setPhoto(c tele.Context, id int64, p *tele.Photo) error {
	f, err := c.Bot().File(&p.File)
	if err != nil {
		return err
	}
	defer f.Close()
	return h.Avatars.Set(id, f)
}
