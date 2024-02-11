package middleware

import (
	"math/rand"
	"strings"

	"github.com/mb-14/gomarkov"
	"github.com/zxy248/nechego/game"
	tu "github.com/zxy248/nechego/teleutil"
	tele "gopkg.in/zxy248/telebot.v3"
)

type MarkovChain struct {
	Universe *game.Universe
}

func (c *MarkovChain) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		go func() {
			w := tu.Lock(ctx, c.Universe)
			if rand.Float64() < w.MarkovProb {
				chain := gomarkov.NewChain(1)
				for _, exp := range w.Expressions {
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
				if len(strings.TrimRight(responseText, " ")) != len(ctx.Text()) {
					ctx.Send(responseText)
				}
			}
			w.Unlock()
		}()
		return next(ctx)
	}
}
