package plant

import (
	"fmt"
	"math"
	"math/rand"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Type int

const (
	Void Type = iota
	Grapes
	Melon
	Watermelon
	Tangerine
	Lemon
	Banana
	Pineapple
	Mango
	RedApple
	GreenApple
	Pear
	Peach
	Cherry
	Strawberry
	Blueberry
	Kiwi
	Tomato
	Olive
	Coconut
	Avocado
	Eggplant
	Potato
	Carrot
	Corn
	HotPepper
	BellPepper
	Cucumber
	Lettuce
	Broccoli
	Garlic
	Onion
	Mushroom
	Peanuts
	Beans
	Chestnut
	Laminaria
)

// All types except Void.
var Types = func() []Type {
	var r []Type
	for t := range data {
		if t != Void {
			r = append(r, t)
		}
	}
	return r
}()

var data = map[Type]struct {
	emoji string
	mult  int
	name  string
	price float64
}{
	Void:       {"〰", 0, "ничего", 0},
	Grapes:     {"🍇", 2, "виноград", 300},
	Melon:      {"🍈", 2, "дыня", 200},
	Watermelon: {"🍉", 2, "арбуз", 200},
	Tangerine:  {"🍊", 3, "мандарин", 150},
	Lemon:      {"🍋", 3, "лимон", 125},
	Banana:     {"🍌", 4, "банан", 125},
	Pineapple:  {"🍍", 2, "ананас", 400},
	Mango:      {"🥭", 2, "манго", 300},
	RedApple:   {"🍎", 4, "яблоко", 100},
	GreenApple: {"🍏", 4, "яблоко", 100},
	Pear:       {"🍐", 3, "груша", 200},
	Peach:      {"🍑", 3, "персик", 150},
	Cherry:     {"🍒", 2, "вишня", 300},
	Strawberry: {"🍓", 2, "клубника", 500},
	Blueberry:  {"🫐", 2, "голубика", 700},
	Kiwi:       {"🥝", 3, "киви", 125},
	Tomato:     {"🍅", 4, "помидор", 200},
	Olive:      {"🫒", 2, "олива", 400},
	Coconut:    {"🥥", 2, "кокос", 300},
	Avocado:    {"🥑", 2, "авокадо", 500},
	Eggplant:   {"🍆", 2, "баклажан", 200},
	Potato:     {"🥔", 5, "картофель", 25},
	Carrot:     {"🥕", 5, "морковь", 30},
	Corn:       {"🌽", 5, "кукуруза", 50},
	HotPepper:  {"🌶", 2, "халапеньо", 500},
	BellPepper: {"🫑", 3, "перец сладкий", 200},
	Cucumber:   {"🥒", 4, "огурец", 250},
	Lettuce:    {"🥬", 3, "салат", 200},
	Broccoli:   {"🥦", 2, "брокколи", 300},
	Garlic:     {"🧄", 2, "чеснок", 400},
	Onion:      {"🧅", 5, "лук", 40},
	Mushroom:   {"🍄", 3, "гриб", 200},
	Peanuts:    {"🥜", 2, "арахис", 500},
	Beans:      {"🫘", 4, "бобы", 200},
	Chestnut:   {"🌰", 2, "фундук", 600},
	Laminaria:  {"🍥", 2, "ламинария", 300},
}

func (t Type) Emoji() string  { return data[t].emoji }
func (t Type) Yield() int     { return data[t].mult }
func (t Type) String() string { return data[t].name }
func (t Type) Price() float64 { return data[t].price }

type Plant struct {
	Type
	Count int
}

func Random() *Plant {
	c := math.Abs(rand.NormFloat64() * 3)
	return &Plant{
		Type:  Types[rand.Intn(len(Types))],
		Count: 1 + int(c),
	}
}

func (p *Plant) String() string {
	name := cases.Title(language.Russian).String(data[p.Type].name)
	var count string
	if p.Count > 1 {
		count = fmt.Sprintf(" (%d шт.)", p.Count)
	}
	return fmt.Sprintf("%s %s%s", data[p.Type].emoji, name, count)
}
