package modifiers

import (
	"golang.org/x/exp/slices"
)

var (
	NoModifier            = &Modifier{"", +0.00, "", ""}
	AdminModifier         = &Modifier{"👑", +0.20, "Вы ощущаете власть над остальными.", "администратор"}
	EblanModifier         = &Modifier{"😸", -0.20, "Вы чувствуете себя оскорбленным.", "еблан"}
	MuchEnergyModifier    = &Modifier{"🍥", +0.20, "Вы хорошо поели.", ""}
	FullEnergyModifier    = &Modifier{"⚡️", +0.10, "Вы полны сил.", ""}
	NoEnergyModifier      = &Modifier{"😣", -0.25, "Вы чувствуете себя уставшим.", ""}
	TerribleLuckModifier  = &Modifier{"☠️", -0.50, "Вас преследуют неудачи.", ""}
	BadLuckModifier       = &Modifier{"", -0.10, "Вам не везет.", ""}
	GoodLuckModifier      = &Modifier{"🤞", +0.10, "Вам везет.", ""}
	ExcellentLuckModifier = &Modifier{"🍀", +0.30, "Сегодня ваш день.", ""}
	RichModifier          = &Modifier{"🎩", +0.05, "Вы богаты.", "магнат"}
	PoorModifier          = &Modifier{"", -0.05, "Вы бедны.", ""}
	FisherModifier        = &Modifier{"🎣", +0.05, "Вы можете рыбачить.", ""}
	DebtorModifier        = &Modifier{"💳", -0.25, "У вас есть кредит.", ""}
)

type Modifier struct {
	Icon        string
	Multiplier  float64
	Description string
	Title       string
}

type Set []*Modifier

func NewSet() *Set {
	return &Set{}
}

func (s *Set) Present(m *Modifier) bool {
	return slices.Contains(*s, m)
}

func (s *Set) Add(m *Modifier) {
	if slices.Contains(*s, m) {
		return
	}
	if m == NoModifier {
		return
	}
	*s = append(*s, m)
}

func (s *Set) List() []*Modifier {
	return *s
}

func (s *Set) Sum() float64 {
	sum := 0.0
	for _, m := range s.List() {
		sum += m.Multiplier
	}
	return sum
}

func (s *Set) Icons() []string {
	out := []string{}
	for _, m := range s.List() {
		if m.Icon == "" {
			continue
		}
		out = append(out, m.Icon)
	}
	return out
}

func (s *Set) Descriptions() []string {
	out := []string{}
	for _, m := range s.List() {
		out = append(out, m.Description)
	}
	return out
}

func (s *Set) Titles() []string {
	out := []string{}
	for _, m := range s.List() {
		if m.Title == "" {
			continue
		}
		out = append(out, m.Title)
	}
	return out
}
