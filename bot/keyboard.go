package bot

import tele "gopkg.in/telebot.v3"

const buttonMasyunyaText = "Масюня 🎀"
const buttonPoppyText = "Паппи 🦊"

// keyboard returns a new keyboard.
func keyboard() *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{ResizeKeyboard: true}
	masyunyaBtn := kb.Text(buttonMasyunyaText)
	poppyBtn := kb.Text(buttonPoppyText)
	kb.Reply(kb.Row(masyunyaBtn, poppyBtn))
	return kb
}
