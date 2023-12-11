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
	ScandinavianBurger
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
	Croissant
	Egg
	Cheese
	Bacon
	Salad
	Can
	Sushi
	Cake
	Lolipop
	Candy
	Popcorn
	Cookie
	Coffee
	Dumpling
)

var beverages = map[Type]bool{
	AdrenalineRush: true,
	Burn:           true,
	HotCat:         true,
	Jaguar:         true,
	Beer:           true,
	Juice:          true,
	Coffee:         true,
}

func (t Type) Emoji() string      { return data[t].emoji }
func (t Type) Nutrition() float64 { return data[t].nutrition }
func (t Type) String() string     { return data[t].name }
func (t Type) Price() float64     { return data[t].price }
func (t Type) Beverage() bool     { return beverages[t] }

var data = map[Type]struct {
	emoji     string
	nutrition float64
	name      string
	price     float64
}{
	Bread:              {"🍞", 0.15, "Хлеб", 50},
	ChickenLeg:         {"🍗", 0.15, "Куриная ножка", 100},
	BigTasty:           {"🍔", 0.25, "Биг Тейсти", 250},
	BigMac:             {"🍔", 0.20, "Биг Мак", 150},
	ScandinavianBurger: {"🍔", 0.20, "Скандинавский бургер", 250},
	Fries:              {"🍟", 0.15, "Картофель фри", 100},
	PizzaFourCheese:    {"🍕", 0.20, "Пицца (4 сыра)", 150},
	PizzaPepperoni:     {"🍕", 0.20, "Пицца (пеперони)", 150},
	PizzaCheeseChicken: {"🍕", 0.20, "Пицца (сырный цыплёнок)", 150},
	Toast:              {"🥪", 0.15, "Бутерброд", 100},
	Shawarma:           {"🌯", 0.25, "Шаурма", 200},
	SuperKontik:        {"🍩", 0.10, "Супер-Контик", 50},
	AdrenalineRush:     {"🦎", 0.20, "Напиток Adrenaline Rush", 100},
	Burn:               {"🔥", 0.20, "Напиток Burn", 100},
	Ramen:              {"🍜", 0.20, "Доширак", 50},
	Hotdog:             {"🌭", 0.20, "Хот-дог", 100},
	RitterSport:        {"🍫", 0.20, "Риттер Спорт", 100},
	HotCat:             {"🐱", 0.20, "Напиток HotCat", 50},
	Jaguar:             {"🐾", 0.20, "Напиток Jaguar", 100},
	Beer:               {"🍺", 0.10, "Пиво", 75},
	IceCream:           {"🍦", 0.15, "Мороженое", 50},
	Juice:              {"🧃", 0.10, "Сок", 30},
	Croissant:          {"🥐", 0.15, "Круасан", 50},
	Egg:                {"🥚", 0.15, "Яйцо", 30},
	Cheese:             {"🧀", 0.15, "Сыр", 150},
	Bacon:              {"🥓", 0.15, "Бекон", 100},
	Salad:              {"🥗", 0.05, "Салат", 300},
	Can:                {"🥫", 0.20, "Консерва", 100},
	Sushi:              {"🍣", 0.15, "Суши", 300},
	Cake:               {"🍰", 0.15, "Торт", 150},
	Lolipop:            {"🍭", 0.10, "Леденец", 50},
	Candy:              {"🍬", 0.05, "Конфета", 10},
	Popcorn:            {"🍿", 0.10, "Попкорн", 100},
	Cookie:             {"🍪", 0.10, "Печенье", 50},
	Coffee:             {"☕️", 0.15, "MacCoffee", 100},
	Dumpling:           {"🥟", 0.20, "Чебупели", 150},
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
