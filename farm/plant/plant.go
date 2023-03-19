package plant

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math/rand"
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
	r := []Type{}
	for t := range data {
		if t != Void {
			r = append(r, t)
		}
	}
	return r
}()

var data = map[Type]struct {
	emoji string
	name  string
}{
	Void:       {"〰", "ничего"},
	Grapes:     {"🍇", "виноград"},
	Melon:      {"🍈", "дыня"},
	Watermelon: {"🍉", "арбуз"},
	Tangerine:  {"🍊", "мандарин"},
	Lemon:      {"🍋", "лимон"},
	Banana:     {"🍌", "банан"},
	Pineapple:  {"🍍", "ананас"},
	Mango:      {"🥭", "манго"},
	RedApple:   {"🍎", "яблоко"},
	GreenApple: {"🍏", "яблоко"},
	Pear:       {"🍐", "груша"},
	Peach:      {"🍑", "персик"},
	Cherry:     {"🍒", "вишня"},
	Strawberry: {"🍓", "клубника"},
	Blueberry:  {"🫐", "голубика"},
	Kiwi:       {"🥝", "киви"},
	Tomato:     {"🍅", "помидор"},
	Olive:      {"🫒", "олива"},
	Coconut:    {"🥥", "кокос"},
	Avocado:    {"🥑", "авокадо"},
	Eggplant:   {"🍆", "баклажан"},
	Potato:     {"🥔", "картофель"},
	Carrot:     {"🥕", "морковь"},
	Corn:       {"🌽", "кукуруза"},
	HotPepper:  {"🌶", "халапеньо"},
	BellPepper: {"🫑", "перец сладкий"},
	Cucumber:   {"🥒", "огурец"},
	Lettuce:    {"🥬", "салат"},
	Broccoli:   {"🥦", "брокколи"},
	Garlic:     {"🧄", "чеснок"},
	Onion:      {"🧅", "лук"},
	Mushroom:   {"🍄", "гриб"},
	Peanuts:    {"🥜", "арахис"},
	Beans:      {"🫘", "бобы"},
	Chestnut:   {"🌰", "фундук"},
	Laminaria:  {"🍥", "ламинария"},
}

func (t Type) String() string {
	return data[t].emoji
}

type Plant struct {
	Type
	Count int
}

func Random() *Plant {
	return &Plant{
		Type:  Types[rand.Intn(len(Types))],
		Count: 1,
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
