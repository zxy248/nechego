package app

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

type Response string

func (r Response) Send(b *tele.Bot, to tele.Recipient, opt *tele.SendOptions) (*tele.Message, error) {
	return b.Send(to, string(r), opt)
}

func (r Response) Fill(a ...any) Response {
	return Response(fmt.Sprintf(string(r), a...))
}

type UserError string

func (e UserError) Send(b *tele.Bot, to tele.Recipient, opt *tele.SendOptions) (*tele.Message, error) {
	return Response(e).Send(b, to, opt)
}

func (e UserError) Fill(a ...any) UserError {
	return UserError(Response(e).Fill(a...))
}

func respond(c tele.Context, r Response, opts ...interface{}) error {
	return c.Send(r, append(opts, tele.ModeHTML)...)
}

func respondUserError(c tele.Context, err UserError) error {
	return respond(c, Response(formatWarning(string(err))))
}

func respondInternalError(c tele.Context, err error) error {
	send := respond(c, Response(formatError("Ошибка сервера")))
	return serverError{send, err}
}

func respondVideo(c tele.Context, path string) error {
	return c.Send(&tele.Video{File: tele.FromDisk(path)})
}

type serverError struct {
	send, actual error
}

func (e serverError) Error() string {
	return fmt.Sprintf("send error: %s; actual error: %s", e.send, e.actual)
}
