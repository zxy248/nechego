package economy

import (
	"fmt"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type Capital struct {
	Universe *game.Universe
}

var capitalRe = handlers.Regexp("^!капитал")

func (h *Capital) Match(c tele.Context) bool {
	return capitalRe.MatchString(c.Text())
}

func (h *Capital) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	total, avg := world.Capital()
	magnate := world.TopUser(game.ByWealth)
	balance := magnate.Balance().Total()
	list := []string{
		fmt.Sprintf("<b>💸 Капитал беседы «%s»</b>: %s\n",
			c.Chat().Title, format.Money(total)),
		fmt.Sprintf("⚖️ В среднем на счету: %s\n",
			format.Money(avg)),
		fmt.Sprintf("🎩 В руках магната %s: %s,",
			format.User(magnate), format.Money(balance)),
		fmt.Sprintf("или <code>%s</code> от общего количества средств.",
			format.Percentage(float64(balance)/float64(total))),
	}
	return c.Send(strings.Join(list, "\n"), tele.ModeHTML)
}
