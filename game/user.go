package game

import (
	"nechego/farm"
	"nechego/fishing"
	"nechego/game/reputation"
	"nechego/item"
	"time"
)

// User represents a player.
type User struct {
	ID               int64        // Telegram ID.
	Name             string       // Cached name.
	Energy           Energy       // Energy level.
	Rating           float64      // Elo rating in fights.
	Messages         int          // Number of messages sent.
	BannedUntil      time.Time    // Time after which the user is unbanned.
	Status           string       // Status displayed in the profile.
	Inventory        *item.Set    // Personal items.
	Net              *fishing.Net // Net if currently cast, else nil.
	LastMessage      time.Time    // When was the last message sent?
	Funds            Funds        // Collectable items.
	Farm             *farm.Farm   // The source of vegetables.
	Retired          time.Time    // When the job shift should finish.
	Friends          Friends      // The list of friends' TUIDs.
	SlotBet          int          // The bet for slots.
	Activity         float64      // Message count coefficient.
	Reputation       reputation.Reputation
	ReputationFactor float64
	RatingPosition   int
}

func NewUser(id int64) *User {
	u := &User{
		ID:        id,
		Rating:    1500,
		Inventory: item.NewSet(),
		Funds:     Funds{},
		Farm:      farm.New(2, 3),
		Friends:   Friends{},
	}
	u.Balance().Add(10000)
	return u
}

// Transfer moves the item x from the sender's inventory to the
// receiver's funds.
func (u *User) Transfer(to *User, x *item.Item) bool {
	if !x.Transferable {
		return false
	}
	u.Inventory.Remove(x)
	to.Funds.Add("обмен", x)
	return true
}
