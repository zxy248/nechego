package main

import (
	"log"
	"nechego/handlers"
	"os"
	"regexp"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Command struct {
	re      *regexp.Regexp
	handler tele.HandlerFunc
}

func NewCommand(re string, handler any) *Command {
	var f tele.HandlerFunc
	switch h := handler.(type) {
	case func(tele.Context) error:
		f = h
	case handlers.Handler:
		f = handlers.Func(h)
	default:
		panic("bad handler type: " + re)
	}
	return &Command{
		re:      regexp.MustCompile(re),
		handler: f,
	}
}

func (cmd *Command) Match(s string) bool {
	return cmd.re.MatchString(s)
}

func (cmd *Command) Handle(c tele.Context) error {
	return cmd.handler(c)
}

type Router struct {
	commands []*Command
}

func (r *Router) Register(h *Command) {
	r.commands = append(r.commands, h)
}

func (r *Router) OnText(c tele.Context) error {
	for _, cmd := range r.commands {
		if cmd.Match(c.Message().Text) {
			return cmd.Handle(c)
		}
	}
	return nil
}

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("NECHEGO_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	r := &Router{}
	handlers := [...]*Command{
		NewCommand("^!мыш", &handlers.Mouse{Path: "data/mouse.mp4"}),
		NewCommand("^!тикток", &handlers.Tiktok{Path: "data/tiktok/"}),
		NewCommand("^!игр", handlers.HandleGame),
	}
	for _, h := range handlers {
		r.Register(h)
	}

	bot.Handle(tele.OnText, r.OnText)
	bot.Start()
}
