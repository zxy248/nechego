package bot

import (
	"nechego/model"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Config struct {
	Token string
	DB    *model.DB
	Owner int64
}

type Bot struct {
	config   *Config
	bot      *tele.Bot
	keyboard *tele.ReplyMarkup

	admins    *model.Admins
	whitelist *model.Whitelist
	users     *model.Users
	pairs     *model.Pairs
	eblans    *model.Eblans
	bans      *model.Bans
	status    *model.Status
}

func NewBot(c *Config) (*Bot, error) {
	pref := tele.Settings{
		Token:  c.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	tb, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	kb := keyboard()

	return &Bot{
		config:   c,
		bot:      tb,
		keyboard: kb,

		admins:    &model.Admins{DB: c.DB},
		whitelist: &model.Whitelist{DB: c.DB},
		users:     &model.Users{DB: c.DB},
		pairs:     &model.Pairs{DB: c.DB},
		eblans:    &model.Eblans{DB: c.DB},
		bans:      &model.Bans{DB: c.DB},
		status:    &model.Status{DB: c.DB},
	}, nil
}

func (b *Bot) Start() {
	b.bot.Handle(tele.OnText, b.route, b.check)
	b.bot.Start()
}
