package app

import (
	"errors"
	"nechego/service"

	tele "gopkg.in/telebot.v3"
)

const (
	yourPet   = Response("🐾 Ваш питомец — %s")
	petBought = Response(`🐾 Вы приобрели питомца за %s.

%s Это <code>%s (%s)</code>.

<i>Используйте команду <code>!назвать &lt;имя&gt;</code>.</i>`)
	petNamed        = Response("🐾 Вы назвали питомца <b>%s</b>.")
	petDropped      = Response("😥 Вы выкинули своего питомца.")
	petBadName      = UserError("Такое имя не подходит для питомца.")
	petAlreadyNamed = UserError("У вашего питомца уже есть имя.")
	petAlreadyTaken = UserError("У вас уже есть питомец.")
	youHaveNoPet    = UserError("У вас нет питомца.")
	nameYourPet     = UserError("Назовите вашего питомца.")
)

func (a *App) handlePet(c tele.Context) error {
	pet, err := a.service.GetPet(getUser(c))
	if err != nil {
		if errors.Is(err, service.ErrNoPet) {
			return respondUserError(c, youHaveNoPet)
		}
		return respondInternalError(c, err)
	}
	if !pet.HasName() {
		return respondUserError(c, nameYourPet)
	}
	return respond(c, yourPet.Fill(formatPet(pet)))
}

func (a *App) handleBuyPet(c tele.Context) error {
	pet, err := a.service.BuyPet(getUser(c))
	if err != nil {
		var moneyErr service.NotEnoughMoneyError
		if errors.As(err, &moneyErr) {
			return respondUserError(c, notEnoughMoneyDelta.Fill(formatMoney(moneyErr.Delta)))
		}
		if errors.Is(err, service.ErrPetAlreadyTaken) {
			return respondUserError(c, petAlreadyTaken)
		}
		return respondInternalError(c, err)
	}
	return respond(c, petBought.Fill(
		formatMoney(a.service.Config.PetPrice),
		pet.Species.Emoji(),
		pet.Species.String(),
		pet.Gender.Emoji(),
	))
}

func (a *App) handleNamePet(c tele.Context) error {
	name := getMessage(c).Argument()
	if err := a.service.NamePet(getUser(c), name); err != nil {
		if errors.Is(err, service.ErrPetAlreadyNamed) {
			return respondUserError(c, petAlreadyNamed)
		}
		if errors.Is(err, service.ErrPetBadName) {
			return respondUserError(c, petBadName)
		}
		return respondInternalError(c, err)
	}
	return respond(c, petNamed.Fill(name))
}

func (a *App) handleDropPet(c tele.Context) error {
	if err := a.service.DropPet(getUser(c)); err != nil {
		if errors.Is(err, service.ErrNoPet) {
			return respondUserError(c, youHaveNoPet)
		}
		return respondInternalError(c, err)
	}
	return respond(c, petDropped)
}
