package context

import (
	tele "gopkg.in/telebot.v3"
	"nechego/handlers"
)

var handlerID = "handlerID"

func SetHandlerID(c tele.Context, h handlers.HandlerID) {
	c.Set(handlerID, h)
}

func HandlerID(c tele.Context) handlers.HandlerID {
	h, ok := c.Get(handlerID).(handlers.HandlerID)
	if !ok {
		return handlers.NoHandler
	}
	return h
}
