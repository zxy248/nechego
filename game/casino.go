package game

import (
	"errors"
	"time"
)

// DiceGameOutcome is either a Draw, a Win, or a Lose.
type DiceGameOutcome int

const (
	Draw DiceGameOutcome = iota
	Win
	Lose
)

// DiceGameResult represents a result of a dice game.
type DiceGameResult struct {
	Prize   int
	Outcome DiceGameOutcome
}

// DiceGame represents an ongoing dice game.
type DiceGame struct {
	playerID    int64
	casinoScore int
	bet         int
	finish      func()
}

// Verify checks if the given player ID matches the ID of a game.
func (d *DiceGame) Verify(playerID int64) bool {
	return playerID == d.playerID
}

// Going returns true if the game is in progress.
func (d *DiceGame) Going() bool {
	return d != nil
}

// Finish calculates a game result from the specified player score and
// stops the game.
func (d *DiceGame) Finish(playerScore int) DiceGameResult {
	d.finish()
	outcome := diceGameOutcome(playerScore, d.casinoScore)
	prize := diceGamePrize(d.bet, outcome)
	return DiceGameResult{prize, outcome}
}

func diceGameOutcome(playerScore, casinoScore int) DiceGameOutcome {
	if playerScore > casinoScore {
		return Win
	}
	if playerScore < casinoScore {
		return Lose
	}
	return Draw
}

func diceGamePrize(bet int, o DiceGameOutcome) int {
	if o == Win {
		return 2 * bet
	}
	if o == Lose {
		return 0
	}
	return bet
}

// Casino holds the current dice game.
type Casino struct {
	Timeout  time.Duration
	diceGame *DiceGame
}

// DiceThrowFunc represents a function provided by caller that is used
// to get the casino score.
type DiceThrowFunc func() (score int, err error)

// PlayDice starts a dice game. If a game is already going, returns an
// error.
func (c *Casino) PlayDice(playerID int64, bet int, throw DiceThrowFunc, timeout func()) error {
	if c.diceGame.Going() {
		return errors.New("casino: game already going")
	}
	score, err := throw()
	if err != nil {
		return err
	}
	done := make(chan struct{}, 1)
	c.diceGame = &DiceGame{
		playerID:    playerID,
		casinoScore: score,
		bet:         bet,
		finish:      func() { done <- struct{}{} },
	}
	go func() {
		timer := time.NewTimer(c.Timeout)
		select {
		case <-timer.C:
			timeout()
		case <-done:
			timer.Stop()
		}
		c.diceGame = nil
	}()
	return nil
}

// Game returns the current dice game.
func (c *Casino) Game() *DiceGame {
	return c.diceGame
}
