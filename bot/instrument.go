package main

import (
	"context"
	"reflect"
	"time"

	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"
	tele "gopkg.in/zxy248/telebot.v3"
)

type InstrumentedHandler struct {
	Queries *data.Queries
	Handler
}

func (h *InstrumentedHandler) Match(c tele.Context) bool {
	return h.Handler.Match(c)
}

func (h *InstrumentedHandler) Handle(c tele.Context) error {
	t1 := time.Now()
	err := h.Handler.Handle(c)
	t2 := time.Now()

	messageID := c.Get(handlers.MessageIDKey).(int64)
	handlerName := reflect.TypeOf(h.Handler).String()
	handlerTime := t2.Sub(t1)
	handlerError := ""
	if err != nil {
		handlerError = err.Error()
	}

	ctx := context.Background()
	arg := data.InstrumentMessageParams{
		MessageID: messageID,
		Handler:   handlerName,
		Time:      handlerTime,
		Error:     handlerError,
	}
	if err := h.Queries.InstrumentMessage(ctx, arg); err != nil {
		return err
	}
	return err
}
