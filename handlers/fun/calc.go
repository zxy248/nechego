package fun

import (
	"fmt"
	"github.com/zxy248/nechego/handlers"
	"math/rand"

	"github.com/antonmedv/expr"
	tele "gopkg.in/zxy248/telebot.v3"
)

type Calc struct{}

var calcRe = handlers.NewRegexp("^!(калькул|вычисл)[а-я]* (.+)")

func (h *Calc) Match(c tele.Context) bool {
	return calcRe.MatchString(c.Text())
}

func (h *Calc) Handle(c tele.Context) error {
	in := calcExpression(c.Text())
	out, err := expr.Eval(in, nil)
	if err != nil {
		return c.Send("😵‍💫")
	}

	es := [...]string{"🧠", "🧮", "🤔", "💭", "🤓"}
	e := es[rand.Intn(len(es))]
	s := fmt.Sprintf("%s %s <b>= %v</b>", e, in, out)
	return c.Send(s, tele.ModeHTML)
}

func calcExpression(s string) string {
	return calcRe.FindStringSubmatch(s)[2]
}
