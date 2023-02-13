package middleware

import (
	"nechego/bot/context"
	"nechego/handlers"
	"time"

	tele "gopkg.in/telebot.v3"
)

func deleteMessage(c tele.Context, m *tele.Message) {
	f := func() { c.Bot().Delete(m) }
	switch context.HandlerID(c) {
	case
		handlers.NoHandler,
		handlers.WhoHandler,
		handlers.ListHandler,
		handlers.TopHandler,
		handlers.InfaHandler,
		handlers.PicHandler,
		handlers.BasiliHandler,
		handlers.CasperHandler,
		handlers.ZeusHandler,
		handlers.CatHandler,
		handlers.AnimeHandler,
		handlers.FurryHandler,
		handlers.FlagHandler,
		handlers.PersonHandler,
		handlers.HorseHandler,
		handlers.ArtHandler,
		handlers.CarHandler,
		handlers.SoyHandler,
		handlers.DanbooruHandler,
		handlers.FapHandler:
		return
	case
		handlers.SendSMSHandler,
		handlers.SpamHandler:
		if m.Sender.IsBot {
			time.AfterFunc(7*time.Second, f)
		} else {
			f()
		}
	default:
		time.AfterFunc(10*time.Minute, f)
	}
}

type deleteMessageContext struct {
	tele.Context
}

func (c deleteMessageContext) Send(what interface{}, opts ...interface{}) error {
	userMsg := c.Message()
	botMsg, err := c.Bot().Send(c.Recipient(), what, opts...)
	if err != nil {
		return err
	}
	deleteMessage(c, userMsg)
	deleteMessage(c, botMsg)
	return nil
}

type DeleteMessage struct{}

func (m *DeleteMessage) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		return next(deleteMessageContext{c})
	}
}
