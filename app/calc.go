package app

import (
	"math/rand"

	"github.com/antonmedv/expr"
	tele "gopkg.in/telebot.v3"
)

const calculatorResponse = Response("%s %s <b>= %v</b>.")

var calculatorEmoji = []string{"🧠", "🧮", "🤔", "💭", "🤓"}

func handleCalculator(c tele.Context) error {
	msg := getMessage(c)
	result, err := expr.Eval(msg.RawArgument(), nil)
	if err != nil {
		return respond(c, Response("😵‍💫"))
	}
	emoji := calculatorEmoji[rand.Intn(len(calculatorEmoji))]
	return respond(c, calculatorResponse.Fill(emoji, msg.Argument(), result))
}
