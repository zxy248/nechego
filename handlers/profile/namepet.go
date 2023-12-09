package profile

import (
	"fmt"
	"nechego/game"
	"nechego/handlers"
	"nechego/pets"
	tu "nechego/teleutil"
	"nechego/valid"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type NamePet struct {
	Universe *game.Universe
}

var namePetRe = handlers.Regexp("^!назвать (.+)")

func (h *NamePet) Match(c tele.Context) bool {
	return namePetRe.MatchString(c.Text())
}

func (h *NamePet) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	pet, ok := user.Pet()
	if !ok {
		return c.Send(petNotFound)
	}
	n := petName(c.Text())
	if !valid.Name(n) {
		return c.Send(badPetName)
	}
	pet.Name = formatPetName(n)
	return c.Send(petNamed(pet), tele.ModeHTML)
}

const (
	petNotFound = "🐱 У вас нет питомца."
	badPetName  = "🐱 Такое имя не подходит для питомца."
)

func petNamed(p *pets.Pet) string {
	const format = "%s Вы назвали питомца <code>%s</code>."
	return fmt.Sprintf(format, p.Species.Emoji(), p.Name)
}

func petName(s string) string {
	return namePetRe.FindStringSubmatch(s)[1]
}

func formatPetName(s string) string {
	return strings.Title(s)
}
