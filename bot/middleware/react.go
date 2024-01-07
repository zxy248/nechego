package middleware

import (
	"math/rand"

	tele "gopkg.in/telebot.v3"
)

type RandomReact struct {
	Prob float64
}

func (m *RandomReact) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if rand.Float64() < m.Prob {
			go react(c)
		}
		return next(c)
	}
}

var emojis = []string{"👍", "👎", "❤", "🔥", "🥰", "👏", "😁", "🤔", "🤯", "😱",
	"🤬", "😢", "🎉", "🤩", "🤮", "💩", "🙏", "👌", "🕊", "🤡", "🥱", "🥴",
	"😍", "🐳", "❤‍🔥", "🌚", "🌭", "💯", "🤣", "⚡", "🍌", "🏆", "💔", "🤨",
	"😐", "🍓", "🍾", "💋", "🖕", "😈", "😴", "😭", "🤓", "👻", "👨‍💻", "👀",
	"🎃", "🙈", "😇", "😨", "🤝", "✍", "🤗", "🫡", "🎅", "🎄", "☃", "💅",
	"🤪", "🗿", "🆒", "💘", "🙉", "🦄", "😘", "💊", "🙊", "😎", "👾", "🤷‍♂",
	"🤷", "🤷‍♀", "😡"}

func react(c tele.Context) error {
	emoji := map[string]any{
		"type":  "emoji",
		"emoji": emojis[rand.Intn(len(emojis))],
	}
	params := map[string]any{
		"chat_id":    c.Chat().ID,
		"message_id": c.Message().ID,
		"reaction":   []any{emoji},
		"is_big":     true,
	}
	_, err := c.Bot().Raw("setMessageReaction", params)
	return err
}
