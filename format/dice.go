package format

import (
	"fmt"
	"nechego/game"
)

const (
	NoDice    = "🎲 У вас нет костей."
	GameGoing = "🎲 Игра уже идет."
)

func DiceGame(mention string, bet int, seconds int) string {
	c := NewConnector("\n")
	c.Add(fmt.Sprintf("🎲 %s играет на %s", Name(mention), Money(bet)))
	c.Add(fmt.Sprintf("У вас <code>%d секунд</code> на то, чтобы бросить кости!", seconds))
	return c.String()
}

func DiceGameResult(r game.DiceGameResult) string {
	if r.Outcome == game.Win {
		return fmt.Sprintf("💥 Вы выиграли %s", Money(r.Prize))
	}
	if r.Outcome == game.Lose {
		return "😵 Вы проиграли."
	}
	return "🎲 Ничья."
}

func DiceTimeout(bet int) string {
	return fmt.Sprintf("<i>⏰ Время вышло: вы потеряли %s</i>", Money(bet))
}

func MinBet(n int) string {
	return fmt.Sprintf("💵 Минимальная ставка %s", Money(n))
}

func SlotWin(mention string, prize int) string {
	return fmt.Sprintf("🎰 %s выигрывает %s 💥", Name(mention), Money(prize))
}

func SlotRoll(mention string, bet int) string {
	return fmt.Sprintf("🎰 %s прокручивает слоты на %s", Name(mention), Money(bet))
}

func BetSet(mention string, n int) string {
	return fmt.Sprintf("🎰 %s устанавливает ставку %s", Name(mention), Money(n))
}
