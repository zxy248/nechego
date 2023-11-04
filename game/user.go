package game

import (
	"nechego/buff"
	"nechego/farm"
	"nechego/fishing"
	"nechego/game/reputation"
	"nechego/item"
	"time"
)

// User represents a player.
type User struct {
	TUID        int64        // Telegram ID.
	Energy      Energy       // Energy level.
	Rating      float64      // Elo rating in fights.
	Messages    int          // Number of messages sent.
	BannedUntil time.Time    // Time after which the user is unbanned.
	Status      string       // Status displayed in the profile.
	Inventory   *item.Set    // Personal items.
	Net         *fishing.Net // Net if currently cast, else nil.
	LastMessage time.Time    // When was the last message sent?
	Buffs       buff.Set     // Active buffs.
	Developer   bool         // Flag of a game developer.
	Funds       Funds        // Collectable items.
	Farm        *farm.Farm   // The source of vegetables.
	Retired     time.Time    // When the job shift should finish.
	Friends     Friends      // The list of friends' TUIDs.
	SlotBet     int          // The bet for slots.
	Reputation  reputation.Reputation
}

func NewUser(tuid int64) *User {
	return &User{
		TUID:      tuid,
		Rating:    1500,
		Inventory: item.NewSet(),
		Buffs:     buff.Set{},
		Funds:     Funds{},
		Farm:      farm.New(2, 3),
		Friends:   Friends{},
	}
}

// ID returns the unique user's ID.
func (u *User) ID() int64 {
	return u.TUID
}

// Transfer moves the item x from the sender's inventory to the
// receiver's funds.
func (u *User) Transfer(to *User, x *item.Item) bool {
	if !x.Transferable || !u.Inventory.Remove(x) {
		return false
	}
	to.Funds.Add("обмен", x)
	return true
}
