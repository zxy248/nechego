package fun

import (
	"fmt"
	"math/rand/v2"

	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Infa struct{}

var infaRe = handlers.NewRegexp("^!(инфа|вероятность)(.*)")

func (h *Infa) Match(c tele.Context) bool {
	return infaRe.MatchString(c.Text())
}

func (h *Infa) Handle(c tele.Context) error {
	return c.Send(formatInfa(parseInfa(c.Text())))
}

func parseInfa(s string) string {
	return infaRe.FindStringSubmatch(s)[2]
}

func formatInfa(s string) string {
	templates := [...]string{
		"Здравый смысл говорит мне о том, что",
		"Благодаря чувственному опыту я определил, что",
		"Я думаю, что",
		"Используя диалектическую логику, я пришел к выводу, что",
		"Проведя некие изыскания, я выяснил, что",
		"Я провёл мысленный эксперимент и выяснил, что",
		"Мои интеллектуальные потуги привели меня к тому, что",
		"С помощью фактов и логики я доказал, что",
		"Как показывает практика,",
		"Прикинув раз на раз, я определился с тем, что",
		"Уверяю вас в том, что",
	}
	t := templates[rand.N(len(templates))]

	var p int
	if rand.N(100) == 0 {
		p = 100
	} else {
		d1 := rand.N(10)
		d2 := rand.N(10)
		p = d1*10 + d2
	}
	return fmt.Sprint(t, s, " с вероятностью ", p, "%")
}
