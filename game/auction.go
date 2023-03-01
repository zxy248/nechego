package game

import (
	"errors"
	"nechego/auction"
	"nechego/item"
	"nechego/money"
	"time"
)

func (u *User) AuctionBuy(w *World, key int) (l *auction.Lot, err error) {
	lot, ok := w.Auction.Get(key)
	if !ok {
		return nil, ErrNoKey
	}
	cost := lot.Price()
	if !u.Balance().Spend(cost) {
		return nil, ErrNoMoney
	}
	w.Auction.Remove(key)
	u.Inventory.Add(lot.Item)
	seller := w.UserByID(lot.SellerID)
	seller.Funds.Add("аукцион", item.New(&money.Cash{Money: lot.Price()}))
	return lot, nil
}

func (u *User) AuctionSell(w *World, i *item.Item, initialPrice int) error {
	if !i.Transferable {
		return errors.New("item is not transferable")
	}
	if !u.Inventory.Remove(i) {
		return errors.New("item cannot be removed")
	}
	lot := &auction.Lot{
		SellerID: u.TUID,
		Item:     i,
		MinPrice: 0,
		MaxPrice: initialPrice,
		Duration: 2 * time.Hour,
	}
	if err := w.Auction.Place(lot); err != nil {
		u.Inventory.Add(i)
		return err
	}
	return nil
}
