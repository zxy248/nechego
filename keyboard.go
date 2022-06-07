package main

import tele "gopkg.in/telebot.v3"

const buttonMasyunyaText = "🎀 Масюня 🎀"

func newKeyboard() *tele.ReplyMarkup {
	k := &tele.ReplyMarkup{ResizeKeyboard: true}
	b := k.Text(buttonMasyunyaText)
	k.Reply(k.Row(b))
	return k
}
