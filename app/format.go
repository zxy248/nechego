package app

import (
	"bytes"
	"fmt"
	"html"
	"math/rand"
	"nechego/input"
	"nechego/model"
	"nechego/numbers"
	"nechego/pets"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func photoFromBytes(data []byte) *tele.Photo {
	return &tele.Photo{File: tele.FromReader(bytes.NewReader(data))}
}

func formatError(s string) string {
	return "⭕️ " + s
}

func formatWarning(s string) string {
	return "⚠️ " + s
}

func formatMoney(n int) HTML {
	var s string
	switch p0 := n % 10; {
	case n >= 10 && n <= 20:
		s = fmt.Sprintf("%d рублей", n)
	case p0 == 1:
		s = fmt.Sprintf("%d рубль", n)
	case p0 >= 2 && p0 <= 4:
		s = fmt.Sprintf("%d рубля", n)
	default:
		s = fmt.Sprintf("%d рублей", n)
	}
	return HTML(fmt.Sprintf("<code>%s 🪙</code>", s))
}

func formatWeight(n float64) HTML {
	return HTML(fmt.Sprintf("<code>%.2f кг ⚖️</code>", n))
}

func formatEnergy(n int) HTML {
	return HTML(fmt.Sprintf("<code>%d ⚡️</code>", n))
}

func formatStrength(n float64) HTML {
	return HTML(fmt.Sprintf("<code>%.2f 💪</code>", n))
}

func formatMessages(n int) HTML {
	return HTML(fmt.Sprintf("<code>%d ✉️</code>", n))
}

func formatFood(n int) HTML {
	return HTML(fmt.Sprintf("<code>%d 🍊</code>", n))
}

func formatPercentage(v float64) HTML {
	return HTML(fmt.Sprintf("<code>%d%%</code>", int(v*100)))
}

func formatStatus(s ...string) HTML {
	var out string
	for _, t := range s {
		out += "<i>" + t + "</i>\n"
	}
	return HTML(strings.TrimSpace(out))
}

func formatCommand(c input.Command) HTML {
	return HTML("<code>" + c.String() + "</code>")
}

func formatTitles(s ...string) string {
	if len(s) > 0 {
		s[0] = strings.Title(s[0])
	}
	titles := joinWords(s...)
	if titles == "" {
		titles = "Пользователь"
	}
	return titles
}

func formatIcons(icon ...string) HTML {
	icons := strings.Join(icon, "·")
	return HTML("<code>" + icons + "</code>")
}

func itemize(s ...string) HTML {
	var out string
	for _, t := range s {
		out += "<b>•</b> " + t + "\n"
	}
	return ellipsizeEmpty(strings.TrimSpace(out))
}

func enumerate(s ...string) HTML {
	var out string
	for i, t := range s {
		out += fmt.Sprintf("<i>%d.</i> %s\n", i+1, t)
	}
	return ellipsizeEmpty(strings.TrimSpace(out))
}

func energyRemaining(n int) HTML {
	return HTML(fmt.Sprintf("<i>Энергии осталось: %s</i>", formatEnergy(n)))
}

const maxTopNumber = 5

// clampTopNumber returns x if x < maxTopNumber; otherwise returns maxTopNumber.
func clampTopNumber(x int) int {
	return numbers.Min(x, maxTopNumber)
}

func (a *App) itemizeUsers(u ...model.User) HTML {
	mentions := []string{}
	for _, uu := range u {
		mentions = append(mentions, string(a.mustMentionUser(uu)))
	}
	return itemize(mentions...)
}

func (a *App) enumerateUsers(u ...model.User) HTML {
	mentions := []string{}
	for _, uu := range u {
		mentions = append(mentions, string(a.mustMentionUser(uu)))
	}
	return enumerate(mentions...)
}

func itemizeCommands(c ...input.Command) HTML {
	s := []string{}
	for _, cc := range c {
		s = append(s, string(formatCommand(cc)))
	}
	return itemize(s...)
}

func ellipsizeEmpty(s string) HTML {
	if s == "" {
		return HTML("<code>. . .</code>")
	}
	return HTML(s)
}

func joinSections(s ...string) string {
	return strings.Join(s, "\n\n")
}

func joinLines(s ...string) string {
	return strings.Join(s, "\n")
}

func joinWords(s ...string) string {
	return strings.Join(s, " ")
}

func mention(uid int64, name string) HTML {
	return HTML(fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, uid, html.EscapeString(name)))
}

var (
	emojisActive   = []string{"🔈", "🔔", "✅", "🆗", "▶️"}
	emojisInactive = []string{"🔇", "🔕", "💤", "❌", "⛔️", "🚫", "⏹"}
)

func activeEmoji() string {
	return emojisActive[rand.Intn(len(emojisActive))]
}

func inactiveEmoji() string {
	return emojisInactive[rand.Intn(len(emojisInactive))]
}

var meals = []string{"завтрак", "полдник", "обед", "ужин", "перекус"}

func randomMeal() string {
	return meals[rand.Intn(len(meals))]
}

func formatPet(p *pets.Pet) HTML {
	return HTML(fmt.Sprintf("<code>%s %s %s (%s)</code>",
		p.Species.Icon(), strings.Title(p.Species.String()), p.Name, p.Gender.Icon()))
}
