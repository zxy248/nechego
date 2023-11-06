package handlers

import (
	"errors"
	"fmt"
	"math/rand"
	"nechego/avatar"
	"nechego/fishing"
	"nechego/format"
	"nechego/game"
	"nechego/game/recipes"
	"nechego/handlers/parse"
	"nechego/item"
	"nechego/money"
	tu "nechego/teleutil"
	"nechego/valid"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	tele "gopkg.in/telebot.v3"
)

const inventoryCapacity = 20

type Inventory struct {
	Universe *game.Universe
}

var inventoryRe = Regexp("^!(инвентарь|лут)")

func (h *Inventory) Match(s string) bool {
	return inventoryRe.MatchString(s)
}

func (h *Inventory) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	items := user.Inventory.HkList()
	warn := ""
	if fullInventory(user.Inventory) {
		warn = " (!)"
	}
	head := fmt.Sprintf("<b>🗄 %s: Инвентарь <code>[%d/%d%s]</code></b>\n",
		tu.Link(c, user), len(items), inventoryCapacity, warn)
	list := format.Items(items)
	return c.Send(head+list, tele.ModeHTML)
}

func fullInventory(i *item.Set) bool {
	return i.Count() >= inventoryCapacity
}

type Sort struct {
	Universe *game.Universe
}

var sortRe = Regexp("^!сорт (.*)")

func (h *Sort) Match(s string) bool {
	return sortRe.MatchString(s)
}

func (h *Sort) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	items := []*item.Item{}
	seen := map[*item.Item]bool{}
	for _, k := range tu.NumArg(c, sortRe, 1) {
		x, ok := user.Inventory.ByKey(k)
		if !ok {
			return c.Send(format.BadKey(k), tele.ModeHTML)
		}
		if !seen[x] {
			items = append(items, x)
		}
		seen[x] = true
	}

	for _, x := range items {
		if !user.Inventory.Remove(x) {
			panic(fmt.Sprintf("sort: cannot remove %v", x))
		}
	}
	user.Inventory.AddFront(items...)
	return c.Send(format.InventorySorted)
}

type Catch struct {
	Universe *game.Universe
}

var catchRe = Regexp("^!улов")

func (h *Catch) Match(s string) bool {
	return catchRe.MatchString(s)
}

func (h *Catch) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if net, ok := user.FishingNet(); ok {
		caught := user.UnloadNet(net)
		for _, f := range caught {
			world.History.Add(user.TUID, f)
		}
	}
	head := fmt.Sprintf("<b>🐟 %s: Улов</b>\n", tu.Link(c, user))
	list := format.Catch(user.Inventory.HkList())
	return c.Send(head+list, tele.ModeHTML)
}

type Drop struct {
	Universe *game.Universe
}

var dropRe = Regexp("^!(выкинуть|выбросить|выложить|дроп|положить) (.*)")

func (h *Drop) Match(s string) bool {
	return dropRe.MatchString(s)
}

func (h *Drop) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	dropped := []*item.Item{}
	for _, key := range tu.NumArg(c, dropRe, 2) {
		item, ok := user.Inventory.ByKey(key)
		if !ok {
			c.Send(format.BadKey(key), tele.ModeHTML)
			break
		}
		if !user.Inventory.Move(world.Floor, item) {
			c.Send(format.CannotDrop(item), tele.ModeHTML)
			break
		}
		dropped = append(dropped, item)
	}
	world.Floor.Trim(10)
	return c.Send(format.Dropped(tu.Link(c, user), dropped...), tele.ModeHTML)
}

type Pick struct {
	Universe *game.Universe
}

var pickRe = Regexp("^!(взять|подобрать|поднять) (.*)")

func (h *Pick) Match(s string) bool {
	return pickRe.MatchString(s)
}

func (h *Pick) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if fullInventory(user.Inventory) {
		return c.Send(format.InventoryOverflow)
	}

	picked := []*item.Item{}
	for _, key := range tu.NumArg(c, pickRe, 2) {
		item, ok := world.Floor.ByKey(key)
		if !ok {
			c.Send(format.BadKey(key), tele.ModeHTML)
			break
		}
		if !world.Floor.Move(user.Inventory, item) {
			c.Send(format.CannotPick(item), tele.ModeHTML)
			break
		}
		picked = append(picked, item)
	}
	return c.Send(format.Picked(tu.Link(c, user), picked...), tele.ModeHTML)
}

type Floor struct {
	Universe *game.Universe
}

var floorRe = Regexp("^!(пол|мусор|вещи|предметы)")

func (h *Floor) Match(s string) bool {
	return floorRe.MatchString(s)
}

func (h *Floor) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	head := "<b>🗃️ Предметы</b>\n"
	list := format.Items(world.Floor.HkList())
	return c.Send(head+list, tele.ModeHTML)
}

type Market struct {
	Universe *game.Universe
}

var marketRe = Regexp("^!(магаз|шоп)")

func (h *Market) Match(s string) bool {
	return marketRe.MatchString(s)
}

func (h *Market) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	var who string
	if id, ok := world.Market.Shift.Worker(); ok {
		who = tu.Link(c, id)
	}
	return c.Send(format.Market(who, world.Market), tele.ModeHTML)
}

type PriceList struct {
	Universe *game.Universe
}

var priceListRe = Regexp("^!(прайс-?лист|цен)")

func (h *PriceList) Match(s string) bool {
	return priceListRe.MatchString(s)
}

func (h *PriceList) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	world.Market.PriceList.Refresh()
	return c.Send(format.PriceList(world.Market.PriceList), tele.ModeHTML)
}

type NameMarket struct {
	Universe *game.Universe
}

var nameMarketRe = Regexp("^!назвать магазин (.+)")

func (h *NameMarket) Match(s string) bool {
	return nameMarketRe.MatchString(s)
}

func (h *NameMarket) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	n := marketName(c.Text())
	if n == "" {
		return c.Send(format.BadMarketName)
	}
	world.Market.Name = n
	return c.Send(format.MarketRenamed)
}

func marketName(s string) string {
	n := nameMarketRe.FindStringSubmatch(s)[1]
	if !valid.Name(n) {
		return ""
	}
	return strings.Title(n)
}

type GetJob struct {
	Universe *game.Universe
}

var getJobRe = Regexp("^!(рохля|работа)")

func (h *GetJob) Match(s string) bool {
	return getJobRe.MatchString(s)
}

func (h *GetJob) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	const shiftHours = 2
	if time.Since(user.Retired) < 2*time.Hour || !world.Market.Shift.Begin(user.TUID, shiftHours*time.Hour) {
		return c.Send(format.CannotGetJob)
	}
	user.Retired = time.Now().Add(shiftHours * time.Hour)
	return c.Send(format.GetJob(tu.Link(c, user), shiftHours), tele.ModeHTML)
}

type QuitJob struct {
	Universe *game.Universe
}

var quitJobRe = Regexp("^!(уволиться|увольнение)")

func (h *QuitJob) Match(s string) bool {
	return quitJobRe.MatchString(s)
}

func (h *QuitJob) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if id, ok := world.Market.Shift.Worker(); ok && id == user.TUID {
		world.Market.Shift.Cancel()
		return c.Send(format.FireJob(tu.Link(c, id)), tele.ModeHTML)
	}
	return c.Send(format.CannotFireJob)
}

type Buy struct {
	Universe *game.Universe
}

var buyRe = Regexp("^!купить (.*)")

func (h *Buy) Match(s string) bool {
	return buyRe.MatchString(s)
}

func (h *Buy) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if fullInventory(user.Inventory) {
		return c.Send(format.InventoryOverflow)
	}

	bought := []*item.Item{}
	cost := 0
	for _, key := range tu.NumArg(c, buyRe, 1) {
		p, err := user.Buy(world, key)
		if errors.Is(err, game.ErrNoKey) {
			c.Send(format.BadKey(key), tele.ModeHTML)
			break
		} else if err != nil {
			c.Send(format.NoMoney, tele.ModeHTML)
			break
		}
		bought = append(bought, p.Item)
		cost += p.Price
	}
	return c.Send(format.Bought(tu.Link(c, user), cost, bought...), tele.ModeHTML)
}

type Eat struct {
	Universe *game.Universe
}

var eatRe = Regexp("^!(с[ъь]есть|еда) (.*)")

func (h *Eat) Match(s string) bool {
	return eatRe.MatchString(s)
}

func (h *Eat) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if user.Energy.Full() {
		return c.Send(format.NotHungry)
	}
	eaten := []*item.Item{}
	for _, key := range tu.NumArg(c, eatRe, 2) {
		item, ok := user.Inventory.ByKey(key)
		if !ok {
			c.Send(format.BadKey(key), tele.ModeHTML)
			break
		}
		if !user.Eat(item) {
			c.Send(format.CannotEat(item), tele.ModeHTML)
			break
		}
		eaten = append(eaten, item)
	}
	return c.Send(format.Eaten(tu.Link(c, user), eaten...)+"\n\n"+
		format.EnergyRemaining(user.Energy), tele.ModeHTML)

}

type EatQuick struct {
	Universe *game.Universe
}

var eatQuickRe = Regexp("^!еда")

func (h *EatQuick) Match(s string) bool {
	return eatQuickRe.MatchString(s)
}

func (h *EatQuick) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if user.Energy.Full() {
		return c.Send(format.NotHungry)
	}
	eaten := []*item.Item{}
	for !user.Energy.Full() {
		x, ok := user.EatQuick()
		if !ok {
			break
		}
		eaten = append(eaten, x)
	}
	return c.Send(format.Eaten(tu.Link(c, user), eaten...)+"\n\n"+
		format.EnergyRemaining(user.Energy), tele.ModeHTML)
}

type Fish struct {
	Universe *game.Universe
}

var fishRe = Regexp("^!(р[ыі]балка|ловля рыб)")

func (h *Fish) Match(s string) bool {
	return fishRe.MatchString(s)
}

func (h *Fish) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if fullInventory(user.Inventory) {
		return c.Send(format.InventoryOverflow)
	}

	rod, ok := user.FishingRod()
	if !ok {
		return c.Send(format.BuyFishingRod)
	}
	if !user.Energy.Spend(0.2) {
		return c.Send(format.NoEnergy)
	}
	item, caught := user.Fish(rod)
	if rod.Broken() {
		c.Send(format.FishingRodBroke)
	}
	if !caught {
		return c.Send(format.BadFishOutcome())
	}
	if f, ok := item.Value.(*fishing.Fish); ok {
		world.History.Add(user.TUID, f)
	}
	user.Inventory.Add(item)
	return c.Send(format.FishCatch(tu.Link(c, user), item), tele.ModeHTML)
}

type CastNet struct {
	Universe *game.Universe
}

var castNetRe = Regexp("^!закинуть")

func (h *CastNet) Match(s string) bool {
	return castNetRe.MatchString(s)
}

func (h *CastNet) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	err := user.CastNet()
	if errors.Is(err, game.ErrNoNet) {
		return c.Send(format.NoNet)
	} else if errors.Is(err, game.ErrNetAlreadyCast) {
		return c.Send(format.NetAlreadyCast)
	} else if err != nil {
		return err
	}
	return c.Send(format.CastNet)
}

type DrawNet struct {
	Universe *game.Universe
}

var drawNetRe = Regexp("^!вытянуть")

func (h *DrawNet) Match(s string) bool {
	return drawNetRe.MatchString(s)
}

func (h *DrawNet) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if fullInventory(user.Inventory) {
		return c.Send(format.InventoryOverflow)
	}

	net, ok := user.DrawNew()
	if !ok {
		return c.Send(format.NetNotCasted)
	}
	err := c.Send(format.DrawNet(net), tele.ModeHTML)
	caught := user.UnloadNet(net)
	for _, f := range caught {
		world.History.Add(user.TUID, f)
	}
	return err
}

type Net struct {
	Universe *game.Universe
}

var netRe = Regexp("^!сеть")

func (h *Net) Match(s string) bool {
	return netRe.MatchString(s)
}

func (h *Net) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	net, ok := user.FishingNet()
	if !ok {
		return c.Send(format.NoNet)
	}
	return c.Send(format.Net(net), tele.ModeHTML)
}

// RecordAnnouncer returns a fishing.RecordAnnouncer that sends a
// message to the chat specified by tgid.
func RecordAnnouncer(bot *tele.Bot, tgid tele.Recipient) fishing.RecordAnnouncer {
	return func(e *fishing.Entry, p fishing.Parameter) {
		bot.Send(tgid, format.NewRecord(e, p), tele.ModeHTML)
	}
}

type FishingRecords struct {
	Universe *game.Universe
}

var fishingRecordsRe = Regexp("^!рекорды")

func (h *FishingRecords) Match(s string) bool {
	return fishingRecordsRe.MatchString(s)
}

func (h *FishingRecords) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()
	byPrice := world.History.Top(fishing.Price, 10)
	byWeight := world.History.Top(fishing.Weight, 1)
	byLength := world.History.Top(fishing.Length, 1)
	for _, top := range [][]*fishing.Entry{byPrice, byWeight, byLength} {
		if len(top) == 0 {
			return c.Send(format.NoFishingRecords)
		}
	}
	return c.Send(format.FishingRecords(byPrice, byWeight[0], byLength[0]), tele.ModeHTML)
}

type Craft struct {
	Universe *game.Universe
}

var craftRe = Regexp("^!крафт (.*)")

func (h *Craft) Match(s string) bool {
	return craftRe.MatchString(s)
}

func (h *Craft) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	keys := tu.NumArg(c, craftRe, 1)
	recipe := []*item.Item{}
	for _, k := range keys {
		i, ok := user.Inventory.ByKey(k)
		if !ok {
			return c.Send(format.BadKey(k), tele.ModeHTML)
		}
		recipe = append(recipe, i)
	}
	crafted, ok := recipes.Craft(user.Inventory, recipe)
	if !ok {
		return c.Send(format.CannotCraft)
	}
	return c.Send(format.Crafted(tu.Link(c, user), crafted...), tele.ModeHTML)
}

type Status struct {
	Universe  *game.Universe
	MaxLength int
}

var statusRe = Regexp("^!статус (.*)")

func (h *Status) Match(s string) bool {
	return statusRe.MatchString(s)
}

func (h *Status) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if reply, ok := tu.Reply(c); ok {
		// If the user has admin rights, they can set a status
		// for other users.
		if !user.Admin() {
			return c.Send("💬 Нельзя установить статус другому пользователю.")
		}
		user = world.UserByID(reply.ID)
	}

	status := tu.Args(c, statusRe)[1]
	if utf8.RuneCountInString(status) > h.MaxLength {
		return c.Send(fmt.Sprintf("💬 Максимальная длина статуса %d символов.", h.MaxLength))
	}
	user.Status = status
	return c.Send("✅ Статус установлен.")
}

type Sell struct {
	Universe *game.Universe
}

func (h *Sell) Match(s string) bool {
	_, ok := sellCommand(s)
	return ok
}

func sellCommand(s string) (keys []int, ok bool) {
	return numCommand(parse.Match("!продать"), s)
}

func (h *Sell) Handle(c tele.Context) error {
	keys, _ := sellCommand(c.Text())
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	total := 0
	sold := []*item.Item{}
	for _, key := range keys {
		item, ok := user.Inventory.ByKey(key)
		if !ok {
			c.Send(format.BadKey(key), tele.ModeHTML)
			break
		}
		profit, ok := user.Sell(world, item)
		if !ok {
			c.Send(format.CannotSell(item), tele.ModeHTML)
			break
		}
		total += profit
		sold = append(sold, item)
	}
	return c.Send(format.Sold(tu.Link(c, user), total, sold...), tele.ModeHTML)
}

type SellQuick struct {
	Universe *game.Universe
}

var sellQuickRe = Regexp("^!продать")

func (h *SellQuick) Match(s string) bool {
	return sellQuickRe.MatchString(s)
}

func (h *SellQuick) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	total := 0
	sold := []*item.Item{}
	for _, item := range user.Inventory.List() {
		fish, ok := item.Value.(*fishing.Fish)
		if !ok || fish.Price() < 2000 {
			continue
		}
		profit, ok := user.Sell(world, item)
		if !ok {
			c.Send(format.CannotSell(item), tele.ModeHTML)
			break
		}
		total += profit
		sold = append(sold, item)
	}
	return c.Send(format.Sold(tu.Link(c, user), total, sold...), tele.ModeHTML)
}

type Stack struct {
	Universe *game.Universe
}

var stackRe = Regexp("^!сложить")

func (h *Stack) Match(s string) bool {
	return stackRe.MatchString(s)
}

func (h *Stack) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	user.Inventory.Stack()
	return c.Send("🗄 Вы сложили вещи.")
}

type Split struct {
	Universe *game.Universe
}

var splitRe = Regexp(`^!(отложить|разделить) (\d*) (\d*)`)

func (h *Split) Match(s string) bool {
	return splitRe.MatchString(s)
}

func (h *Split) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if fullInventory(user.Inventory) {
		return c.Send(format.InventoryOverflow)
	}

	args := tu.Args(c, splitRe)
	key, err := strconv.Atoi(args[2])
	if err != nil {
		return nil
	}
	count, err := strconv.Atoi(args[3])
	if err != nil {
		return nil
	}
	whole, ok := user.Inventory.ByKey(key)
	if !ok {
		return c.Send(format.BadKey(key), tele.ModeHTML)
	}
	part, ok := item.Split(whole, count)
	if !ok {
		return c.Send(format.CannotSplit(whole), tele.ModeHTML)
	}
	user.Inventory.Add(part)
	return c.Send(format.Splitted(tu.Link(c, user), part), tele.ModeHTML)
}

type Cashout struct {
	Universe *game.Universe
}

var cashoutRe = Regexp("^!(отложить|обнал|снять) (.*)")

func (h *Cashout) Match(s string) bool {
	return cashoutRe.MatchString(s)
}

func (h *Cashout) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()
	args := tu.NumArg(c, cashoutRe, 2)
	if len(args) != 1 {
		return c.Send(format.SpecifyMoney)
	}
	amount := args[0]
	if err := user.Balance().Cashout(amount); errors.Is(err, money.ErrBadMoney) {
		return c.Send(format.BadMoney)
	} else if errors.Is(err, money.ErrNoMoney) {
		return c.Send(format.NoMoney)
	} else if err != nil {
		return err
	}
	return c.Send(fmt.Sprintf("💵 Вы отложили %s.",
		format.Money(amount)), tele.ModeHTML)
}

type Fight struct {
	Universe *game.Universe
}

var fightRe = Regexp("^!(драка|дуэль|поединок|атака|битва|схватка|сражение|бой|борьба)")

func (h *Fight) Match(s string) bool {
	return fightRe.MatchString(s)
}

func (h *Fight) Handle(c tele.Context) error {
	// Sanity check before locking the world.
	reply, ok := tu.Reply(c)
	if !ok {
		return c.Send(format.RepostMessage)
	}
	if c.Sender().ID == reply.ID {
		return c.Send(format.CannotAttackYourself)
	}

	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()
	opnt := world.UserByID(reply.ID)

	// Can opponent fight back?
	if time.Since(opnt.LastMessage) > 10*time.Minute {
		return c.Send(format.NotOnline)
	}

	if !user.Energy.Spend(0.33) {
		return c.Send(format.NoEnergy)
	}

	// Fight begins.
	win, lose, elo := world.Fight(user, opnt)

	msg := format.NewConnector("\n\n")
	msg.Add(format.Fight(
		tu.Link(c, user.TUID),
		tu.Link(c, opnt.TUID),
		user.Strength(world),
		opnt.Strength(world)))

	// The winner takes a random item.
	if rand.Float64() < 0.02 {
		if x, ok := moveRandomItem(win.Inventory, lose.Inventory); ok {
			msg.Add(format.WinnerTook(tu.Link(c, win), x))
		}
	}
	// The attacker drops a random item.
	if rand.Float64() < 0.04 {
		if x, ok := moveRandomItem(world.Floor, user.Inventory); ok {
			msg.Add(format.AttackerDrop(tu.Link(c, user), x))
		}
	}
	msg.Add(format.Win(tu.Link(c, win), elo))
	return c.Send(msg.String(), tele.ModeHTML)
}

func moveRandomItem(dst, src *item.Set) (i *item.Item, ok bool) {
	i, ok = src.Random()
	if !ok {
		return nil, false
	}
	return i, src.Move(dst, i)
}

type Profile struct {
	Universe *game.Universe
	Avatars  *avatar.Storage
}

var profileRe = Regexp("^!(профиль|стат)")

func (h *Profile) Match(s string) bool {
	return profileRe.MatchString(s)
}

func (h *Profile) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if u, ok := tu.Reply(c); ok {
		user = world.UserByID(u.ID)
	}

	out := format.Profile(tu.Link(c, user), user, world)
	if a, ok := h.Avatars.Get(user.TUID); ok {
		a.Caption = out
		return c.Send(a, tele.ModeHTML)
	}
	return c.Send(out, tele.ModeHTML)
}

type TopStrong struct {
	Universe *game.Universe
}

var topStrongRe = Regexp("^!топ сил")

func (h *TopStrong) Match(s string) bool {
	return topStrongRe.MatchString(s)
}

func (h *TopStrong) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	users := world.SortedUsers(game.ByStrength)
	users = users[:min(len(users), 5)]
	list := []string{"🏋️‍♀️ <b>Самые сильные пользователи</b>"}
	for i, u := range users {
		list = append(list, fmt.Sprintf("<b><i>%d.</i></b> %s %s",
			i+1, tu.Link(c, u.TUID), format.Strength(u.Strength(world))))
	}
	return c.Send(strings.Join(list, "\n"), tele.ModeHTML)
}

type TopRating struct {
	Universe *game.Universe
}

var topRating = Regexp("^!(рейтинг|ммр|эло)")

func (h *TopRating) Match(s string) bool {
	return topRating.MatchString(s)
}

func (h *TopRating) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	who := func(u *game.User) string { return tu.Link(c, u.TUID) }
	users := world.SortedUsers(game.ByElo)
	const limit = 10
	if len(users) > limit {
		users = users[:limit]
	}
	return c.Send(format.TopRating(who, users...), tele.ModeHTML)
}

type TopRich struct {
	Universe *game.Universe
}

var topRich = Regexp("^!топ бога[тч]")

func (h *TopRich) Match(s string) bool {
	return topRich.MatchString(s)
}

func (h *TopRich) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	users := world.SortedUsers(game.ByWealth)
	users = users[:min(len(users), 5)]
	list := []string{"💵 <b>Самые богатые пользователи</b>"}
	for i, u := range users {
		list = append(list, fmt.Sprintf("<b><i>%d.</i></b> %s %s",
			i+1, tu.Link(c, u.TUID), format.Money(u.Balance().Total())))
	}
	return c.Send(strings.Join(list, "\n"), tele.ModeHTML)
}

type Capital struct {
	Universe *game.Universe
}

var capitalRe = Regexp("^!капитал")

func (h *Capital) Match(s string) bool {
	return capitalRe.MatchString(s)
}

func (h *Capital) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	total, avg := world.Capital()
	users := world.SortedUsers(game.ByWealth)
	users = users[:min(len(users), 5)]
	rich := users[0]
	balance := rich.Balance().Total()
	list := []string{
		fmt.Sprintf("💸 Капитал беседы <b>%s</b>: %s\n",
			c.Chat().Title, format.Money(total)),
		fmt.Sprintf("<i>В среднем на счету: %s</i>\n",
			format.Money(avg)),
		fmt.Sprintf("<i>В руках магната %s %s,</i>",
			tu.Link(c, users[0].TUID), format.Money(balance)),
		fmt.Sprintf("<i>или %s от общего количества средств.</i>\n",
			format.Percentage(float64(balance)/float64(total))),
	}
	return c.Send(strings.Join(list, "\n"), tele.ModeHTML)
}

type Balance struct {
	Universe *game.Universe
}

var balanceRe = Regexp("^!(баланс|деньги)")

func (h *Balance) Match(s string) bool {
	return balanceRe.MatchString(s)
}

func (h *Balance) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()
	return c.Send(fmt.Sprintf("💵 Ваш баланс: %s",
		format.Money(user.Balance().Total())), tele.ModeHTML)
}

type Funds struct {
	Universe *game.Universe
}

var fundsRe = Regexp("^!(зарплата|средства|получить)")

func (h *Funds) Match(s string) bool {
	return fundsRe.MatchString(s)
}

func (h *Funds) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if fullInventory(user.Inventory) {
		return c.Send(format.InventoryOverflow)
	}

	collected := user.Funds.Collect()
	for _, f := range collected {
		user.Inventory.Add(f.Item)
	}
	user.Inventory.Stack()
	return c.Send(format.FundsCollected(tu.Link(c, user), collected...), tele.ModeHTML)
}

type Energy struct {
	Universe *game.Universe
}

var energyRe = Regexp("^!энергия")

func (h *Energy) Match(s string) bool {
	return energyRe.MatchString(s)
}

func (h *Energy) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	emoji := "🔋"
	if user.Energy < 0.5 {
		emoji = "🪫"
	}
	return c.Send(fmt.Sprintf("%s Запас энергии: %s",
		emoji, format.Energy(user.Energy)), tele.ModeHTML)
}

type NamePet struct {
	Universe *game.Universe
}

var namePetRe = Regexp("^!назвать (.+)")

func (h *NamePet) Match(s string) bool {
	return namePetRe.MatchString(s)
}

func (h *NamePet) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	pet, ok := user.Pet()
	if !ok {
		return c.Send("🐱 У вас нет питомца.")
	}

	e := pet.Species.Emoji()
	n := petName(c.Text())
	if n == "" {
		return c.Send(fmt.Sprintf("%s Такое имя не подходит для питомца.", e))
	}
	pet.Name = n
	s := fmt.Sprintf("%s Вы назвали питомца <code>%s</code>.", e, n)
	return c.Send(s, tele.ModeHTML)
}

func petName(s string) string {
	n := namePetRe.FindStringSubmatch(s)[1]
	if !valid.Name(n) {
		return ""
	}
	return strings.Title(n)
}
