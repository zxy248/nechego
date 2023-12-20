package game

import (
	"nechego/farm"
	"nechego/game/reputation"
	"nechego/item"
	"time"
)

// User represents a player.
type User struct {
	ID               int64      // Telegram ID.
	Name             string     // Cached name.
	Energy           Energy     // Energy level.
	Rating           float64    // Elo rating in fights.
	Messages         int        // Number of messages sent.
	Status           string     // Status displayed in the profile.
	Inventory        *item.Set  // Personal items.
	Mail             *item.Set  // Incoming mail.
	LastMessage      time.Time  // When was the last message sent?
	Farm             *farm.Farm // The source of vegetables.
	Friends          Friends    // The list of friends' TUIDs.
	SlotBet          int        // The bet for slots.
	Activity         float64    // Message count coefficient.
	Reputation       reputation.Reputation
	ReputationFactor float64
	RatingPosition   int
}

func NewUser(id int64) *User {
	u := &User{
		ID:        id,
		Rating:    1500,
		Inventory: item.NewSet(),
		Mail:      item.NewSet(),
		Farm:      farm.New(2, 3),
		Friends:   Friends{},
	}
	u.Balance().Add(5000)
	return u
}
