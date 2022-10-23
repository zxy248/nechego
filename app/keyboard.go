package app

import tele "gopkg.in/telebot.v3"

const (
	buttonMasyunyaText = "Масюня 🎀"
	buttonPoppyText    = "Паппи 🦊"
	buttonSimaText     = "Сима 💖"
)

var keyboard = func() *tele.ReplyMarkup {
	k := &tele.ReplyMarkup{ResizeKeyboard: true}
	buttonSima := k.Text(buttonSimaText)
	buttonMasyunya := k.Text(buttonMasyunyaText)
	buttonPoppy := k.Text(buttonPoppyText)
	k.Reply(k.Row(buttonSima, buttonMasyunya, buttonPoppy))
	return k
}()

func openKeyboard(c tele.Context) error {
	return c.Send("Клавиатура ⌨️", keyboard)
}

func closeKeyboard(c tele.Context) error {
	return c.Send("Клавиатура закрыта 😣", tele.RemoveKeyboard)
}
