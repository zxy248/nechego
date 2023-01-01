package app

import (
	"bytes"
	"fmt"
	"math/rand"
	"nechego/input"
	"nechego/model"
	"nechego/numbers"
	"nechego/pets"
	"strings"
	"time"

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

func formatMoney(n int) string {
	return fmt.Sprintf("<code>%d ₽</code>", n)
}

func formatDebtStatus(u model.User) string {
	if u.Debtor() {
		return "У вас нет кредитов."
	}
	return fmt.Sprintf("Вы должны банку %s", formatMoney(u.Debt))
}

func formatWeight(n float64) string {
	return fmt.Sprintf("<code>%.2f кг ⚖️</code>", n)
}

func formatEnergy(n int) string {
	return fmt.Sprintf("<code>%d ⚡️</code>", n)
}

func formatStrength(n float64) string {
	return fmt.Sprintf("<code>%.2f 💪</code>", n)
}

func formatElo(n float64) string {
	return fmt.Sprintf("<code>%.1f ⚜️</code>", n)
}

func formatEloDelta(n float64) string {
	sign := "+"
	if n < 0 {
		sign = "-"
	}
	return fmt.Sprintf("<code>%s%.1f</code>", sign, n)
}

func formatMessages(n int) string {
	return fmt.Sprintf("<code>%d ✉️</code>", n)
}

func formatFood(n int) string {
	return fmt.Sprintf("<code>%d 🍊</code>", n)
}

func formatPercentage(v float64) string {
	return fmt.Sprintf("<code>%d%%</code>", int(v*100))
}

func formatStatus(s ...string) string {
	var out string
	for _, t := range s {
		out += fmt.Sprintf("<i>%s</i>\n", t)
	}
	return strings.TrimSpace(out)
}

func formatCommand(c input.Command) string {
	return fmt.Sprintf("<code>%s</code>", c)
}

func formatTitles(s ...string) string {
	if len(s) > 0 {
		s[0] = strings.Title(s[0])
	}
	titles := joinSpace(s...)
	if titles == "" {
		titles = "Пользователь"
	}
	return titles
}

func formatIcons(icon ...string) string {
	icons := strings.Join(icon, "·")
	return fmt.Sprintf("<code>%s</code>", icons)
}

const longListThreshold = 10

func ellipsizeLong(ss []string) []string {
	if len(ss) > longListThreshold {
		ss = ss[:longListThreshold]
		ss[longListThreshold-1] = ellipsis
	}
	return ss
}

func itemize(s ...string) string {
	var out string
	for _, t := range s {
		out += fmt.Sprintf("<b>•</b> %s\n", t)
	}
	return strings.TrimSpace(ellipsizeEmpty(out))
}

func enumerate(s ...string) string {
	var out string
	for i, t := range s {
		out += fmt.Sprintf("<i>%d.</i> %s\n", i+1, t)
	}
	return strings.TrimSpace(ellipsizeEmpty(out))
}

const ellipsis = "<code>. . .</code>"

func ellipsizeEmpty(s string) string {
	if s == "" {
		return ellipsis
	}
	return s
}

func formatEnergyRemaining(n int) string {
	return fmt.Sprintf("<i>Энергии осталось: %s</i>", formatEnergy(n))
}

func formatEnergyCooldown(d time.Duration) string {
	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	return fmt.Sprintf("⏰ До восстановления энергии: <code>%d минут %d секунд</code>.",
		mins, secs)
}

const maxTopNumber = 5

// clampTopNumber returns x if x < maxTopNumber; otherwise returns maxTopNumber.
func clampTopNumber(x int) int {
	return numbers.Min(x, maxTopNumber)
}

func (a *App) itemizeUsers(u ...model.User) string {
	s := []string{}
	for _, uu := range u {
		s = append(s, a.mention(uu))
	}
	return itemize(s...)
}

func (a *App) enumerateUsers(u ...model.User) string {
	s := []string{}
	for _, uu := range u {
		s = append(s, a.mention(uu))
	}
	return enumerate(s...)
}

func itemizeCommands(c ...input.Command) string {
	s := []string{}
	for _, cc := range c {
		s = append(s, formatCommand(cc))
	}
	return itemize(s...)
}

func joinSections(s ...string) string {
	return strings.Join(s, "\n\n")
}

func joinLines(s ...string) string {
	return strings.Join(s, "\n")
}

func joinSpace(s ...string) string {
	return strings.Join(s, " ")
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

func formatPet(p *pets.Pet) string {
	return fmt.Sprintf("<code>%s %s %s (%s)</code>",
		p.Species.Emoji(),
		strings.Title(p.Species.String()),
		p.Name,
		p.Gender.Emoji(),
	)
}
