package format

import (
	"fmt"
	"math/rand"
	"nechego/fishing"
	"nechego/food"
	"nechego/game"
	"nechego/item"
	"nechego/money"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/message"
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
	InventoryOverflow    = "🗄 Инвентарь переполнен."
	BadMarketName        = "🏪 Такое название не подходит для магазина."
	MarketRenamed        = "🏪 Магазин переименован."
	SpecifyMoney         = "💵 Укажите количество средств."
	BadMoney             = "💵 Некорректное количество средств."
	CannotCraft          = "🛠 Эти предметы нельзя объединить."
	InventorySorted      = "🗃 Инвентарь отсортирован."
	BuyFishingRod        = "🎣 Приобретите удочку в магазине, прежде чем рыбачить."
	FishingRodBroke      = "🎣 Удочка сломалась."
	NoFishingRecords     = "🏆 Рекордов пока нет."
	NotOnline            = "🚫 Этот пользователь не в сети."
	CannotGetJob         = "💼 Такую работу получить пока нельзя."
	CannotFriend         = "👤 С этим пользователем нельзя подружиться."
	NonFriendTransfer    = "📦 Вещи можно передавать только тем, кто с вами дружит."
	ItemNotFound         = "🔖 Предмет не найден."
)

func Link(id int64, text string) string {
	return fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, id, text)
}

func Item(i *item.Item) string {
	return fmt.Sprintf("<code>%s</code>", i.Value)
}

func Selector(key int, s string) string {
	return fmt.Sprintf("<code>%2d ≡ </code>%s", key, s)
}

func Items(is []*item.Item) string {
	const (
		limit        = 30
		amortization = 5
	)
	if len(is) == 0 {
		return Empty
	}
	c := NewConnector("\n")
	for k, i := range is {
		if k >= limit && len(is) > limit+amortization {
			c.Add(fmt.Sprintf("<i>...и ещё %d предметов.</i>", len(is)-k))
			break
		}
		c.Add(Selector(k, Item(i)))
	}
	return c.String()
}

func Catch(is []*item.Item) string {
	if len(is) == 0 {
		return Empty
	}
	c := NewConnector("\n")
	price, weight := 0.0, 0.0
	for k, i := range is {
		if f, ok := i.Value.(*fishing.Fish); ok {
			price += f.Price()
			weight += f.Weight
			c.Add(Selector(k, Item(i)))
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
		c.Add(fmt.Sprintf("%s <code>⟨%s⟩</code>", Selector(k, Item(p.Item)), Money(p.Price)))
	}
	return c.String()
}

func Money(q int) string {
	p := message.NewPrinter(message.MatchLanguage("ru"))
	return p.Sprintf("<code>%d %s</code>", q, money.Currency)
}

func Name(s string) string {
	return fmt.Sprintf("<b>%s</b>", s)
}

func User(u *game.User) string {
	return Name(Link(u.ID, u.Name))
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

func Eaten(who string, is []*item.Item) string {
	if len(is) == 0 {
		return NoFood
	}
	emoji, verb := "🥤", "выпил(а)"
	c := NewConnector(", ")
	for _, i := range is {
		if f, ok := i.Value.(*food.Food); !ok || !f.Beverage() {
			emoji, verb = "🍊", "съел(а)"
		}
		c.Add(Item(i))
	}
	return fmt.Sprintf("%s %s %s %s.", emoji, Name(who), verb, c.String())
}

func CannotEat(is ...*item.Item) string {
	c := NewConnector(", ")
	for _, i := range is {
		c.Add(Item(i))
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

func Mods(ms []*game.Mod) string {
	c := NewConnector("\n")
	for _, m := range ms {
		c.Add(fmt.Sprintf("<i>%s %s</i>", m.Emoji, m.Description))
	}
	return c.String()
}

func Percentage(p float64) string {
	return fmt.Sprintf("%.1f%%", p*100)
}

func UserBanned(hours int) string {
	return fmt.Sprintf("🚫 Пользователь заблокирован на %d %s.", hours, declHours(hours))
}

func CannotDrop(i *item.Item) string {
	return fmt.Sprintf("♻ Нельзя выложить %s.", Item(i))
}

func Dropped(who string, is []*item.Item) string {
	if len(is) == 0 {
		return "♻ Ничего не выложено."
	}
	c := NewConnector(", ")
	for _, i := range is {
		c.Add(Item(i))
	}
	return fmt.Sprintf("♻ %s выкладывает %s.", Name(who), c.String())
}

func CannotPick(i *item.Item) string {
	return fmt.Sprintf("♻ Нельзя взять %s.", Item(i))
}

func Picked(who string, is []*item.Item) string {
	if len(is) == 0 {
		return "🫳 Ничего не взято."
	}
	c := NewConnector(", ")
	for _, i := range is {
		c.Add(Item(i))
	}
	return fmt.Sprintf("🫳 %s берёт %s.", Name(who), c.String())
}

func Cashout(who string, n int) string {
	return fmt.Sprintf("💵 %s откладывает %s.", Name(who), Money(n))
}

func CannotSell(i *item.Item) string {
	return fmt.Sprintf("🏪 Нельзя продать %s.", Item(i))
}

func Sold(who string, profit int, is []*item.Item) string {
	if len(is) == 0 {
		return "💵 Ничего не продано."
	}
	c := NewConnector(", ")
	for _, i := range is {
		c.Add(Item(i))
	}
	return fmt.Sprintf("💵 %s продаёт %s и зарабатывает %s.",
		Name(who), c.String(), Money(profit))
}

func Bought(who string, cost int, is []*item.Item) string {
	if len(is) == 0 {
		return "💵 Ничего не куплено."
	}
	c := NewConnector(", ")
	for _, i := range is {
		c.Add(Item(i))
	}
	return fmt.Sprintf("🛒 %s покупает %s за %s.",
		Name(who), c.String(), Money(cost))
}

func Crafted(who string, is ...*item.Item) string {
	if len(is) == 0 {
		return "🛠 Ничего не сделано."
	}
	c := NewConnector(", ")
	for _, i := range is {
		c.Add(Item(i))
	}
	return fmt.Sprintf("🛠 %s получает %s.", Name(who), c.String())
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

func FishCatch(who string, i *item.Item) string {
	return fmt.Sprintf("🎣 %s получает %s.", Name(who), Item(i))
}

func RecordCatch(p fishing.Parameter, e *fishing.Entry) string {
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
		n := fmt.Sprintf("<b><i>%s</i></b>. ", Link(e.ID, strconv.Itoa(1+i)))
		c.Add(n + Fish(e.Fish) + ", " + Money(int(e.Fish.Price())))
	}
	c.Add("")
	c.Add("<b>⚖ Самая тяжёлая рыба:</b>")
	c.Add(fmt.Sprintf("<b><i>%s</i></b> %s", Link(weight.ID, "→"), Fish(weight.Fish)))
	c.Add("")
	c.Add("<b>📐 Самая большая рыба:</b>")
	c.Add(fmt.Sprintf("<b><i>%s</i></b> %s", Link(length.ID, "→"), Fish(length.Fish)))
	return c.String()
}

func Fight(u1, u2 *game.User) string {
	const fighter = "%s <code>[%.2f]</code>"
	const versus = "<b><i>vs.</i></b>"
	const fight = "⚔️ " + fighter + " " + versus + " " + fighter
	return fmt.Sprintf(fight,
		Name(Link(u1.ID, u1.Name)),
		u1.Strength(),
		Name(Link(u2.ID, u2.Name)),
		u2.Strength())
}

func WinnerTook(who string, i *item.Item) string {
	return fmt.Sprintf("🥊 %s забирает %s у проигравшего.", Name(who), Item(i))
}

func AttackerDrop(who string, i *item.Item) string {
	return fmt.Sprintf("🌀 %s уронил %s во время драки.", Name(who), Item(i))
}

func Win(who string, elo float64) string {
	return fmt.Sprintf("🏆 %s <code>(+%.1f)</code> выигрывает в поединке.", Name(who), elo)
}

func Profile(u *game.User) string {
	head := fmt.Sprintf("<b>📇 %s: Профиль</b>", Name(Link(u.ID, u.Name)))
	entries := []string{
		Energy(u.Energy),
		Reputation{u.Reputation.Score(), u.ReputationFactor}.lhsEmoji(),
		Luck(u.Luck()),
		Strength(u.Strength()),
		Rating(u.Rating),
		Messages(u.Messages),
		Balance(u.Balance().Total()),
	}
	table := profileTable(entries)
	mods := Mods(u.Mods())
	status := Status(u.Status)
	return fmt.Sprintf("%s\n%s\n\n%s\n\n%s", head, table, mods, status)
}

func profileTable(entries []string) string {
	lines := []string{}
	for i, e := range entries {
		if i%2 == 0 {
			x := fmt.Sprintf("%-21s", e)
			lines = append(lines, x)
		} else {
			x := fmt.Sprintf(" %s", e)
			lines[len(lines)-1] += x
		}
	}
	for i, line := range lines {
		lines[i] = "<code>" + line + "</code>"
	}
	return strings.Join(lines, "\n")
}

func FundsCollected(who string, fs []*game.Fund) string {
	if len(fs) == 0 {
		return "🧾 Средств пока нет."
	}
	c := NewConnector("\n")
	c.Add(fmt.Sprintf("<b>🧾 %s получает средства:</b>", Name(who)))
	for i, f := range fs {
		if rest := len(fs) - i; i >= 15 && rest >= 5 {
			c.Add(fmt.Sprintf("<i>...и ещё <code>%d</code> шт.</i>", rest))
			break
		}
		c.Add(fmt.Sprintf("<code> • </code>%s <i>%s</i>", Item(f.Item), f.Source))
	}
	return c.String()
}

func GetJob(who string, hours int) string {
	const format = "💼 %s получает работу на <code>%d %s</code>."
	return fmt.Sprintf(format, Name(who), hours, declHours(hours))
}

func MarketShift(who string, s game.Shift) string {
	const clock = "<code>%02d:%02d</code>"
	const format = "🪪 С " + clock + " по " + clock + " вас обслуживает %s."
	return fmt.Sprintf(format,
		s.From.Hour(), s.From.Minute(),
		s.To.Hour(), s.To.Minute(),
		Name(who))
}

func Market(who string, m *game.Market) string {
	c := NewConnector("\n")
	c.Add(fmt.Sprintf("<b>%v</b>", m))
	if who != "" {
		c.Add(MarketShift(who, m.Shift))
	}
	c.Add(Products(m.Products()))
	return c.String()
}

func CannotSplit(i *item.Item) string {
	return fmt.Sprintf("🗃 Нельзя разделить %s.", Item(i))
}

func Splitted(who string, i *item.Item) string {
	return fmt.Sprintf("🗃 %s откладывает %s.", Name(who), Item(i))
}

func Index(i int) string {
	return fmt.Sprintf("<b><i>%d.</i></b>", 1+i)
}

func FriendRemoved(who1, who2 string) string {
	return fmt.Sprintf("😰 %s теперь не дружит с %s.", Name(who1), Name(who2))
}

func FriendAdded(who1, who2 string) string {
	return fmt.Sprintf("😊 %s теперь дружит с %s.", Name(who1), Name(who2))
}

func MutualFriends(who1, who2 string) string {
	return fmt.Sprintf("🤝 %s и %s теперь друзья.", Name(who1), Name(who2))
}

type Friend struct {
	Who    string
	Mutual bool
}

func FriendList(who string, friends []Friend) string {
	mutual := 0
	c := NewConnector("\n")
	for _, f := range friends {
		e := "💔"
		if f.Mutual {
			mutual++
			e = "❤️"
		}
		c.Add(e + " " + Name(f.Who))
	}
	header := fmt.Sprintf("<b>👥 %s: Друзья <code>[%d/%d]</code></b>",
		Name(who), mutual, len(friends))
	return header + "\n" + c.String()
}

func CannotTransfer(i *item.Item) string {
	return fmt.Sprintf("📦 Нельзя передать %s.", Item(i))
}

func Transfered(sender, receiver string, is ...*item.Item) string {
	if len(is) == 0 {
		return "📦 Ничего не передано."
	}
	c := NewConnector(", ")
	for _, i := range is {
		c.Add(Item(i))
	}
	const help = "<i>Используйте команду <code>!получить</code>, чтобы взять предметы.</i>"
	message := fmt.Sprintf("📦 %s передаёт %s %s.", Name(sender), Name(receiver), c.String())
	return message + "\n\n" + help
}

func Duration(d time.Duration) string {
	c := NewConnector(" ")
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		c.Add(fmt.Sprintf("%d ч.", h))
	}
	if m > 0 {
		c.Add(fmt.Sprintf("%d мин.", m))
	}
	if s > 0 {
		c.Add(fmt.Sprintf("%d сек.", s))
	}
	return c.String()
}

func Title(s string) string {
	return fmt.Sprintf("<b>«%s»</b>", s)
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
