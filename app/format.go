package app

import (
	"bytes"
	"fmt"
	"math/rand"
	"nechego/input"
	"nechego/model"
	"strings"

	tele "gopkg.in/telebot.v3"
)

// photoFromBytes converts the image data to Photo.
func photoFromBytes(data []byte) *tele.Photo {
	return &tele.Photo{File: tele.FromReader(bytes.NewReader(data))}
}

// markdownEscaper escapes any character with the code between 1 and 126
// inclusively with a preceding backslash.
var markdownEscaper = func() *strings.Replacer {
	var table []string
	for i := 1; i <= 126; i++ {
		c := string(rune(i))
		table = append(table, c, "\\"+c)
	}
	return strings.NewReplacer(table...)
}()

var errorSigns = []string{"❌", "🚫", "⭕️", "🛑", "⛔️", "📛", "💢", "❗️", "‼️", "⚠️"}

// errorSign returns a random error sign.
func errorSign() string {
	return errorSigns[rand.Intn(len(errorSigns))]
}

// makeError formats the error message.
func makeError(s string) string {
	return errorSign() + " " + s
}

func internalError(c tele.Context, err error) error {
	c.Send(makeError("Ошибка сервера"))
	return err
}

func userError(c tele.Context, msg string) error {
	return c.Send(makeError(msg))
}

func userErrorMarkdown(c tele.Context, msg string) error {
	return c.Send(makeError(msg), tele.ModeMarkdownV2)
}

// randInRange returns a random value in range [min, max].
func randInRange(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// formatMoney formats the specified amount of money.
func formatMoney(n int) string {
	var out string
	switch p0 := n % 10; {
	case n >= 10 && n <= 20:
		out = fmt.Sprintf("%v монет", n)
	case p0 == 1:
		out = fmt.Sprintf("%v монета", n)
	case p0 >= 2 && p0 <= 4:
		out = fmt.Sprintf("%v монеты", n)
	default:
		out = fmt.Sprintf("%v монет", n)
	}
	return fmt.Sprintf("`%s 💰`", out)
}

func formatEnergy(n int) string {
	return fmt.Sprintf("`%v ⚡️`", n)
}

func formatStrength(n float64) string {
	return fmt.Sprintf("`%.2f 💪`", n)
}

func formatMessages(n int) string {
	return fmt.Sprintf("`%d ✍️`", n)
}

func formatFishes(n int) string {
	return fmt.Sprintf("`%d 🎣`", n)
}

func formatRatio(v float64) string {
	return fmt.Sprintf("`%d%%`", int(v*100))
}

func (a *App) formatUnorderedList(users []model.User) string {
	var list string
	for _, u := range users {
		list += fmt.Sprintf("— %s\n", a.mustMentionUser(u))
	}
	if list == "" {
		list = "…\n"
	}
	return list
}

func (a *App) formatOrderedList(users []model.User) string {
	var list string
	for i, u := range users {
		list += fmt.Sprintf("_%d\\._ %s\n", i+1, a.mustMentionUser(u))
	}
	if list == "" {
		list = "…\n"
	}
	return list
}

func formatCommandList(commands []input.Command) string {
	var list string
	for _, c := range commands {
		list += fmt.Sprintf("— `%s`\n", markdownEscaper.Replace(input.CommandText(c)))
	}
	if list == "" {
		list = "…\n"
	}
	return list
}

func (a *App) formatTopStrength(users []model.User) (string, error) {
	var top string
	for i, u := range users {
		s, err := a.actualUserStrength(u)
		if err != nil {
			return "", err
		}
		top += fmt.Sprintf("_%d\\._ %s, %s\n",
			i+1, a.mustMentionUser(u), formatStrength(s))
	}
	return top, nil
}

func (a *App) formatTopRich(users []model.User) string {
	var top string
	for i := 0; i < len(users); i++ {
		top += fmt.Sprintf("_%d\\._ %s, %s\n",
			i+1, a.mustMentionUser(users[i]), formatMoney(users[i].Summary()))
	}
	return top
}

func energyRemaining(energy int) string {
	return fmt.Sprintf("_Энергии осталось: %s_", formatEnergy(energy))
}

func appendEnergyRemaining(s string, energy int) string {
	return s + "\n\n" + energyRemaining(energy)
}

// topNumber returns l if l < maxTopNumber; otherwise returns maxTopNumber.
func topNumber(l int) int {
	if l < maxTopNumber {
		return l
	}
	return maxTopNumber
}

func formatStatus(desc ...string) string {
	status := strings.Join(desc, "\n")
	if status != "" {
		status = fmt.Sprintf("_%s_", markdownEscaper.Replace(status))
	}
	return status
}

func formatTitles(title ...string) string {
	if len(title) > 0 {
		title[0] = strings.Title(title[0])
	}
	titles := strings.Join(title, " ")
	if titles == "" {
		titles = "Пользователь"
	}
	return titles
}

func formatIcons(icon ...string) string {
	icons := strings.Join(icon, "·")
	return "`" + icons + "`"
}

func itemize(s ...string) string {
	var out string
	for _, ss := range s {
		out += "· " + ss + "\n"
	}
	return out
}

func (a *App) itemizeUsers(u ...model.User) string {
	mentions := []string{}
	for _, uu := range u {
		m := a.mustMentionUser(uu)
		mentions = append(mentions, m)
	}
	return itemize(mentions...)
}
