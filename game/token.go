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

type Dice struct{}

func (d Dice) String() string {
	return "🎲 Игральные кости"
}

func (u *User) Dice() (d *Dice, ok bool) {
	for _, x := range u.Inventory.normalize() {
		if d, ok = x.Value.(*Dice); ok {
			return
		}
	}
	return nil, false
}
