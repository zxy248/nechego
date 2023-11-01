package format

import (
	"fmt"
	"nechego/farm"
	"nechego/farm/plant"
	"nechego/game"
	"nechego/item"
)

const (
	MaxSizeFarm = "🏡 Вы достигли максимального размера фермы."
	BadFarmName = "🏡 Такое название не подходит для фермы."
)

func Farm(mention string, f *farm.Farm, upgradeCost int) string {
	c := NewConnector("\n")
	c.Add(farmHeader(mention, f))
	if until := f.Until(); until > 0 {
		c.Add(fmt.Sprintf("<i>🌾 До урожая осталось %s</i>", Duration(until)))
	}
	if free := f.Free(); free > 0 {
		c.Add(fmt.Sprintf("<i>🌱 Можно посадить ещё %d %s</i>.",
			free, declPlant(free)))
	}
	if pending := f.Pending(); pending > 0 {
		c.Add(fmt.Sprintf("<i>🧺 Можно собрать урожай.</i>"))
	}
	if upgradeCost > 0 {
		c.Add(fmt.Sprintf("<i>💰 Можно купить землю за %s.</i>",
			Money(upgradeCost)))
	}
	return c.String()
}

func farmHeader(mention string, f *farm.Farm) string {
	name := ""
	if f.Name != "" {
		name = " " + Title(f.Name)
	}
	return fmt.Sprintf("<b>🏡 %s: Ферма%s (%d × %d)</b>",
		Name(mention), name, f.Rows, f.Columns)
}

func declPlant(n int) string {
	suffix := "й"
	switch n {
	case 1:
		suffix = "е"
	case 2, 3, 4:
		suffix = "я"
	}
	return "растени" + suffix
}

func Plant(p *plant.Plant) string {
	return fmt.Sprintf("<code>%s</code>", p)
}

func CannotPlant(i *item.Item) string {
	return fmt.Sprintf("🌱 Нельзя посадить %s.", Item(i))
}

func Planted(mention string, p ...*plant.Plant) string {
	if len(p) == 0 {
		return "🌱 Ничего не посажено."
	}
	c := NewConnector(", ")
	for _, x := range p {
		c.Add(Plant(x))
	}
	return fmt.Sprintf("🌱 %s посадил(а) %s.", Name(mention), c.String())
}

func Harvested(mention string, p ...*plant.Plant) string {
	if len(p) == 0 {
		return "🧺 Ничего не собрано."
	}
	c := NewConnector(", ")
	for _, x := range p {
		c.Add(Plant(x))
	}
	return fmt.Sprintf("🧺 %s собрал(а) %s.", Name(mention), c.String())
}

func FarmUpgraded(mention string, f *farm.Farm, cost int) string {
	c := NewConnector("\n")
	c.Add(fmt.Sprintf("💸 %s приобрел(а) землю за %s.", Name(mention), Money(cost)))
	c.Add(fmt.Sprintf("🏡 Новый размер фермы: <b>%d × %d</b>.", f.Rows, f.Columns))
	return c.String()
}

func FarmNamed(mention string, f *farm.Farm) string {
	return fmt.Sprintf("🏡 %s называет ферму %s.", Name(mention), Title(f.Name))
}

func PriceList(p *game.PriceList) string {
	out := fmt.Sprintf("<b>📊 Цены на %s</b>\n", p.Updated.Format("2006.01.02"))
	var table string
	for i, t := range plant.Types {
		table += fmt.Sprintf("<code>%s %-20s</code>", t, Money(p.Plants[t]))
		if i%2 == 0 {
			table += "<code>    </code>"
		} else {
			table += "\n"
		}
	}
	return out + table
}
