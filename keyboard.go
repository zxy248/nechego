package main

import tele "gopkg.in/telebot.v3"

const buttonMasyunyaText = "🎀 Масюня 🎀"

// newKeyboard returns a new keyboard.
func newKeyboard() *tele.ReplyMarkup {
	k := &tele.ReplyMarkup{ResizeKeyboard: true}
	b := k.Text(buttonMasyunyaText)
	k.Reply(k.Row(b))
	return k
}
