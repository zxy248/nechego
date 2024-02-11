package command

import (
	"strconv"
	"strings"

	tele "gopkg.in/zxy248/telebot.v3"
)

type MarkovConfig struct {
	Prob *float64
}

func (m *MarkovConfig) Match(c tele.Context) bool {
	return cfgPattern.MatchString(c.Text())
}

func (m *MarkovConfig) Handle(c tele.Context) error {
	probString := strings.TrimSpace(cfgPattern.FindStringSubmatch(c.Text())[1])
	probFloat, err := strconv.ParseFloat(probString, 64)
	if err != nil {
		return err
	}
	*m.Prob = probFloat / 100
	return nil
}
