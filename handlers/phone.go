package handlers

import (
	"nechego/format"
	"nechego/game"
	"nechego/phone"
	"nechego/teleutil"
	"unicode/utf8"

	tele "gopkg.in/telebot.v3"
)

const smsMaxLen = 120

type ReceiveSMS struct {
	Universe *game.Universe
}

var receiveSMSRe = re("^!смс")

func (h *ReceiveSMS) Match(s string) bool {
	return receiveSMSRe.MatchString(s)
}

func (h *ReceiveSMS) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	p, ok := user.Phone()
	if !ok {
		return c.Send(format.NoPhone)
	}
	mention := teleutil.Mention(c, user)
	smses := world.SMS.Receive(p.Number)
	return c.Send(format.SMSes(mention, smses), tele.ModeHTML)
}

type SendSMS struct {
	Universe *game.Universe
}

var sendSMSRe = re("^!смс *(" + phone.NumberExpr + ") *(.+)")

func (h *SendSMS) Match(s string) bool {
	return sendSMSRe.MatchString(s)
}

func (h *SendSMS) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	p, ok := user.Phone()
	if !ok {
		return c.Send(format.NoPhone)
	}

	a := teleutil.Args(c, sendSMSRe)
	num, msg := a[1], a[2]
	if utf8.RuneCountInString(msg) > smsMaxLen {
		return c.Send(format.SMSMaxLen(smsMaxLen))
	}
	receiver, err := phone.MakeNumber(num)
	if err != nil {
		return c.Send(format.BadPhone)
	}
	world.SMS.Send(p.Number, receiver, msg)
	return c.Send(format.MessageSent(p.Number, receiver), tele.ModeHTML)
}

type Contacts struct {
	Universe *game.Universe
}

var contactsRe = re("^!контакты")

func (h *Contacts) Match(s string) bool {
	return contactsRe.MatchString(s)
}

func (h *Contacts) Handle(c tele.Context) error {
	world, _ := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	contacts := []format.Contact{}
	for _, u := range world.Users {
		if p, ok := u.Phone(); ok {
			member := teleutil.Member(c, tele.ChatID(u.TUID))
			contacts = append(contacts, format.Contact{
				Name:   teleutil.Name(member),
				Number: p.Number,
			})
		}
	}
	return c.Send(format.Contacts(contacts), tele.ModeHTML)
}

type Spam struct {
	Universe *game.Universe
}

var spamRe = re("^!(спам|рассылка) (.*)")

func (h *Spam) Match(s string) bool {
	return spamRe.MatchString(s)
}

func (h *Spam) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	p, ok := user.Phone()
	if !ok {
		return c.Send(format.NoPhone)
	}

	msg := teleutil.Args(c, spamRe)[2]
	if utf8.RuneCountInString(msg) > smsMaxLen {
		return c.Send(format.SMSMaxLen(smsMaxLen))
	}
	const price = 2000
	if !user.Balance().Spend(price) {
		return c.Send(format.NoMoney)
	}
	for _, u := range world.Users {
		if q, ok := u.Phone(); ok {
			world.SMS.Send(p.Number, q.Number, msg)
		}
	}
	return c.Send(format.SpamSent(price), tele.ModeHTML)
}
