package format

import (
	"fmt"
	"nechego/farm"
	"nechego/game"
	"nechego/item"
)

type Use struct {
	c *Connector
}

func NewUse() *Use {
	return &Use{NewConnector(" ")}
}

func (u *Use) String() string {
	return u.c.String()
}

func (u *Use) Callback(mention string) game.UseCallback {
	return game.UseCallback{
		Fertilizer: func(f *farm.Fertilizer) {
			u.c.Add(Fertilize(mention, f))
		},
	}
}

func Fertilize(mention string, f *farm.Fertilizer) string {
	return fmt.Sprintf("🛢 %s выливает <code>%v л.</code> удобрений на ферму.", mention, f.Volume)
}

func CannotUse(x *item.Item) string {
	return fmt.Sprintf("💡 Нельзя использовать %s.", Item(x))
}
