package handlers

import (
	"fmt"
	"math/rand"
	"nechego/teleutil"
	"regexp"

	"github.com/antonmedv/expr"
	tele "gopkg.in/telebot.v3"
)

type Calculator struct{}

var calculatorRe = regexp.MustCompile("!калькулятор (.*)")

func (h *Calculator) Match(s string) bool {
	return calculatorRe.MatchString(s)
}

func (h *Calculator) Handle(c tele.Context) error {
	arg := teleutil.Args(c, calculatorRe)[1]
	result, err := expr.Eval(arg, nil)
	if err != nil {
		return c.Send("😵‍💫")
	}
	emojis := [...]string{"🧠", "🧮", "🤔", "💭", "🤓"}
	emoji := emojis[rand.Intn(len(emojis))]
	out := fmt.Sprintf("%s %s <b>= %v</b>.", emoji, arg, result)
	return c.Send(out, tele.ModeHTML)
}
