package app

import (
	"errors"
	"fmt"
	"math/rand"

	tele "gopkg.in/telebot.v3"
)

const (
	boughtFishingRod = "🎣 Вы приобрели удочку за `%s 💰`"
	alreadyCanFish   = "Вы уже приобрели удочку"
	fishingRodCost   = 50
)

// !удочка
func (a *App) handleFishingRod(c tele.Context) error {
	user := getUser(c)
	if user.CanFish {
		return userError(c, alreadyCanFish)
	}
	ok := a.model.UpdateMoney(user, -fishingRodCost)
	if !ok {
		return userError(c, notEnoughMoney)
	}
	a.model.AllowFishing(user)
	return c.Send(fmt.Sprintf(boughtFishingRod, formatAmount(fishingRodCost)),
		tele.ModeMarkdownV2)
}

type catchFishType int

const (
	catchFishSell catchFishType = iota
	catchFishRelease
	catchFishBad
	catchFishEat
	catchFishCount
)

const (
	buyFishingRod           = "Приобретите удочку, прежде чем рыбачить."
	catchFishSellMessage    = "🎣 Вы поймали рыбу `%v` и продали ее за `%v 💰`"
	catchFishReleaseMessage = "🎣 Вы поймали рыбу `%v`, но решили отпустили ее\\."
	catchFishBadMessage     = "🎣 Вы не смогли выудить рыбу из воды\\."
	catchFishEatMessage     = "🎣 Вы поймали рыбу `%v` и съели ее\\."
	fishSellMinPrice        = 1
	fishSellMaxPrice        = 40
)

// !рыбалка
func (a *App) handleFishing(c tele.Context) error {
	user := getUser(c)
	if !user.CanFish {
		return userError(c, buyFishingRod)
	}
	ok := a.model.UpdateEnergy(user, -energyDelta, energyCap)
	if !ok {
		return userError(c, notEnoughEnergy)
	}
	fish := randomFish()
	reward := randInRange(fishSellMinPrice, fishSellMaxPrice)

	switch catchFishType(rand.Intn(int(catchFishCount))) {
	case catchFishSell:
		a.model.UpdateMoney(user, reward)
		return c.Send(fmt.Sprintf(catchFishSellMessage, fish, formatAmount(reward)), tele.ModeMarkdownV2)
	case catchFishRelease:
		return c.Send(fmt.Sprintf(catchFishReleaseMessage, fish), tele.ModeMarkdownV2)
	case catchFishBad:
		return c.Send(catchFishBadMessage, tele.ModeMarkdownV2)
	case catchFishEat:
		a.model.UpdateEnergy(user, energyDelta, energyCap)
		return c.Send(fmt.Sprintf(catchFishEatMessage, fish), tele.ModeMarkdownV2)
	default:
		return internalError(c, errors.New("unknown fish type"))
	}
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
