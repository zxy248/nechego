package app

import (
	"errors"
	"fmt"
	"nechego/fishing"
	"nechego/model"
	"nechego/service"

	tele "gopkg.in/telebot.v3"
)

const (
	foodEaten        = HTML("🍊 Вы поели.")
	boughtFishingRod = Response("🎣 Вы приобрели удочку за %s")
	notEnoughFood    = UserError("Недостаточно еды.")
	youAreFull       = UserError("Вы не хотите есть.")
	alreadyCanFish   = UserError("Вы уже приобрели удочку.")
	buyFishingRod    = UserError("Приобретите удочку, прежде чем рыбачить.")
)

// !еда
func (a *App) handleEatFood(c tele.Context) error {
	user := getUser(c)
	energyRestored, err := a.service.EatFish(user)
	if err != nil {
		if errors.Is(err, service.ErrNotHungry) {
			return respondUserError(c, youAreFull)
		}
		if errors.Is(err, service.ErrNotEnoughFish) {
			return respondUserError(c, notEnoughFood)
		}
		return respondInternalError(c, err)
	}
	return respond(c, eatFoodResponse(user, energyRestored))
}

func eatFoodResponse(u model.User, energyRestored int) Response {
	return Response(joinSections(
		string(foodEaten),
		string(energyRemaining(u.Energy+energyRestored))))
}

// !удочка
func (a *App) handleFishingRod(c tele.Context) error {
	if err := a.service.BuyFishingRod(getUser(c)); err != nil {
		if errors.Is(err, service.ErrAlreadyFisher) {
			return respondUserError(c, alreadyCanFish)
		}
		var moneyErr service.NotEnoughMoneyError
		if errors.As(err, &moneyErr) {
			return respondUserError(c, notEnoughMoneyDelta.Fill(formatMoney(moneyErr.Delta)))
		}
		return respondInternalError(c, err)
	}
	return respond(c, boughtFishingRod.Fill(formatMoney(a.service.Config.FishingRodPrice)))
}

const (
	catchFish = Response("<i>%s получает рыбу: <code>%s</code></i>")
	foodFish  = HTML("<i>🍊 Вы отложили улов на %s.</i>")
)

// !рыбалка
func (a *App) handleFishing(c tele.Context) error {
	user := getUser(c)
	session, err := a.service.Fish(user)
	if err != nil {
		if errors.Is(err, service.ErrNotFisher) {
			return respondUserError(c, buyFishingRod)
		}
		if errors.Is(err, service.ErrNotEnoughEnergy) {
			return respondUserError(c, notEnoughEnergy)
		}
		return respondInternalError(c, err)
	}
	return respond(c, a.fishingResponse(user, session))
}

func (a *App) fishingResponse(u model.User, s fishing.Session) Response {
	out := s.Outcome.String()
	sections := []string{out, string(catchFish.Fill(a.mustMentionUser(u), s.Fish))}
	if s.Fish.Light() {
		sections = append(sections, fmt.Sprintf(string(foodFish), randomMeal()))
	}
	if s.Success() {
		out = joinSections(sections...)
	}
	return Response(out)
}

// !рыба
func (a *App) handleFish(c tele.Context) error {
	user := getUser(c)
	fishes, err := a.service.FreshFish(user)
	if err != nil {
		return respondInternalError(c, err)
	}
	return respond(c, freshFish.Fill(a.mustMentionUser(user), formatFish(fishes)))
}

// !продажа
func (a *App) handleSellFish(c tele.Context) error {
	fishes, err := a.service.SellFish(getUser(c))
	if err != nil {
		return respondInternalError(c, err)
	}
	price := fishes.Price()
	if price > 0 {
		return respond(c, soldFish.Fill(formatMoney(price)))
	}
	return respondUserError(c, noFish)
}

const (
	freshFish    = Response("<b>🐟 Улов %s</b>\n%s")
	freezerFish  = Response("<b>🧊 Холодильник %s</b>\n%s")
	fishFrozen   = Response("❄️ Рыба заморожена.")
	fishUnfrozen = Response("💧 Рыба разморожена.")
	soldFish     = Response("🐟 Рыбы продано на %s")
	noFish       = UserError("У вас нет свежей рыбы.")
)

func (a *App) handleFreeze(c tele.Context) error {
	a.service.FreezeFish(getUser(c))
	return respond(c, fishFrozen)
}

func (a *App) handleUnfreeze(c tele.Context) error {
	a.service.UnfreezeFish(getUser(c))
	return respond(c, fishUnfrozen)
}

// !холодильник
func (a *App) handleFreezer(c tele.Context) error {
	user := getUser(c)
	fishes, err := a.service.Freezer(user)
	if err != nil {
		return respondInternalError(c, err)
	}
	return respond(c, freezerFish.Fill(a.mustMentionUser(user), formatFish(fishes)))
}

func formatFish(f fishing.Fishes) HTML {
	lines := []string{}
	for _, ff := range f {
		lines = append(lines, "<code>"+ff.String()+"</code>")
	}
	sections := []string{string(itemize(lines...))}
	if len(f) > 0 {
		sections = append(sections, string(formatFishSum(f)))
	}
	return HTML(joinSections(sections...))
}

func formatFishSum(f fishing.Fishes) HTML {
	lines := []string{
		string("<i>Стоимость: </i>" + formatMoney(f.Price())),
		string("<i>Вес: </i>" + formatWeight(f.Weight())),
	}
	return HTML(joinLines(lines...))
}
