package token

import "fmt"

type Eblan struct{}

func (e Eblan) String() string {
	return "😸 Токен еблана"
}

type Admin struct{}

func (a Admin) String() string {
	return "👑 Токен администратора"
}

type Pair struct{}

func (p Pair) String() string {
	return "💘 Токен пары"
}

type Letter struct {
	Author string
	Text   string
}

func (l Letter) String() string {
	return fmt.Sprintf("✉️ Письмо (%s)", l.Author)
}
