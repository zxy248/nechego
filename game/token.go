package game

type EblanToken struct{}

func (e EblanToken) String() string {
	return "😸 Токен еблана дня"
}

type AdminToken struct{}

func (a AdminToken) String() string {
	return "👑 Токен админа дня"
}

type PairToken struct{}

func (p PairToken) String() string {
	return "💘 Токен пары дня"
}
