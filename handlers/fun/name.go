package fun

import (
	"fmt"
	"html"
	"strings"
	"unicode/utf8"

	"github.com/zxy248/nechego/handlers"
	tu "github.com/zxy248/nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Name struct{}

func (h *Name) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!имя")
}

func (h *Name) Handle(c tele.Context) error {
	user := c.Sender()
	if reply := c.Message().ReplyTo; reply != nil {
		user = reply.Sender
	}

	name := getNameArgument(c.Text())
	if name == "" {
		// Check the user's name.
		prefix := "Имя пользователя"
		if user.ID == c.Sender().ID {
			prefix = "Ваше имя"
		}
		link := tu.Link(c, user)
		out := fmt.Sprintf("%s: <b>%s</b> 🔖", prefix, link)
		return c.Send(out, tele.ModeHTML)
	}

	if n := 16; utf8.RuneCountInString(name) > n {
		out := fmt.Sprintf("⚠️ Максимальная длина имени %d символов.", n)
		return c.Send(out)
	}

	if user.ID == c.Bot().Me.ID {
		// Set the bot's name.
		payload := map[string]any{"name": name}
		if _, err := c.Bot().Raw("setMyName", payload); err != nil {
			return c.Send("⚠️ Подождите, прежде чем использовать эту команду снова.")
		}
		out := fmt.Sprintf("Теперь меня зовут <b>%s</b> ✅", name)
		return c.Send(out, tele.ModeHTML)
	}

	// User must be an admin to have a name.
	if err := tu.Promote(c, tu.Member(c, user)); err != nil {
		return err
	}

	if err := c.Bot().SetAdminTitle(c.Chat(), user, name); err != nil {
		return c.Send("⚠️ Не удалось установить имя.")
	}
	out := fmt.Sprintf("Имя <b>%s</b> установлено ✅", name)
	return c.Send(out, tele.ModeHTML)
}

func getNameArgument(s string) string {
	trim := strings.TrimPrefix(s, "!имя")
	if len(trim) == 0 || trim[0] != ' ' {
		return ""
	}
	return html.EscapeString(trim[1:])
}
