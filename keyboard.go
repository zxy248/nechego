package main

import tele "gopkg.in/telebot.v3"

var (
	keyboard    = &tele.ReplyMarkup{ResizeKeyboard: true}
	btnMasyunya = keyboard.Text("🎀 Масюня 🎀")
)

func initializeKeyboard() {
	keyboard.Reply(keyboard.Row(btnMasyunya))
}
