package token

type Eblan struct{}

func (e Eblan) String() string {
	return "😸 Токен еблана дня"
}

type Admin struct{}

func (a Admin) String() string {
	return "👑 Токен админа дня"
}

type Pair struct{}

func (p Pair) String() string {
	return "💘 Токен пары дня"
}

type Dice struct{}

func (d Dice) String() string {
	return "🎲 Игральные кости"
}
