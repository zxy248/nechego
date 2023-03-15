package food

import (
	"fmt"
	"math/rand"
)

type Type int

const (
	Bread Type = iota
	ChickenLeg
	BigTasty
	BigMac
	Fries
	PizzaFourCheese
	PizzaPepperoni
	PizzaCheeseChicken
	Toast
	Shawarma
	SuperKontik
	AdrenalineRush
	Burn
	Ramen
	Hotdog
	RitterSport
	HotCat
	Jaguar
	Beer
	IceCream
	Juice
)

var beverages = map[Type]bool{
	AdrenalineRush: true,
	Burn:           true,
	HotCat:         true,
	Jaguar:         true,
	Beer:           true,
	Juice:          true,
}

func (t Type) Emoji() string      { return data[t].Emoji }
func (t Type) Nutrition() float64 { return data[t].Nutrition }
func (t Type) String() string     { return data[t].Description }
func (t Type) Beverage() bool     { return beverages[t] }

var data = map[Type]struct {
	Emoji       string
	Nutrition   float64
	Description string
}{
	Bread:              {"🍞", 0.08, "Хлеб"},
	ChickenLeg:         {"🍗", 0.12, "Куриная ножка"},
	BigTasty:           {"🍔", 0.16, "Биг Тейсти"},
	BigMac:             {"🍔", 0.14, "Биг Мак"},
	Fries:              {"🍟", 0.08, "Картофель фри"},
	PizzaFourCheese:    {"🍕", 0.16, "Пицца (4 сыра)"},
	PizzaPepperoni:     {"🍕", 0.16, "Пицца (пеперони)"},
	PizzaCheeseChicken: {"🍕", 0.16, "Пицца (сырный цыплёнок)"},
	Toast:              {"🥪", 0.10, "Бутерброд"},
	Shawarma:           {"🌯", 0.16, "Шаурма"},
	SuperKontik:        {"🍩", 0.10, "Супер-Контик"},
	AdrenalineRush:     {"🦎", 0.20, "Напиток Adrenaline Rush"},
	Burn:               {"🔥", 0.20, "Напиток Burn"},
	Ramen:              {"🍜", 0.20, "Доширак"},
	Hotdog:             {"🌭", 0.16, "Хот-дог"},
	RitterSport:        {"🍫", 0.16, "Риттер Спорт"},
	HotCat:             {"🐱", 0.20, "Напиток HotCat"},
	Jaguar:             {"🐾", 0.20, "Напиток Jaguar"},
	Beer:               {"🍺", 0.10, "Пиво"},
	IceCream:           {"🍦", 0.08, "Мороженое"},
	Juice:              {"🧃", 0.08, "Сок"},
}

type Food struct {
	Type
}

func (f Food) String() string {
	return fmt.Sprintf("%s %s", f.Type.Emoji(), f.Type)
}

func Random() *Food {
	return &Food{Type: Type(rand.Intn(len(data)))}
}
