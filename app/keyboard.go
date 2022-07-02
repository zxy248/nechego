package app

import tele "gopkg.in/telebot.v3"

const buttonMasyunyaText = "Масюня 🎀"
const buttonPoppyText = "Паппи 🦊"

var keyboard = func() *tele.ReplyMarkup {
	k := &tele.ReplyMarkup{ResizeKeyboard: true}
	buttonMasyunya := k.Text(buttonMasyunyaText)
	buttonPoppy := k.Text(buttonPoppyText)
	k.Reply(k.Row(buttonMasyunya, buttonPoppy))
	return k
}()
