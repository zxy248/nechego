package auction

import (
	"errors"
	"nechego/item"
	"sort"
	"time"
)

// Lot represents a lot at the auction.
type Lot struct {
	// ID of the seller.
	SellerID int64
	// Actual item.
	Item *item.Item
	// Price interval.
	MinPrice, MaxPrice int
	// Duration after which the Lot is removed from the Auction
	// and the underlying item is returned to the seller.
	Duration time.Duration
	// Time when the Lot was placed at the auction.
	Created time.Time
	// Key (if the Lot is added or returned from the Auction)
	Key int
}

// Expire returns the time when the Lot should end.
func (l *Lot) Expire() time.Time {
	return l.Created.Add(l.Duration)
}

// Price returns the current price of the Lot.
func (l *Lot) Price() int {
	start, end := l.Created, l.Expire()
	now := time.Now()
	if now.After(end) {
		return l.MinPrice
	}
	level := float64(now.Sub(start)) / float64(end.Sub(start))
	price := float64(l.MaxPrice) - float64(l.MaxPrice-l.MinPrice)*level
	return int(price)
}

// Auction represents a list of Lots indexed by a set of unique keys.
type Auction struct {
	Lots        map[int]*Lot
	MaxLots     int
	MinDuration time.Duration
}

// New returns an empty Auction.
func New() *Auction {
	return &Auction{
		Lots:        map[int]*Lot{},
		MaxLots:     12,
		MinDuration: time.Minute,
	}
}

// Place validates the Lot, gives it a unique key and places it at the Auction.
func (a *Auction) Place(l *Lot) error {
	if err := a.validateLot(l); err != nil {
		return err
	}
	key, ok := a.freeKey()
	if !ok {
		return errors.New("too many lots")
	}
	l.Key = key
	l.Created = time.Now()
	a.Lots[key] = l
	return nil
}

// Get tries to get a Lot from the Auction by the given key.
func (a *Auction) Get(key int) (l *Lot, ok bool) {
	l, ok = a.Lots[key]
	return
}

// Remove removes the Lot specified by key from the Auction.
func (a *Auction) Remove(key int) {
	delete(a.Lots, key)
}

// List returns a sorted list of Lots at the Auction.
func (a *Auction) List() []*Lot {
	lots := []*Lot{}
	for _, l := range a.Lots {
		lots = append(lots, l)
	}
	sort.Slice(lots, func(i, j int) bool {
		return lots[i].Key < lots[j].Key
	})
	return lots
}

// Full returns true if there is no space at the Auction for new Lots.
func (a *Auction) Full() bool {
	_, ok := a.freeKey()
	return !ok
}

// validateLot checks if the specified Lot is correct.
func (a *Auction) validateLot(l *Lot) error {
	if l.SellerID == 0 {
		return errors.New("lot seller is not specified")
	}
	if l.Item == nil {
		return errors.New("lot item is not specified")
	}
	if l.MinPrice < 0 {
		return errors.New("min price is negative")
	}
	if l.MinPrice > l.MaxPrice {
		return errors.New("min price is greater than max price")
	}
	if l.Duration < a.MinDuration {
		return errors.New("lot duration is too short")
	}
	return nil
}

// freeKey returns the first free key.
// If there is no free keys left, returns (0, false).
func (a *Auction) freeKey() (key int, ok bool) {
	for i := 0; i < a.MaxLots; i++ {
		if _, ok := a.Lots[i]; !ok {
			return i, true
		}
	}
	return 0, false
}
