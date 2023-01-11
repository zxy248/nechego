package game

import (
	"errors"
	"time"
)

type DiceGame struct {
	Player      *User
	CasinoScore int
	Bet         int
	time        time.Time
	done        chan struct{}
}

func (d *DiceGame) Finish() {
	d.done <- struct{}{}
}

type Casino struct {
	Timeout  time.Duration
	diceGame *DiceGame
}

func (c *Casino) PlayDice(player *User, bet int, throw func() (int, error), timeout func()) error {
	if c.GameGoing() {
		return errors.New("game already going")
	}
	score, err := throw()
	if err != nil {
		return err
	}
	c.diceGame = &DiceGame{
		Player:      player,
		CasinoScore: score,
		Bet:         bet,
		time:        time.Now(),
		done:        make(chan struct{}, 1),
	}
	go func() {
		timer := time.NewTimer(c.Timeout)
		select {
		case <-timer.C:
			timeout()
		case <-c.diceGame.done:
			timer.Stop()
		}
		c.diceGame = nil
	}()
	return nil
}

func (c *Casino) DiceGame() (d *DiceGame, ok bool) {
	if c.GameGoing() {
		return c.diceGame, true
	}
	return nil, false
}

func (c *Casino) GameGoing() bool {
	return c.diceGame != nil
}
