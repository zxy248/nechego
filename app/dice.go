package app

import (
	"errors"
	"fmt"
	"math/rand"
	"nechego/input"
	"nechego/model"
	"sync"
	"time"

	tele "gopkg.in/telebot.v3"
)

// map[int64]diceGame
var diceGames = &sync.Map{}

type diceGame struct {
	id    time.Time
	user  model.User
	money int
	roll  int
}

func makeDiceGame(u model.User, money, roll int) diceGame {
	return diceGame{
		id:    time.Now(),
		user:  u,
		money: money,
		roll:  roll,
	}
}

func (g diceGame) key() int64 {
	return g.user.GID
}

func currentDiceGame(g model.Group) (diceGame, bool) {
	return loadDiceGame(g.GID)
}

func loadDiceGame(key int64) (diceGame, bool) {
	game, ok := diceGames.Load(key)
	if !ok {
		return diceGame{}, false
	}
	return game.(diceGame), true
}

func (g diceGame) storeDiceGame() (ok bool) {
	_, loaded := diceGames.LoadOrStore(g.key(), g)
	return !loaded
}

func (g diceGame) finish() {
	diceGames.Delete(g.key())
}

func (g diceGame) startDiceGame(notify func()) error {
	ok := g.storeDiceGame()
	if !ok {
		return errors.New("game already going")
	}
	time.AfterFunc(diceRollTime, func() { g.cancelDiceGame(notify) })
	return nil
}

func (g diceGame) cancelDiceGame(notify func()) {
	game, ok := loadDiceGame(g.key())
	if ok {
		if g.id == game.id {
			g.finish()
			notify()
		}
	}
}

const (
	diceStart       = "🎲 %s играет на `%s 💰`\nУ вас `%d секунд` на то, чтобы кинуть кости\\!"
	diceWin         = "💥 Вы выиграли `%v 💰`"
	diceDraw        = "Ничья."
	diceLose        = "Вы проиграли."
	diceBonus       = "_🎰 *%s* получает бонус за риск: `%s 💰`_"
	diceTimeout     = "Время вышло: вы потеряли `%s`\\."
	diceMinBonus    = 1
	diceMaxBonus    = 5
	diceBetForBonus = 5
	diceBonusChance = 0.2
	diceRollTime    = time.Second * 25
	diceMinBet      = 1
)

var handleDiceMutex = &sync.Mutex{}

// handleDice rolls a dice.
func (a *App) handleDice(c tele.Context) error {
	handleDiceMutex.Lock()
	defer handleDiceMutex.Unlock()
	group := getGroup(c)
	user := getUser(c)

	_, ok := currentDiceGame(group)
	if ok {
		return c.Send(makeError("Игра уже идет"))
	}

	arg, err := getMessage(c).Dynamic()
	if err != nil {
		if errors.Is(err, input.ErrSpecifyAmount) {
			return c.Send(makeError("Укажите количество средств"))
		}
		if errors.Is(err, input.ErrNotPositive) {
			return c.Send(makeError("Некорректная ставка"))
		}
		return err
	}
	bet := arg.(int)
	if bet < diceMinBet {
		return c.Send(makeError("Поставьте больше средств"))
	}

	ok = a.model.UpdateMoney(user, -bet)
	if !ok {
		return c.Send(makeError("Недостаточно средств"))
	}

	dice := &tele.Dice{Type: tele.Cube.Type}
	msg, err := dice.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
	if err != nil {
		return err
	}
	roll := msg.Dice.Value

	game := makeDiceGame(user, bet, roll)
	game.startDiceGame(func() {
		c.Send(fmt.Sprintf(diceTimeout,
			formatAmount(game.money)),
			tele.ModeMarkdownV2)
	})

	out := fmt.Sprintf(diceStart,
		a.mustMentionUser(user), formatAmount(bet, diceRollTime/time.Second)
	return c.Send(out, tele.ModeMarkdownV2)
}

func (a *App) handleRoll(c tele.Context) error {
	group := getGroup(c)
	user := getUser(c)

	game, ok := currentDiceGame(group)
	if !ok {
		return nil
	}
	if game.user.ID != user.ID {
		return nil
	}
	game.finish()

	defer func() {
		if rand.Float64() <= diceBonusChance && game.money >= diceBetForBonus {
			bonus := randInRange(diceMinBonus, diceMaxBonus)
			a.model.UpdateMoney(user, bonus)
			c.Send(fmt.Sprintf(diceBonus,
				a.mustMentionUser(user), formatAmount(bonus)),
				tele.ModeMarkdownV2)
		}
	}()

	switch roll := c.Message().Dice.Value; {
	case roll > game.roll:
		win := game.money * 2
		a.model.UpdateMoney(user, win)
		return c.Send(fmt.Sprintf(diceWin, formatAmount(win)), tele.ModeMarkdownV2)
	case roll == game.roll:
		a.model.UpdateMoney(user, game.money)
		return c.Send(diceDraw)
	default:
		return c.Send(diceLose)
	}
}
