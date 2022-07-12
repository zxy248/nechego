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
	fishEaten     = "🐟 Вы съели рыбу\\."
	youAreFull    = "🐟 Вы не хотите есть."
)

func (a *App) handleEatFish(c tele.Context) error {
	user := getUser(c)
	if hasFullEnergy(user) {
		return c.Send(youAreFull)
	}
	ok := a.model.EatFish(user, eatFishEnergyDelta, energyTrueCap)
	if !ok {
		return c.Send(notEnoughFish)
	}
	out := appendEnergyRemaining(fishEaten, user.Energy+eatFishEnergyDelta)
	return c.Send(out, tele.ModeMarkdownV2)
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
	catchFishRelease catchFishType = iota
	catchFishLost
	catchFishSell
	catchFishEat
	catchFishRetain
)

const (
	buyFishingRod             = "Приобретите удочку, прежде чем рыбачить."
	catchFishSellMessage      = "🎣 Вы поймали рыбу %s и продали ее за %s"
	catchFishReleaseMessage   = "🎣 Вы поймали рыбу %s, но решили отпустить ее\\."
	catchFishLostMessage      = "🎣 Вы не смогли выудить рыбу из воды\\."
	catchFishEatMessage       = "🎣 Вы поймали рыбу %s и съели ее\\."
	catchFishRetainMessage    = "🎣 Вы поймали рыбу %s и оставили ее себе\\."
	catchFishSuccessThreshold = 0.5
	eatFishEnergyDelta        = 2
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

	switch randomFishType(user) {
	case catchFishSell:
		return a.sellFish(c, user)
	case catchFishRelease:
		return releaseFish(c, user)
	case catchFishLost:
		return lostFish(c, user)
	case catchFishEat:
		return a.eatFish(c, user)
	case catchFishRetain:
		return a.retainFish(c, user)
	default:
		return internalError(c, errors.New("unknown fish type"))
	}
}

func randomFishType(u model.User) catchFishType {
	success := []catchFishType{catchFishSell, catchFishEat, catchFishRetain}
	failure := []catchFishType{catchFishRelease, catchFishLost}
	r := rand.Float64()
	switch luckModifier(u) {
	case terribleLuckModifier:
		r -= .20
	case badLuckModifier:
		r -= .10
	case goodLuckModifier:
		r += .05
	case excellentLuckModifier:
		r += .10
	}
	if r >= catchFishSuccessThreshold {
		return success[rand.Intn(len(success))]
	}
	return failure[rand.Intn(len(failure))]
}

func (a *App) sellFish(c tele.Context, u model.User) error {
	fish := randomFish()
	a.model.UpdateMoney(u, fish.price())
	out := fmt.Sprintf(catchFishSellMessage, fish, formatMoney(fish.price()))
	out = appendEnergyRemaining(out, u.Energy-energyDelta)
	return c.Send(out, tele.ModeMarkdownV2)
}

func releaseFish(c tele.Context, u model.User) error {
	fish := randomFish()
	out := fmt.Sprintf(catchFishReleaseMessage, fish)
	out = appendEnergyRemaining(out, u.Energy-energyDelta)
	return c.Send(out, tele.ModeMarkdownV2)

}

func lostFish(c tele.Context, u model.User) error {
	out := appendEnergyRemaining(catchFishLostMessage, u.Energy-energyDelta)
	return c.Send(out, tele.ModeMarkdownV2)
}

func (a *App) eatFish(c tele.Context, u model.User) error {
	if hasFullEnergy(u) {
		return a.retainFish(c, u)
	}
	fish := randomFish()
	a.model.UpdateEnergy(u, eatFishEnergyDelta, energyTrueCap)
	out := fmt.Sprintf(catchFishEatMessage, fish)
	out = appendEnergyRemaining(out, u.Energy-energyDelta+eatFishEnergyDelta)
	return c.Send(out, tele.ModeMarkdownV2)
}

func (a *App) retainFish(c tele.Context, u model.User) error {
	fish := randomFish()
	a.model.AddFish(u)
	out := fmt.Sprintf(catchFishRetainMessage, fish)
	out = appendEnergyRemaining(out, u.Energy-energyDelta)
	return c.Send(out, tele.ModeMarkdownV2)
}

var fishNames = []string{
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

const (
	fishPricePerKg = 10
	minFishWeight  = 100
	maxFishWeight  = 5000
)

type fish struct {
	weight int // in grams
	name   string
}

func (f *fish) String() string {
	weight := float64(f.weight) / 1000
	return fmt.Sprintf("`%s (%.2f кг)`", f.name, weight)
}

func (f *fish) price() int {
	return int(float64(f.weight) / 1000 * fishPricePerKg)
}

func randomFish() *fish {
	return &fish{
		weight: randInRange(minFishWeight, maxFishWeight),
		name:   fishNames[rand.Intn(len(fishNames))],
	}
}
