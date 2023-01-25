package format

import (
	"fmt"
	"html"
	"nechego/fishing"
	"nechego/food"
	"nechego/game"
	"nechego/item"
	"nechego/modifier"
	"nechego/money"
	"nechego/phone"
	"strings"
	"time"
)

const (
	Empty                = "<code>. . .</code>"
	NoMoney              = "💵 Недостаточно средств."
	NoEnergy             = "⚡ Недостаточно энергии."
	AdminsOnly           = "⚠️ Эта команда доступна только администрации."
	RepostMessage        = "✉️ Перешлите сообщение пользователя."
	UserBanned           = "🚫 Пользователь заблокирован."
	UserUnbanned         = "✅ Пользователь разблокирован."
	CannotAttackYourself = "🛡️ Вы не можете напасть на самого себя."
	NoFood               = "🍊 У вас закончилась подходящая еда."
	NotHungry            = "🍊 Вы не хотите есть."
	InventoryFull        = "🗄 Ваш инвентарь заполнен."
	BadMarketName        = "🏪 Такое название не подходит для магазина."
	MarketRenamed        = "🏪 Вы назвали магазин."
	SpecifyMoney         = "💵 Укажите количество средств."
	BadMoney             = "💵 Некорректное количество средств."
	CannotCraft          = "🛠 Эти предметы нельзя объединить."
	InventorySorted      = "🗃 Инвентарь отсортирован."
	NoPhone              = "📱 У вас нет телефона."
	BadPhone             = "☎ Некорректный формат номера."
)

func Mention(uid int64, name string) string {
	return fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, uid, html.EscapeString(name))
}

func Item(i *item.Item) string {
	return fmt.Sprintf("<code>%s</code>", i.Value)
}

func ItemsComma(items []*item.Item) string {
	r := make([]string, 0, len(items))
	for _, x := range items {
		r = append(r, Item(x))
	}
	return strings.Join(r, ", ")
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
	lines := make([]string, 0, len(items))
	for i, v := range items {
		lines = append(lines, NumItem(i, v))
	}
	return strings.Join(lines, "\n")
}

func Catch(items []*item.Item) string {
	lines := []string{}
	price, weight := 0.0, 0.0
	for i, v := range items {
		if f, ok := v.Value.(*fishing.Fish); ok {
			price += f.Price()
			weight += f.Weight
			lines = append(lines, NumItem(i, v))
		}
	}
	if len(lines) == 0 {
		return Empty
	}
	tail := fmt.Sprintf("Стоимость: %s\nВес: %s",
		Money(int(price)), Weight(weight))
	lines = append(lines, tail)
	return strings.Join(lines, "\n")
}

func Products(products []*game.Product) string {
	if len(products) == 0 {
		return Empty
	}
	lines := make([]string, 0, len(products))
	for i, p := range products {
		lines = append(lines, NumItem(i, p.Item)+", "+Money(p.Price))
	}
	return strings.Join(lines, "\n")
}

func Money(q int) string {
	return fmt.Sprintf("<code>%d %s</code>", q, money.Symbol)
}

func Balance(q int) string {
	return "💵 " + Money(q)
}

func Weight(w float64) string {
	return fmt.Sprintf("<code>%.2f кг ⚖️</code>", w)
}

func Energy(e int) string {
	return fmt.Sprintf("<code>⚡ %d</code>", e)
}

func EnergyOutOf(e, max int) string {
	return fmt.Sprintf("<code>%d из %d ⚡</code>", e, max)
}

func EnergyRemaining(e int) string {
	return fmt.Sprintf("<i>Энергии осталось: %s</i>", Energy(e))
}

func Eat(i *item.Item) string {
	emoji, verb := "🍊", "съели"
	if x, ok := i.Value.(*food.Food); ok && x.Beverage() {
		emoji, verb = "🥤", "выпили"
	}
	return fmt.Sprintf("%s Вы %s %s.", emoji, verb, Item(i))
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

func ModifierEmojis(m []*modifier.Mod) string {
	emojis := []string{}
	for _, x := range m {
		if x.Emoji != "" {
			emojis = append(emojis, x.Emoji)
		}
	}
	return fmt.Sprintf("<code>%s</code>", strings.Join(emojis, "·"))
}

func ModifierDescriptions(m []*modifier.Mod) string {
	descs := []string{}
	for _, x := range m {
		descs = append(descs, x.Description)
	}
	return fmt.Sprintf("<i>%s</i>", strings.Join(descs, "\n"))
}

func ModifierTitles(m []*modifier.Mod) string {
	titles := []string{}
	for _, x := range m {
		if x.Title != "" {
			titles = append(titles, x.Title)
		}
	}
	if len(titles) == 0 {
		titles = append(titles, "пользователь")
	}
	titles[0] = strings.Title(titles[0])
	return strings.Join(titles, " ")
}

func Percentage(p float64) string {
	return fmt.Sprintf("%.1f%%", p*100)
}

func SMSes(mention string, smses []*phone.SMS) string {
	if len(smses) == 0 {
		return fmt.Sprintf("<b>✉ %s: Новых сообщений нет.</b>", mention)
	}
	lines := make([]string, 0, len(smses))
	for _, sms := range smses {
		lines = append(lines, SMS(sms))
	}
	head := fmt.Sprintf("<b>✉ %s: Сообщения</b>\n", mention)
	end := strings.Join(lines, "\n")
	return head + end
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
	head := "<b>👥 Контакты</b>\n"
	lines := make([]string, 0, len(cc))
	for _, c := range cc {
		lines = append(lines, c.String())
	}
	tail := strings.Join(lines, "\n")
	return head + tail
}

func MessageSent(sender, receiver phone.Number) string {
	return fmt.Sprintf("📱 Сообщение отправлено.\n\n"+
		"✉ <code>%v</code> → <code>%v</code>", sender, receiver)
}
