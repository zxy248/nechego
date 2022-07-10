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

func (a *App) formatRichTop(users []model.User) string {
	var top string
	for i := 0; i < len(users); i++ {
		top += fmt.Sprintf("_%d\\._ %s, %s\n",
			i+1, a.mustMentionUser(users[i]), formatMoney(users[i].Balance))
	}
	return top
}
