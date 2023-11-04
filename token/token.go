package token

import (
	"nechego/modifier"
)

type Eblan struct{}

func (e Eblan) String() string {
	return "😸 Токен еблана"
}

func (e Eblan) Mod() (m *modifier.Mod, ok bool) {
	return &modifier.Mod{
		Emoji:       "😸",
		Multiplier:  -0.2,
		Description: "Вы еблан.",
	}, true
}

type Admin struct{}

func (a Admin) String() string {
	return "👑 Токен администратора"
}

func (a Admin) Mod() (m *modifier.Mod, ok bool) {
	return &modifier.Mod{
		Emoji:       "👑",
		Multiplier:  0.2,
		Description: "Вы администратор.",
	}, true
}

type Pair struct{}

func (p Pair) String() string {
	return "💘 Токен пары"
}

func (p Pair) Mod() (m *modifier.Mod, ok bool) {
	return &modifier.Mod{
		Emoji:       "💖",
		Multiplier:  0.1,
		Description: "У вас есть пара.",
	}, true
}
