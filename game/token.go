package game

type EblanToken struct{}

func (e EblanToken) String() string {
	return "😸 Жетон еблана дня"
}

type AdminToken struct{}

func (a AdminToken) String() string {
	return "👑 Жетон админа дня"
}

type PairToken struct {
	MateID int
}

func (p PairToken) String() string {
	return "💘 Жетон пары дня"
}
