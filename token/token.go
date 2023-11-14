package token

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
