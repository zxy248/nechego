package app

import (
	"errors"
	"fmt"
	"math/rand"
	"nechego/model"

	tele "gopkg.in/telebot.v3"
)

const (
	notEnoughFish = "🐟 Недостаточно рыбы."
	fishEaten     = "🐟 Вы съели рыбу."
	youAreFull    = "🐟 Вы не хотите есть."
)

func (a *App) handleEatFish(c tele.Context) error {
	user := getUser(c)
	if hasFullEnergy(user) {
		return c.Send(youAreFull)
	}
	ok := a.model.EatFish(user, energyDelta, energyCap)
	if !ok {
		return c.Send(notEnoughFish)
	}
	return c.Send(fishEaten)
}

const (
	boughtFishingRod         = "🎣 Вы приобрели удочку за %s"
	alreadyCanFish           = "Вы уже приобрели удочку."
	notEnoughMoneyFishingRod = "Вам не хватает %s"
)

// !удочка
func (a *App) handleFishingRod(c tele.Context) error {
	user := getUser(c)
	if user.Fisher {
		return userError(c, alreadyCanFish)
	}
	ok := a.model.UpdateMoney(user, -fishingRodPrice)
	if !ok {
		return userErrorMarkdown(c, fmt.Sprintf(notEnoughMoneyFishingRod,
			formatMoney(fishingRodPrice-user.Balance)))
	}
	a.model.AllowFishing(user)
	return c.Send(fmt.Sprintf(boughtFishingRod, formatMoney(fishingRodPrice)),
		tele.ModeMarkdownV2)
}

type catchFishType int

const (
	catchFishSell catchFishType = iota
	catchFishRelease
	catchFishLost
	catchFishEat
	catchFishRetain
	catchFishCount
)

const (
	buyFishingRod           = "Приобретите удочку, прежде чем рыбачить."
	catchFishSellMessage    = "🎣 Вы поймали рыбу `%v` и продали ее за %s"
	catchFishReleaseMessage = "🎣 Вы поймали рыбу `%v`, но решили отпустили ее\\."
	catchFishLostMessage    = "🎣 Вы не смогли выудить рыбу из воды\\."
	catchFishEatMessage     = "🎣 Вы поймали рыбу `%v` и съели ее\\."
	catchFishRetainMessage  = "🎣 Вы поймали рыбу `%v` и оставили ее себе\\."
)

// !рыбалка
func (a *App) handleFishing(c tele.Context) error {
	user := getUser(c)
	if !user.Fisher {
		return userError(c, buyFishingRod)
	}
	ok := a.model.UpdateEnergy(user, -energyDelta, energyCap)
	if !ok {
		return userError(c, notEnoughEnergy)
	}
	fish := randomFish()

	switch catchFishType(rand.Intn(int(catchFishCount))) {
	case catchFishSell:
		return a.sellFish(c, user, fish)
	case catchFishRelease:
		return releaseFish(c, fish)
	case catchFishLost:
		return lostFish(c)
	case catchFishEat:
		return a.eatFish(c, user, fish)
	case catchFishRetain:
		return a.retainFish(c, user, fish)
	default:
		return internalError(c, errors.New("unknown fish type"))
	}
}

func (a *App) sellFish(c tele.Context, u model.User, fish string) error {
	reward := randInRange(fishSellMinPrice, fishSellMaxPrice)
	a.model.UpdateMoney(u, reward)
	return c.Send(fmt.Sprintf(catchFishSellMessage, fish, formatMoney(reward)), tele.ModeMarkdownV2)
}

func releaseFish(c tele.Context, fish string) error {
	return c.Send(fmt.Sprintf(catchFishReleaseMessage, fish), tele.ModeMarkdownV2)

}

func lostFish(c tele.Context) error {
	return c.Send(catchFishLostMessage, tele.ModeMarkdownV2)
}

func (a *App) eatFish(c tele.Context, u model.User, fish string) error {
	if hasFullEnergy(u) {
		return a.retainFish(c, u, fish)
	}
	a.model.UpdateEnergy(u, energyDelta, energyCap)
	return c.Send(fmt.Sprintf(catchFishEatMessage, fish), tele.ModeMarkdownV2)
}

func (a *App) retainFish(c tele.Context, u model.User, fish string) error {
	a.model.AddFish(u)
	return c.Send(fmt.Sprintf(catchFishRetainMessage, fish), tele.ModeMarkdownV2)
}

var fishes = []string{
	"Щука",
	"Окунь",
	"Судак",
	"Ерш",
	"Берш",
	"Жерех",
	"Голавль",
	"Змееголов",
	"Налим",
	"Угорь",
	"Сом",
	"Лосось",
	"Хариус",
	"Форель",
	"Голец",
	"Осетр",
	"Стерлядь",
	"Карп",
	"Карась",
	"Линь",
	"Лещ",
	"Язь",
	"Плотва",
	"Толстолобик",
	"Белоглазка",
	"Красноперка",
	"Уклейка",
	"Подуст",
	"Таймень",
}

func randomFish() string {
	return fishes[rand.Intn(len(fishes))]
}
