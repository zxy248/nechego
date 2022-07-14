package app

import (
	"fmt"
	"nechego/fishing"
	"nechego/model"
	"nechego/numbers"

	tele "gopkg.in/telebot.v3"
)

const (
	// fish
	notEnoughFish = "🐟 Недостаточно рыбы."
	fishEaten     = "🐟 Вы съели рыбу."
	youAreFull    = "🐟 Вы не хотите есть."

	// rod
	boughtFishingRod         = "🎣 Вы приобрели удочку за %s"
	alreadyCanFish           = "Вы уже приобрели удочку."
	notEnoughMoneyFishingRod = "Вам не хватает %s"
	buyFishingRod            = "Приобретите удочку, прежде чем рыбачить."

	// energy
	eatFishMinEnergy = 1
	eatFishMaxEnergy = 2
)

// !еда
func (a *App) handleEatFish(c tele.Context) error {
	user := getUser(c)
	if hasFullEnergy(user) {
		return respondPlain(c, youAreFull)
	}
	energyRestored, ok := a.eatFish(user)
	if !ok {
		return respondPlain(c, notEnoughFish)
	}
	return respondPlain(c, eatFishResponse(user, energyRestored))
}

func (a *App) eatFish(u model.User) (energyRestored int, enoughFish bool) {
	energyRestored = numbers.InRange(eatFishMinEnergy, eatFishMaxEnergy)
	enoughFish = a.model.EatFish(u, energyRestored, energyLimit)
	return
}

func eatFishResponse(u model.User, energyRestored int) string {
	return joinSections(fishEaten, energyRemaining(u.Energy+energyRestored))
}

// !удочка
func (a *App) handleFishingRod(c tele.Context) error {
	user := getUser(c)
	if user.Fisher {
		return userError(c, alreadyCanFish)
	}
	if ok := a.fishingRod(user); !ok {
		return userErrorMarkdown(c, fishingRodNotEnoughMoney(user))
	}
	return respondPlain(c, fishingRodSuccessResponse())
}

func (a *App) fishingRod(u model.User) bool {
	ok := a.model.UpdateMoney(u, -fishingRodPrice)
	if ok {
		a.model.AllowFishing(u)
	}
	return ok
}

func fishingRodNotEnoughMoney(u model.User) string {
	return fmt.Sprintf(notEnoughMoneyFishingRod, formatMoney(fishingRodPrice-u.Balance))
}

func fishingRodSuccessResponse() string {
	return fmt.Sprintf(boughtFishingRod, formatMoney(fishingRodPrice))
}

// !рыбалка
func (a *App) handleFishing(c tele.Context) error {
	user := getUser(c)
	if !user.Fisher {
		return userError(c, buyFishingRod)
	}
	session, ok := a.fishing(user)
	if !ok {
		return userError(c, notEnoughEnergy)
	}
	if session.Success() {
		respondPlain(c, a.catchFishResponse(user, session.Fish))
	}
	return respondPlain(c, session.Outcome.String())
}

func (a *App) fishing(u model.User) (fishing.Session, bool) {
	ok := a.model.UpdateEnergy(u, -energyDelta, energyLimit)
	if !ok {
		return fishing.Session{}, false
	}
	session := fishing.CastChance(fisherWinChance(u))
	if session.Outcome.Success() {
		a.collectFish(u, session.Fish)
	}
	return session, ok
}

func (a *App) collectFish(u model.User, f fishing.Fish) {
	if f.Weight < f.NormalWeight() {
		a.model.AddFish(u)
		return
	}
	a.model.InsertFish(model.MakeCatch(u, f))
}

func fisherWinChance(u model.User) float64 {
	r := fishing.SuccessChance
	switch luckModifier(u) {
	case terribleLuckModifier:
		r -= .12
	case badLuckModifier:
		r -= .06
	case goodLuckModifier:
		r += .04
	case excellentLuckModifier:
		r += .08
	}
	return r
}

const catchFish = "%s получает рыбу: %s"

func (a *App) catchFishResponse(u model.User, f fishing.Fish) string {
	return fmt.Sprintf(catchFish, a.mustMentionUser(u), f)
}

// !рыба
func (a *App) handleFish(c tele.Context) error {
	user := getUser(c)
	fishes, err := a.freshFishList(user)
	if err != nil {
		return internalError(c, err)
	}
	out := a.freshFishResponse(user, fishes)
	return respondPlain(c, out)
}

func (a *App) freshFishList(u model.User) ([]fishing.Fish, error) {
	catch, err := a.model.SelectFish(u)
	if err != nil {
		return nil, err
	}
	fishes := []fishing.Fish{}
	for _, c := range catch {
		if !c.Frozen {
			fishes = append(fishes, c.Fish)
		}
	}
	return fishes, nil
}

func (a *App) freshFishResponse(u model.User, f []fishing.Fish) string {
	return joinLines(fmt.Sprintf(freshFish, a.mustMentionUser(u)), formatFishList(f...))
}

func formatFishList(f ...fishing.Fish) string {
	lines := []string{}
	for _, ff := range f {
		lines = append(lines, fmt.Sprint(ff))
	}
	return joinLines(lines...)
}

// !продажа
func (a *App) handleSellFish(c tele.Context) error {
	user := getUser(c)
	price, err := a.sellFreshFish(user)
	if err != nil {
		return internalError(c, err)
	}
	out := sellFishResponse(price)
	return respondPlain(c, out)
}

func (a *App) sellFreshFish(u model.User) (int, error) {
	fishes, err := a.fishForSell(u)
	if err != nil {
		return 0, err
	}
	price := fishPrice(fishes...)
	a.model.UpdateMoney(u, price)
	return price, nil
}

func (a *App) fishForSell(u model.User) ([]fishing.Fish, error) {
	catch, err := a.model.SellFish(u)
	if err != nil {
		return nil, err
	}
	fishes := []fishing.Fish{}
	for _, c := range catch {
		fishes = append(fishes, c.Fish)
	}
	return fishes, nil
}

func fishPrice(f ...fishing.Fish) int {
	sum := 0
	for _, ff := range f {
		sum += ff.Price()
	}
	return sum
}

func sellFishResponse(price int) string {
	return fmt.Sprintf(soldFish, formatMoney(price))
}

const (
	freshFish    = "Улов %s"
	freezerFish  = "Холодильник %s"
	fishFrozen   = "Рыба заморожена."
	fishUnfrozen = "Рыба разморожена."
	soldFish     = "Рыбы продано на %s"
)

func (a *App) handleFreeze(c tele.Context) error {
	user := getUser(c)
	a.freezeFish(user)
	return respondPlain(c, fishFrozen)
}

func (a *App) freezeFish(u model.User) {
	a.model.FreezeFish(u)
}

func (a *App) handleUnfreeze(c tele.Context) error {
	user := getUser(c)
	a.unfreezeFish(user)
	return respondPlain(c, fishUnfrozen)
}

func (a *App) unfreezeFish(u model.User) {
	a.model.UnfreezeFish(u)
}

// !холодильник
func (a *App) handleFreezer(c tele.Context) error {
	user := getUser(c)
	fishes, err := a.frozenFishList(user)
	if err != nil {
		return internalError(c, err)
	}
	out := a.freezerFishResponse(user, fishes)
	return respondPlain(c, out)
}

func (a *App) frozenFishList(u model.User) ([]fishing.Fish, error) {
	catch, err := a.model.SelectFish(u)
	if err != nil {
		return nil, err
	}
	fishes := []fishing.Fish{}
	for _, c := range catch {
		if c.Frozen {
			fishes = append(fishes, c.Fish)
		}
	}
	return fishes, nil
}

func (a *App) freezerFishResponse(u model.User, f []fishing.Fish) string {
	return joinLines(fmt.Sprintf(freezerFish, a.mustMentionUser(u)), formatFishList(f...))
}
