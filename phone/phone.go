package phone

import (
	"fmt"
	"nechego/modifier"
	"time"
)

type Database map[Number][]*SMS

func (db Database) Count(n Number) int {
	return len(db[n])
}

func (db Database) Receive(n Number) []*SMS {
	r, ok := db[n]
	if !ok {
		return []*SMS{}
	}
	delete(db, n)
	return r
}

func (db Database) Send(sender, receiver Number, text string) {
	msgs, ok := db[receiver]
	if !ok {
		msgs = []*SMS{}
	}
	msgs = append(msgs, NewSMS(sender, text))
	db[receiver] = msgs
}

type SMS struct {
	Time   time.Time
	Sender Number
	Text   string
}

func NewSMS(sender Number, text string) *SMS {
	return &SMS{
		Time:   time.Now(),
		Sender: sender,
		Text:   text,
	}
}

type Phone struct {
	Number Number
}

func NewPhone() *Phone {
	return &Phone{RandomNumber()}
}

func (p *Phone) String() string {
	return fmt.Sprintf("📱 Смартфон (%s)", p.Number)
}

func (p *Phone) Mod() (m *modifier.Mod, ok bool) {
	return &modifier.Mod{
		Emoji:       "📱",
		Multiplier:  +0.05,
		Description: "У вас есть смартфон.",
	}, true
}
