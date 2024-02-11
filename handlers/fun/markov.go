package fun

import (
	"strconv"
	"strings"

	"github.com/zxy248/nechego/game"
	"github.com/zxy248/nechego/handlers"
	tu "github.com/zxy248/nechego/teleutil"
	tele "gopkg.in/zxy248/telebot.v3"
)

var cfgPattern = handlers.NewRegexp("^!конфиг(.*)")

type MarkovConfig struct {
	Universe *game.Universe
}

func (m *MarkovConfig) Match(c tele.Context) bool {
	return cfgPattern.MatchString(c.Text())
}

func (m *MarkovConfig) Handle(c tele.Context) error {
	w := tu.Lock(c, m.Universe)
	probString := strings.TrimSpace(cfgPattern.FindStringSubmatch(c.Text())[1])
	probFloat, err := strconv.ParseFloat(probString, 64)
	if err != nil {
		return err
	}
	w.MarkovProb = probFloat / 100
	w.Unlock()
	return nil
}
