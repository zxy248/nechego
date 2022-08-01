package dice

import (
	"errors"
	"nechego/model"
	"sync"
	"time"
)

var (
	ErrGameInProgress = errors.New("game in progress")
	ErrNoGame         = errors.New("no game")
	ErrWrongUser      = errors.New("wrong user")
)

type Actions struct {
	Throw   func() (int, error)
	Timeout func()
}

type Game struct {
	Start time.Time
	User  model.User
	Bet   int
	Roll  int
}

func newGame(u model.User, bet, roll int) *Game {
	return &Game{time.Now(), u, bet, roll}
}

type Settings struct {
	RollTime time.Duration
}

type Casino struct {
	games    *sync.Map // map[GID]*Game
	Settings Settings
}

func New(s Settings) *Casino {
	return &Casino{&sync.Map{}, s}
}

func (c *Casino) Play(g model.Group, u model.User, bet int, a Actions) error {
	if _, loaded := c.loadGame(g); loaded {
		return ErrGameInProgress
	}
	roll, err := a.Throw()
	if err != nil {
		return err
	}
	game := newGame(u, bet, roll)
	c.saveGame(g, game)
	go func() {
		time.Sleep(c.Settings.RollTime)
		deleted := c.deleteGame(g, game)
		if deleted {
			a.Timeout()
		}
	}()
	return nil
}

type Outcome int

const (
	Win Outcome = iota
	Draw
	Lose
)

type Result struct {
	Outcome
	*Game
}

func (c *Casino) Roll(g model.Group, u model.User, roll int) (Result, error) {
	game, ok := c.loadGame(g)
	if !ok {
		return Result{Draw, game}, ErrNoGame
	}
	if game.User.ID != u.ID {
		return Result{Draw, game}, ErrWrongUser
	}
	c.deleteGame(g, game)
	switch {
	case roll > game.Roll:
		return Result{Win, game}, nil
	case roll == game.Roll:
		return Result{Draw, game}, nil
	default:
		return Result{Lose, game}, nil
	}
}

func (c *Casino) loadGame(g model.Group) (*Game, bool) {
	game, ok := c.games.Load(g.GID)
	if ok {
		return game.(*Game), true
	}
	return nil, false
}

func (c *Casino) saveGame(g model.Group, game *Game) error {
	_, loaded := c.games.LoadOrStore(g.GID, game)
	if loaded {
		return ErrGameInProgress
	}
	return nil
}

func (c *Casino) deleteGame(g model.Group, game *Game) bool {
	loadedGame, loaded := c.loadGame(g)
	if !loaded {
		return false
	}
	if loadedGame.Start != game.Start {
		return false
	}
	c.games.Delete(g.GID)
	return true
}
