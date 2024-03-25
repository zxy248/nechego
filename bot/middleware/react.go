package middleware

import (
	"math/rand/v2"

	tele "gopkg.in/zxy248/telebot.v3"
)

type RandomReact struct {
	Prob float64
}

func (m *RandomReact) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if rand.Float64() < m.Prob {
			go randomReact(c)
		}
		return next(c)
	}
}

var emojis = []tele.EmojiType{"👍", "👎", "❤", "🔥", "🥰", "👏", "😁",
	"🤔", "🤯", "😱", "🤬", "😢", "🎉", "🤩", "🤮", "💩", "🙏",
	"👌", "🕊", "🤡", "🥱", "🥴", "😍", "🐳", "❤‍🔥", "🌚", "🌭",
	"💯", "🤣", "⚡", "🍌", "🏆", "💔", "🤨", "😐", "🍓", "🍾",
	"💋", "🖕", "😈", "😴", "😭", "🤓", "👻", "👨‍💻", "👀", "🎃",
	"🙈", "😇", "😨", "🤝", "✍", "🤗", "🫡", "🎅", "🎄", "☃",
	"💅", "🤪", "🗿", "🆒", "💘", "🙉", "🦄", "😘", "💊", "🙊",
	"😎", "👾", "🤷‍♂", "🤷", "🤷‍♀", "😡"}

func randomReact(c tele.Context) {
	opts := tele.ReactionOptions{
		Reactions: []tele.Reaction{{
			Type:  "emoji",
			Emoji: emojis[rand.N(len(emojis))],
		}},
		Big: true,
	}
	c.Bot().React(c.Chat(), c.Message(), opts)
}
