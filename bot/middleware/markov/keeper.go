package markov

import (
	"slices"

	"github.com/zxy248/nechego/game"
	tele "gopkg.in/zxy248/telebot.v3"
)

type MarkovKeeper struct {
	Universe *game.Universe
}

func (mk *MarkovKeeper) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if !slices.Contains(mk.Universe.Expressions, ctx.Text()) {
			mk.Universe.Expressions = append(mk.Universe.Expressions, ctx.Text())
		}
		return next(ctx)
	}
}
