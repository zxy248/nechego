package bot

import tele "gopkg.in/telebot.v3"

const buttonMasyunyaText = "🎀 Масюня 🎀"

// keyboard returns a new keyboard.
func keyboard() *tele.ReplyMarkup {
	k := &tele.ReplyMarkup{ResizeKeyboard: true}
	b := k.Text(buttonMasyunyaText)
	k.Reply(k.Row(b))
	return k
}
