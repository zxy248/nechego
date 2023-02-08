package format

import (
	"fmt"
	"math/rand"
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
	NothingSold          = "💵 Ничего не продано."
	NothingBought        = "💵 Ничего не куплено."
	NotOnline            = "🚫 Этот пользователь не в сети."
	CannotBan            = "😖 Этого пользователя нельзя забанить."
	CannotFight          = "🛡 С этим пользователем нельзя подраться."
	FightVersusPvE       = "🛡 Оппонент находится в PvE режиме."
	FightFromPvE         = "🛡 Вы находитесь в PvE режиме."
)

func Item(i *item.Item) string {
	return fmt.Sprintf("<code>%s</code>", i.Value)
}

func ItemsComma(items []*item.Item) string {
	c := NewConnector(", ")
	for _, x := range items {
		c.Add(Item(x))
	}
	return c.String()
}

func NumItem(n int, i *item.Item) string {
	return NumString(n, Item(i))
}

func NumString(n int, s string) string {
	return fmt.Sprintf("<code>%d ≡ </code> %s", n, s)
}

func Items(items []*item.Item) string {
	if len(items) == 0 {
		return Empty
	}
	c := NewConnector("\n")
	for i, v := range items {
		c.Add(NumItem(i, v))
	}
	return c.String()
}

func Catch(items []*item.Item) string {
	if len(items) == 0 {
		return Empty
	}
	c := NewConnector("\n")
	price, weight := 0.0, 0.0
	for i, v := range items {
		if f, ok := v.Value.(*fishing.Fish); ok {
			price += f.Price()
			weight += f.Weight
			c.Add(NumItem(i, v))
		}
	}
	c.Add(fmt.Sprintf("Стоимость: %s\nВес: %s",
		Money(int(price)), Weight(weight)))
	return c.String()
}

func Products(products []*game.Product) string {
	if len(products) == 0 {
		return Empty
	}
	c := NewConnector("\n")
	for i, p := range products {
		c.Add(fmt.Sprintf("%s, %s", NumItem(i, p.Item), Money(p.Price)))
	}
	return c.String()
}

func Money(q int) string {
	return fmt.Sprintf("<code>%d %s</code>", q, money.Currency)
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

func Eaten(mention string, eaten ...*item.Item) string {
	if len(eaten) == 0 {
		return NoFood
	}
	emoji, verb := "🥤", "выпил(а)"
	c := NewConnector(", ")
	for _, x := range eaten {
		if f, ok := x.Value.(*food.Food); !ok || !f.Beverage() {
			emoji, verb = "🍊", "съел(а)"
		}
		c.Add(Item(x))
	}
	return fmt.Sprintf("%s %s %s %s.", emoji, mention, verb, c.String())
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
		return fmt.Sprintf("<b>✉ %s: Новых сообщений нет.</b>", mention)
	}
	c := NewConnector("\n")
	c.Add(fmt.Sprintf("<b>✉ %s: Сообщения</b>", mention))
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
	return fmt.Sprintf("✉ %s совершает рассылку за %s.", mention, Money(price))
}

func UserBanned(hours int) string {
	suffix := "ов"
	switch hours {
	case 1:
		suffix = ""
	case 2, 3, 4:
		suffix = "а"

	}
	return fmt.Sprintf("🚫 Пользователь заблокирован на %d час%s.", hours, suffix)
}

func CannotDrop(i *item.Item) string {
	return fmt.Sprintf("♻ Нельзя выложить %s.", Item(i))
}

func Drop(mention string, i *item.Item) string {
	return fmt.Sprintf("♻ %s выкладывает %s.", mention, Item(i))
}

func CannotPick(i *item.Item) string {
	return fmt.Sprintf("♻ Нельзя взять %s.", Item(i))
}

func Pick(mention string, i *item.Item) string {
	return fmt.Sprintf("🫳 %s берёт %s.", mention, Item(i))
}

func NotOnFloor(key int) string {
	return fmt.Sprintf("🗄 Предмета %s нет на полу.", Key(key))
}

func CannotSell(i *item.Item) string {
	return fmt.Sprintf("🏪 Нельзя продать %s.", Item(i))
}

func Sold(mention string, profit int, sold ...*item.Item) string {
	if len(sold) == 0 {
		return NothingSold
	}
	c := NewConnector(", ")
	for _, x := range sold {
		c.Add(Item(x))
	}
	return fmt.Sprintf("💵 %s продаёт %s и зарабатывает %s.",
		mention, c.String(), Money(profit))
}

func Bought(mention string, cost int, bought ...*item.Item) string {
	if len(bought) == 0 {
		return NothingBought
	}
	c := NewConnector(", ")
	for _, x := range bought {
		c.Add(Item(x))
	}
	return fmt.Sprintf("🛒 %s покупает %s за %s.",
		mention, c.String(), Money(cost))
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
	return fmt.Sprintf("🎣 %s получает %s.", mention, Item(i))
}

func DrawNet(n *fishing.Net) string {
	count := n.Count()
	caught := "Поймано"
	if count == 1 {
		caught = "Поймана"
	}
	c := NewConnector("\n")
	c.Add("<b>🕸 Сеть вытянута.</b>")
	c.Add("<i>🐟 %s <code>%s</code>.</i>")
	return fmt.Sprintf(c.String(), caught, fish(count))
}

func Net(n *fishing.Net) string {
	c := NewConnector("\n")
	c.Add("<b>🕸 У вас есть рыболовная сеть на <code>%d</code> слотов.</b>")
	c.Add("<i>🐟 Команды: <code>!закинуть</code>, <code>!вытянуть</code>.</i>")
	return fmt.Sprintf(c.String(), n.Capacity)
}

func fish(count int) string {
	suffix := ""
	switch count {
	case 1:
		suffix = "а"
	case 2, 3, 4:
		suffix = "ы"
	}
	return fmt.Sprintf("%d рыб%s", count, suffix)
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

func mention(id int64, text string) string {
	return fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, id, text)
}

func PvPMode() string {
	return "⚔ PvP режим активирован."
}

func PvEMode() string {
	minutes := pvp.WaitForPvE / time.Minute
	return fmt.Sprintf("🛡 PvE режим активируется через %d минут.", minutes)
}

func Fight(mentionA, mentionB string, strengthA, strengthB float64) string {
	const fighter = "<b>%s</b> <code>[%.2f]</code>"
	const versus = "<b><i>vs.</i></b>"
	const fight = "⚔️ " + fighter + " " + versus + " " + fighter
	return fmt.Sprintf(fight, mentionA, strengthA, mentionB, strengthB)
}

func WinnerTook(mention string, i *item.Item) string {
	return fmt.Sprintf("🥊 %s забирает %s у проигравшего.", mention, Item(i))
}

func AttackerDrop(mention string, i *item.Item) string {
	return fmt.Sprintf("🌀 %s уронил %s во время драки.", mention, Item(i))
}

func Win(mention string, elo float64) string {
	return fmt.Sprintf("🏆 %s <code>(+%.1f)</code> выигрывает в поединке.", mention, elo)
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
		mention, CombatStatus(u.CombatMode.Status()),

		Energy(u.Energy), Balance(u.Balance().Total()),
		Strength(u.Strength(w)), Rating(u.Rating),
		Luck(u.Luck()), Messages(u.Messages),

		Modset(u.Modset(w)),
		Status(u.Status),
	)
}
