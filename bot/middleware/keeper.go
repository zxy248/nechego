package middleware

import (
	"slices"
	"strings"

	"github.com/zxy248/nechego/game"
	tu "github.com/zxy248/nechego/teleutil"
	tele "gopkg.in/zxy248/telebot.v3"
)

type MarkovKeeper struct {
	Universe *game.Universe
}

func (mk *MarkovKeeper) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		w := tu.Lock(ctx, mk.Universe)
		if !slices.Contains(w.Expressions, ctx.Text()) && len(strings.Split(ctx.Text(), " ")) > 1 {
			w.Expressions = append(w.Expressions, ctx.Text())
		}
		w.Unlock()
		return next(ctx)
	}
}
