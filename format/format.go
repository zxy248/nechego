package format

import (
	"fmt"
	"math/rand"
	"nechego/auction"
	"nechego/farm"
	"nechego/farm/plant"
	"nechego/fishing"
	"nechego/food"
	"nechego/game"
	"nechego/game/pvp"
	"nechego/item"
	"nechego/modifier"
	"nechego/money"
	"nechego/phone"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	Empty                = "<code>. . .</code>"
	NoMoney              = "💵 Недостаточно средств."
	NoEnergy             = "⚡ Недостаточно энергии."
	AdminsOnly           = "⚠️ Эта команда доступна только администрации."
	RepostMessage        = "✉️ Перешлите сообщение пользователя."
	UserUnbanned         = "✅ Пользователь разблокирован."
	CannotAttackYourself = "🛡️ Нельзя напасть на самого себя."
	NoFood               = "🍊 Подходящей еды нет."
	NotHungry            = "🍊 Вы не хотите есть."
	InventoryFull        = "🗄 Инвентарь переполнен."
	BadMarketName        = "🏪 Такое название не подходит для магазина."
	MarketRenamed        = "🏪 Магазин переименован."
	SpecifyMoney         = "💵 Укажите количество средств."
	BadMoney             = "💵 Некорректное количество средств."
	CannotCraft          = "🛠 Эти предметы нельзя объединить."
	InventorySorted      = "🗃 Инвентарь отсортирован."
	NoPhone              = "📱 У вас нет телефона."
	BadPhone             = "☎ Некорректный формат номера."
	BuyFishingRod        = "🎣 Приобретите удочку в магазине, прежде чем рыбачить."
	FishingRodBroke      = "🎣 Удочка сломалась."
	NoNet                = "🕸 У вас нет рыболовной сети."
	NetAlreadyCast       = "🕸 Рыболовная сеть уже закинута."
	CastNet              = "🕸 Рыболовная сеть закинута."
	NetNotCasted         = "🕸 Рыболовная сеть ещё не закинута."
	NoFishingRecords     = "🏆 Рекордов пока нет."
	NotOnline            = "🚫 Этот пользователь не в сети."
	CannotBan            = "😖 Этого пользователя нельзя забанить."
	CannotFight          = "🛡 С этим пользователем нельзя подраться."
	FightVersusPvE       = "🛡 Оппонент находится в <b>PvE-режиме</b>."
	FightFromPvE         = "🛡 Вы находитесь в <b>PvE-режиме</b>."
	CannotGetJob         = "💼 Такую работу получить пока нельзя."
	CannotFireJob        = "💼 Вы нигде не работаете."
	MaxSizeFarm          = "🏡 Вы достигли максимального размера фермы."
	NoLot                = "🏦 Лот уже продан."
	AuctionSell          = "🏦 Лот выставлен на продажу."
	AuctionFull          = "🏦 На аукционе нет места."
)

func Item(i *item.Item) string {
	return fmt.Sprintf("<code>%s</code>", i.Value)
}

func Selector(key int, s string) string {
	return fmt.Sprintf("<code>%d ≡ </code> %s", key, s)
}

func Items(i []*item.Item) string {
	const (
		limit        = 30
		amortization = 5
	)
	if len(i) == 0 {
		return Empty
	}
	c := NewConnector("\n")
	for k, x := range i {
		if k >= limit && len(i) > limit+amortization {
			c.Add(fmt.Sprintf("<i>...и ещё %d предметов.</i>", len(i)-k))
			break
		}
		c.Add(Selector(k, Item(x)))
	}
	return c.String()
}

func Catch(items []*item.Item) string {
	if len(items) == 0 {
		return Empty
	}
	c := NewConnector("\n")
	price, weight := 0.0, 0.0
	for k, x := range items {
		if f, ok := x.Value.(*fishing.Fish); ok {
			price += f.Price()
			weight += f.Weight
			c.Add(Selector(k, Item(x)))
		}
	}
	c.Add(fmt.Sprintf("Стоимость: %s", Money(int(price))))
	c.Add(fmt.Sprintf("Вес: %s", Weight(weight)))
	return c.String()
}

func Products(products []*game.Product) string {
	if len(products) == 0 {
		return Empty
	}
	c := NewConnector("\n")
	for k, p := range products {
		c.Add(fmt.Sprintf("%s, %s", Selector(k, Item(p.Item)), Money(p.Price)))
	}
	return c.String()
}

func Money(q int) string {
	return fmt.Sprintf("<code>%d %s</code>", q, money.Currency)
}

func Name(s string) string {
	return fmt.Sprintf("<b>%s</b>", s)
}

func Balance(q int) string {
	return "💵 " + Money(q)
}

func Weight(w float64) string {
	return fmt.Sprintf("<code>%.2f кг ⚖️</code>", w)
}

func Energy(e game.Energy) string {
	return fmt.Sprintf("<code>⚡ %.1f%%</code>", 100*e)
}

func EnergyRemaining(e game.Energy) string {
	return fmt.Sprintf("<i>Энергии осталось: %s</i>", Energy(e))
}

func Eaten(mention string, i ...*item.Item) string {
	if len(i) == 0 {
		return NoFood
	}
	emoji, verb := "🥤", "выпил(а)"
	c := NewConnector(", ")
	for _, x := range i {
		if f, ok := x.Value.(*food.Food); !ok || !f.Beverage() {
			emoji, verb = "🍊", "съел(а)"
		}
		c.Add(Item(x))
	}
	return fmt.Sprintf("%s %s %s %s.", emoji, Name(mention), verb, c.String())
}

func CannotEat(i ...*item.Item) string {
	c := NewConnector(", ")
	for _, x := range i {
		c.Add(Item(x))
	}
	return fmt.Sprintf("🤮 Нельзя съесть %s.", c.String())
}

func Fish(f *fishing.Fish) string {
	return fmt.Sprintf("<code>%s</code>", f)
}

func Plant(p *plant.Plant) string {
	return fmt.Sprintf("<code>%s</code>", p)
}

func Rating(r float64) string {
	return fmt.Sprintf("<code>⚜️ %.1f</code>", r)
}

func Strength(s float64) string {
	return fmt.Sprintf("<code>💪 %.1f</code>", s)
}

func Luck(l float64) string {
	return fmt.Sprintf("<code>🍀 %.1f</code>", 10*l)
}

func Messages(n int) string {
	return fmt.Sprintf("<code>✉️ %d</code>", n)
}

func Status(s string) string {
	return fmt.Sprintf("<i>%s</i>", s)
}

func Key(k int) string {
	return fmt.Sprintf("<code>#%d</code>", k)
}

func BadKey(k int) string {
	return fmt.Sprintf("🔖 Предмет %s не найден.", Key(k))
}

func Modset(s modifier.Set) string {
	c := NewConnector("\n")
	for _, x := range s.List() {
		c.Add(fmt.Sprintf("<i>%s %s</i>", x.Emoji, x.Description))
	}
	return c.String()
}

func Percentage(p float64) string {
	return fmt.Sprintf("%.1f%%", p*100)
}

func SMSes(mention string, smses []*phone.SMS) string {
	if len(smses) == 0 {
		return fmt.Sprintf("<b>✉ %s: Новых сообщений нет.</b>", Name(mention))
	}
	c := NewConnector("\n")
	c.Add(fmt.Sprintf("<b>✉ %s: Сообщения</b>", Name(mention)))
	for _, sms := range smses {
		c.Add(SMS(sms))
	}
	return c.String()
}

func SMS(sms *phone.SMS) string {
	format := "2006/02/01"
	if sms.Time.YearDay() == time.Now().YearDay() {
		format = "15:04"
	}
	return fmt.Sprintf("<code>|%s|</code> <code>%s</code><b>:</b> %s",
		sms.Time.Format(format), sms.Sender, sms.Text)
}

func SMSMaxLen(l int) string {
	return fmt.Sprintf("✉ Максимальная длина сообщения %d символов.", l)
}

type Contact struct {
	Name   string
	Number phone.Number
}

func (c Contact) String() string {
	return fmt.Sprintf("<b>→ <code>%s</code>:</b> %s", c.Number, c.Name)
}

func Contacts(cc []Contact) string {
	if len(cc) == 0 {
		return "👥 Контактов нет."
	}
	c := NewConnector("\n")
	c.Add("<b>👥 Контакты</b>")
	for _, contact := range cc {
		c.Add(contact.String())
	}
	return c.String()
}

func MessageSent(sender, receiver phone.Number) string {
	return fmt.Sprintf("📱 Сообщение отправлено.\n\n"+
		"✉ <code>%v</code> → <code>%v</code>", sender, receiver)
}

func SpamSent(mention string, price int) string {
	return fmt.Sprintf("✉ %s совершает рассылку за %s.", Name(mention), Money(price))
}

func UserBanned(hours int) string {
	return fmt.Sprintf("🚫 Пользователь заблокирован на %d %s.", hours, declHours(hours))
}

func CannotDrop(i *item.Item) string {
	return fmt.Sprintf("♻ Нельзя выложить %s.", Item(i))
}

func Dropped(mention string, i ...*item.Item) string {
	if len(i) == 0 {
		return "♻ Ничего не выложено."
	}
	c := NewConnector(", ")
	for _, x := range i {
		c.Add(Item(x))
	}
	return fmt.Sprintf("♻ %s выкладывает %s.", Name(mention), c.String())
}

func CannotPick(i *item.Item) string {
	return fmt.Sprintf("♻ Нельзя взять %s.", Item(i))
}

func Picked(mention string, i ...*item.Item) string {
	if len(i) == 0 {
		return "🫳 Ничего не взято."
	}
	c := NewConnector(", ")
	for _, x := range i {
		c.Add(Item(x))
	}
	return fmt.Sprintf("🫳 %s берёт %s.", Name(mention), c.String())
}

func CannotSell(i *item.Item) string {
	return fmt.Sprintf("🏪 Нельзя продать %s.", Item(i))
}

func Sold(mention string, profit int, i ...*item.Item) string {
	if len(i) == 0 {
		return "💵 Ничего не продано."
	}
	c := NewConnector(", ")
	for _, x := range i {
		c.Add(Item(x))
	}
	return fmt.Sprintf("💵 %s продаёт %s и зарабатывает %s.",
		Name(mention), c.String(), Money(profit))
}

func Bought(mention string, cost int, i ...*item.Item) string {
	if len(i) == 0 {
		return "💵 Ничего не куплено."
	}
	c := NewConnector(", ")
	for _, x := range i {
		c.Add(Item(x))
	}
	return fmt.Sprintf("🛒 %s покупает %s за %s.",
		Name(mention), c.String(), Money(cost))
}

func Crafted(mention string, i ...*item.Item) string {
	if len(i) == 0 {
		return "🛠 Ничего не сделано."
	}
	c := NewConnector(", ")
	for _, x := range i {
		c.Add(Item(x))
	}
	return fmt.Sprintf("🛠 %s получает %s.", Name(mention), c.String())
}

func BadFishOutcome() string {
	outcomes := [...]string{
		"Вы не смогли выудить рыбу.",
		"Рыба сорвалась с крючка.",
		"Рыба сорвала леску.",
		"Рыба скрылась в водорослях.",
		"Рыба выскользнула из рук.",
		"Вы отпустили рыбу обратно в воду.",
	}
	return "🎣 " + outcomes[rand.Intn(len(outcomes))]
}

func FishCatch(mention string, i *item.Item) string {
	return fmt.Sprintf("🎣 %s получает %s.", Name(mention), Item(i))
}

func DrawNet(n *fishing.Net) string {
	m := n.Count()
	c := NewConnector("\n")
	c.Add("<b>🕸 Сеть вытянута.</b>")
	c.Add("<i>🐟 %s <code>%d %s</code>.</i>")
	return fmt.Sprintf(c.String(), declCaught(m), m, declFish(m))
}

func Net(n *fishing.Net) string {
	c := NewConnector("\n")
	c.Add("<b>🕸 У вас есть рыболовная сеть на <code>%d</code> слотов.</b>")
	c.Add("<i>🐟 Команды: <code>!закинуть</code>, <code>!вытянуть</code>.</i>")
	return fmt.Sprintf(c.String(), n.Capacity)
}

func NewRecord(e *fishing.Entry, p fishing.Parameter) string {
	var p1, p2 string
	switch p {
	case fishing.Weight:
		p1, p2 = "весу", "тяжёлая"
	case fishing.Length:
		p1, p2 = "длине", "большая"
	case fishing.Price:
		p1, p2 = "цене", "дорогая"
	}
	c := NewConnector("\n")
	c.Add("<b>🎉 Установлен новый рекорд по %s рыбы!</b>")
	c.Add("%s это самая %s рыба из всех пойманных.")
	return fmt.Sprintf(c.String(), p1, Fish(e.Fish), p2)
}

func FishingRecords(price []*fishing.Entry, weight, length *fishing.Entry) string {
	c := NewConnector("\n")
	c.Add("<b>🏆 Книга рекордов 🎣</b>")
	c.Add("")
	c.Add("<b>💰 Самые дорогие рыбы:</b>")
	for i, e := range price {
		n := fmt.Sprintf("<b><i>%s</i></b>. ", mention(e.TUID, strconv.Itoa(1+i)))
		c.Add(n + Fish(e.Fish) + ", " + Money(int(e.Fish.Price())))
	}
	c.Add("")
	c.Add("<b>⚖ Самая тяжёлая рыба:</b>")
	c.Add(fmt.Sprintf("<b><i>%s</i></b> %s", mention(weight.TUID, "→"), Fish(weight.Fish)))
	c.Add("")
	c.Add("<b>📐 Самая большая рыба:</b>")
	c.Add(fmt.Sprintf("<b><i>%s</i></b> %s", mention(length.TUID, "→"), Fish(length.Fish)))
	return c.String()
}

func PvPMode() string {
	return "⚔ <b>PvP-режим</b> активирован."
}

func PvEMode() string {
	minutes := pvp.WaitForPvE / time.Minute
	return fmt.Sprintf("🛡 <b>PvE-режим</b> активируется через <code>%d минут</code>.", minutes)
}

func Fight(mentionA, mentionB string, strengthA, strengthB float64) string {
	const fighter = "%s <code>[%.2f]</code>"
	const versus = "<b><i>vs.</i></b>"
	const fight = "⚔️ " + fighter + " " + versus + " " + fighter
	return fmt.Sprintf(fight, Name(mentionA), strengthA, Name(mentionB), strengthB)
}

func WinnerTook(mention string, i *item.Item) string {
	return fmt.Sprintf("🥊 %s забирает %s у проигравшего.", Name(mention), Item(i))
}

func AttackerDrop(mention string, i *item.Item) string {
	return fmt.Sprintf("🌀 %s уронил %s во время драки.", Name(mention), Item(i))
}

func Win(mention string, elo float64) string {
	return fmt.Sprintf("🏆 %s <code>(+%.1f)</code> выигрывает в поединке.", Name(mention), elo)
}

func CombatStatus(s pvp.Status) string {
	return fmt.Sprintf("<code>[%v]</code>", s)
}

func Profile(mention string, u *game.User, w *game.World) string {
	const profile = `<b>📇 %s %s: Профиль</b>
<code>%-22s %s</code>
<code>%-22s %s</code>
<code>%-22s %s</code>

%s

%s`
	return fmt.Sprintf(
		profile,
		Name(mention), CombatStatus(u.CombatMode.Status()),

		Energy(u.Energy), Balance(u.Balance().Total()),
		Strength(u.Strength(w)), Rating(u.Rating),
		Luck(u.Luck()), Messages(u.Messages),

		Modset(u.Modset(w)),
		Status(u.Status),
	)
}

func FundsCollected(mention string, f ...*game.Fund) string {
	if len(f) == 0 {
		return "🧾 Средств пока нет."
	}
	c := NewConnector("\n")
	c.Add(fmt.Sprintf("<b>🧾 %s получает средства:</b>", Name(mention)))
	for i, x := range f {
		if rest := len(f) - i; i >= 15 && rest >= 5 {
			c.Add(fmt.Sprintf("<i>...и ещё <code>%d</code> пунктов.</i>", rest))
			break
		}
		c.Add(fmt.Sprintf("<b>·</b> %s <i>(%s)</i>", Item(x.Item), x.Source))
	}
	return c.String()
}

func GetJob(mention string, hours int) string {
	return fmt.Sprintf("💼 %s получает работу на <code>%d %s</code>.",
		Name(mention), hours, declHours(hours))
}

func MarketShift(mention string, s game.Shift) string {
	const clock = "<code>%02d:%02d</code>"
	const format = "🪪 С " + clock + " по " + clock + " вас обслуживает %s."
	return fmt.Sprintf(format,
		s.Start.Hour(), s.Start.Minute(),
		s.End.Hour(), s.End.Minute(),
		Name(mention))
}

func Market(mention string, m *game.Market) string {
	c := NewConnector("\n")
	c.Add(fmt.Sprintf("<b>%v</b>", m))
	if mention != "" {
		c.Add(MarketShift(mention, m.Shift))
	}
	c.Add(Products(m.Products()))
	return c.String()
}

func FireJob(mention string) string {
	return fmt.Sprintf("💼 %s покидает место работы.", Name(mention))
}

func Farm(mention string, f *farm.Farm, upgradeCost int) string {
	c := NewConnector("\n")
	c.Add(fmt.Sprintf("<b>🏡 %s: Ферма (%d × %d)</b>",
		Name(mention), f.Rows, f.Columns))
	if free := f.Free(); free > 0 {
		c.Add(fmt.Sprintf("<i>🌱 Можно посадить ещё %d %s</i>.",
			free, declPlant(free)))
	}
	if pending := f.Pending(); pending > 0 {
		c.Add(fmt.Sprintf("<i>🧺 Можно собрать урожай.</i>"))
	}
	if upgradeCost > 0 {
		c.Add(fmt.Sprintf("<i>💰 Можно купить землю за %s.</i>",
			Money(upgradeCost)))
	}
	return c.String()
}

func CannotPlant(i *item.Item) string {
	return fmt.Sprintf("🌱 Нельзя посадить %s.", Item(i))
}

func Planted(mention string, p ...*plant.Plant) string {
	if len(p) == 0 {
		return "🌱 Ничего не посажено."
	}
	c := NewConnector(", ")
	for _, x := range p {
		c.Add(Plant(x))
	}
	return fmt.Sprintf("🌱 %s посадил(а) %s.", Name(mention), c.String())
}

func Harvested(mention string, p ...*plant.Plant) string {
	if len(p) == 0 {
		return "🧺 Ничего не собрано."
	}
	c := NewConnector(", ")
	for _, x := range p {
		c.Add(Plant(x))
	}
	return fmt.Sprintf("🧺 %s собрал(а) %s.", Name(mention), c.String())
}

func PriceList(p *game.PriceList) string {
	out := fmt.Sprintf("<b>📊 Цены на %s</b>\n", p.Updated.Format("2006.01.02"))
	var table string
	for i, t := range plant.Types {
		table += fmt.Sprintf("<code>%s %20s</code>", t, Money(p.Price(t)))
		if i%2 == 0 {
			table += "<code>    </code>"
		} else {
			table += "\n"
		}
	}
	return out + table
}

func FarmUpgraded(mention string, f *farm.Farm, cost int) string {
	c := NewConnector("\n")
	c.Add(fmt.Sprintf("💸 %s приобрел(а) землю за %s.", Name(mention), Money(cost)))
	c.Add(fmt.Sprintf("🏡 Новый размер фермы: <b>%d × %d</b>.", f.Rows, f.Columns))
	return c.String()
}

func CannotSplit(i *item.Item) string {
	return fmt.Sprintf("🗃 Нельзя разделить %s.", Item(i))
}

func Splitted(mention string, i *item.Item) string {
	return fmt.Sprintf("🗃 %s откладывает %s.", Name(mention), Item(i))
}

func TopRating(mention func(*game.User) string, users ...*game.User) string {
	if len(users) == 0 {
		return fmt.Sprintf("🏆 Пользователей пока нет.")
	}
	c := NewConnector("\n")
	c.Add("<b>🏆 Боевой рейтинг</b>")
	for i, u := range users {
		c.Add(fmt.Sprintf("%s %s %s %s",
			Index(i),
			Name(mention(u)),
			u.CombatMode.Status().Emoji(),
			Rating(u.Rating)))
	}
	return c.String()
}

func Auction(lots []*auction.Lot, encode func(*auction.Lot) string) (string, *tele.ReplyMarkup) {
	s := "<b>🏦 Аукцион</b>"
	m := &tele.ReplyMarkup{}
	rows := []tele.Row{}
	for _, l := range lots {
		minutes := time.Until(l.Expire()) / time.Minute
		s := fmt.Sprintf("%s · %d %s · %d %s",
			l.Item.Value, l.Price(), money.Currency,
			minutes, declMinutes(int(minutes)))
		data := encode(l)
		rows = append(rows, m.Row(m.Data(s, data)))
	}
	m.Inline(rows...)
	return s, m
}

func AuctionBought(buyer, seller string, cost int, x *item.Item) string {
	return fmt.Sprintf("🤝 %s покупает %s у %s за %s.",
		Name(buyer), Item(x), Name(seller), Money(cost))
}

func Index(i int) string {
	return fmt.Sprintf("<b><i>%d.</i></b>", 1+i)
}

func declHours(n int) string {
	suffix := "ов"
	switch n {
	case 1:
		suffix = ""
	case 2, 3, 4:
		suffix = "а"
	}
	return "час" + suffix
}

func declMinutes(n int) string {
	suffix := ""
	switch n {
	case 1:
		suffix = "а"
	case 2, 3, 4:
		suffix = "ы"
	}
	return "минут" + suffix
}

func declFish(n int) string {
	suffix := ""
	switch n {
	case 1:
		suffix = "а"
	case 2, 3, 4:
		suffix = "ы"
	}
	return "рыб" + suffix
}

func declCaught(n int) string {
	if n == 1 {
		return "Поймана"
	}
	return "Поймано"
}

func declPlant(n int) string {
	suffix := "й"
	switch n {
	case 1:
		suffix = "е"
	case 2, 3, 4:
		suffix = "я"
	}
	return "растени" + suffix
}

func mention(id int64, text string) string {
	return fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, id, text)
}
