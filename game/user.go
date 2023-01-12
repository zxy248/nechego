package game

import (
	"fmt"
	"math/rand"
	"nechego/elo"
	"nechego/fishing"
	"nechego/modifier"
	"nechego/pets"
	"time"
)

type Gender int

const (
	GenderUnknown Gender = iota
	GenderMale
	GenderFemale
	GenderTrans
)

type User struct {
	TUID      int64
	Energy    int
	EnergyCap int
	Rating    float64
	Messages  int
	Banned    bool
	Birthday  time.Time
	Gender    Gender
	Status    string
	Inventory *Items
}

func NewUser(tuid int64) *User {
	return &User{
		TUID:      tuid,
		EnergyCap: 5,
		Rating:    1500,
		Inventory: NewItems(),
	}
}

func (u *User) SpendEnergy(e int) bool {
	if u.Energy < e {
		return false
	}
	u.Energy -= e
	return true
}

func (u *User) RestoreEnergy(e int) {
	u.Energy += e
	if u.Energy > u.EnergyCap {
		u.Energy = u.EnergyCap
	}
}

func (u *User) SpendMoney(n int) bool {
	u.Stack()
	return u.spendWallet(n) || u.spendCash(n)
}

func (u *User) spendWallet(n int) bool {
	w, ok := u.Wallet()
	if !ok {
		return false
	}
	if w.Money < n {
		return false
	}
	w.Money -= n
	return true
}

func (u *User) spendCash(n int) bool {
	c, ok := u.Cash()
	if !ok {
		return false
	}
	if c.Money < n {
		return false
	}
	c.Money -= n
	return true
}

func (u *User) AddMoney(n int) {
	u.Inventory.Add(&Item{
		Type:         ItemTypeCash,
		Transferable: true,
		Value:        &Cash{Money: n},
	})
}

func (u *User) Stack() bool {
	t := 0
	for _, x := range u.Inventory.list() {
		if cash, ok := x.Value.(*Cash); ok {
			// TODO: possible optimization (from n^2 to n)
			// Filter(func (i *Item) (keep bool))
			u.Inventory.Remove(x)
			t += cash.Money
		}
	}
	if t == 0 {
		return false
	}
	wallet, ok := u.Wallet()
	if !ok {
		u.Inventory.Add(&Item{
			Type:         ItemTypeCash,
			Transferable: true,
			Value:        &Cash{Money: t},
		})
		return true
	}
	wallet.Money += t
	return true
}

func (u *User) Total() int {
	t := 0
	for _, v := range u.Inventory.list() {
		switch x := v.Value.(type) {
		case *Cash:
			t += x.Money
		case *Wallet:
			t += x.Money
		case *CreditCard:
			t += x.Money
		case *Debt:
			t -= x.Money
		}
	}
	return t
}

func (u *User) InDebt() bool {
	for _, v := range u.Inventory.list() {
		switch v.Value.(type) {
		case *Debt:
			return true
		}
	}
	return false
}

// TODO: remove is prefixes
func (u *User) IsEblan() bool {
	for _, v := range u.Inventory.list() {
		switch v.Value.(type) {
		case *EblanToken:
			return true
		}
	}
	return false
}

func (u *User) IsAdmin() bool {
	for _, v := range u.Inventory.list() {
		switch v.Value.(type) {
		case *AdminToken:
			return true
		}
	}
	return false
}

func (u *User) IsPair() bool {
	for _, v := range u.Inventory.list() {
		switch v.Value.(type) {
		case *PairToken:
			return true
		}
	}
	return false
}

func (u *User) Pet() (p *pets.Pet, ok bool) {
	for _, x := range u.Inventory.list() {
		if p, ok = x.Value.(*pets.Pet); ok {
			return
		}
	}
	return nil, false
}

func (u *User) Rich() bool {
	return u.Total() > 1000000
}

func (u *User) Poor() bool {
	return u.Total() < 3000
}

func (u *User) Luck() float64 {
	day := time.Now().Truncate(time.Hour * 24)
	const prime = 1000003
	return float64((uint64(day.UnixNano())+uint64(u.TUID))%prime) / prime
}

func (u *User) Eat(i *Item) bool {
	switch x := i.Value.(type) {
	case *fishing.Fish:
		u.Inventory.Remove(i)
		e := 1
		if x.Heavy() {
			e = 2
		}
		u.RestoreEnergy(e)
		return true
	case *pets.Pet:
		u.Inventory.Remove(i)
		u.RestoreEnergy(1 + rand.Intn(2))
		return true
	}
	return false
}

func (u *User) Sell(i *Item) (profit int, ok bool) {
	if !i.Transferable {
		return 0, false
	}
	if ok = u.Inventory.Remove(i); !ok {
		return 0, false
	}
	switch x := i.Value.(type) {
	case *fishing.Fish:
		n := int(x.Price())
		u.AddMoney(n)
		return n, true
	default:
		// cannot sell; return item back
		u.Inventory.Add(i)
	}
	return 0, false
}

func (u *User) Strength() float64 {
	return 1.0 + u.Modset().Sum()
}

func (u *User) Fight(opponent *User) (winner, loser *User, r float64) {
	if u == opponent {
		panic("user cannot be an opponent to themself")
	}
	if u.power() > opponent.power() {
		winner, loser = u, opponent
	} else {
		winner, loser = opponent, u
	}
	r = elo.EloDelta(winner.Rating, loser.Rating, elo.KDefault, elo.ScoreWin)
	winner.Rating += r
	loser.Rating -= r
	return
}

func (u *User) power() float64 {
	return u.Strength() * rand.Float64()
}

func (u *User) Modset() modifier.Set {
	set := modifier.Set{}
	if u.IsAdmin() {
		set.Add(modifier.Admin)
	}
	if u.IsEblan() {
		set.Add(modifier.Eblan)
	}
	if u.Energy == 0 {
		set.Add(modifier.NoEnergy)
	}
	if u.Energy == u.EnergyCap {
		set.Add(modifier.FullEnergy)
	}
	if u.Energy > u.EnergyCap {
		set.Add(modifier.MuchEnergy)
	}
	if u.Rich() {
		set.Add(modifier.Rich)
	}
	if u.Poor() {
		set.Add(modifier.Poor)
	}
	if u.InDebt() {
		set.Add(modifier.Debtor)
	}

	switch l := u.Luck(); {
	case l < 0.01:
		set.Add(modifier.TerribleLuck)
	case l < 0.02:
		set.Add(modifier.ExcellentLuck)
	case l < 0.10:
		set.Add(modifier.BadLuck)
	case l < 0.18:
		set.Add(modifier.GoodLuck)
	}

	if _, ok := u.FishingRod(); ok {
		set.Add(modifier.Fisher)
	}
	if p, ok := u.Pet(); ok {
		q := 0.05
		switch p.Species.Quality() {
		case pets.Rare:
			q = 0.10
		case pets.Exotic:
			q = 0.15
		case pets.Legendary:
			q = 0.20
		}
		r := ""
		if p.Species.Quality() != pets.Common {
			r = fmt.Sprintf("%s ", p.Species.Quality())
		}
		set.Add(&modifier.Mod{
			Emoji:       p.Species.Emoji(),
			Multiplier:  q,
			Description: fmt.Sprintf("У вас есть %sпитомец: <code>%s</code>", r, p),
		})
	}
	return set
}
