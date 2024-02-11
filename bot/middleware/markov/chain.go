package markov

import (
	"math/rand"
	"strings"

	"github.com/mb-14/gomarkov"
	"github.com/zxy248/nechego/game"
	tele "gopkg.in/zxy248/telebot.v3"
)

type Chain struct {
	Universe *game.Universe
	Prob     *float64
}

func (c *Chain) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		go func() {
			if rand.Float64() < *c.Prob {
				chain := gomarkov.NewChain(1)
				for _, exp := range c.Universe.Expressions {
					chain.Add(strings.Split(exp, " "))
				}
				var (
					err       error
					generated string
				)

				responseText := ctx.Text()

				for err == nil {
					tmp := strings.Split(responseText, " ")
					generated, err = chain.Generate(tmp[len(tmp)-1:])

					responseText += " " + strings.TrimSpace(strings.Trim(generated, "$"))

				}
				if len(responseText) != len(ctx.Text()) {
					ctx.Send(responseText)
				}
			}
		}()
		return next(ctx)
	}
}
