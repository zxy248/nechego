package handlers

import (
	"errors"
	"fmt"
	"nechego/auction"
	"nechego/format"
	"nechego/game"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Auction struct {
	Universe *game.Universe
}

var auctionRe = re("^!аук")

func (h *Auction) Match(s string) bool {
	return auctionRe.MatchString(s)
}

func (h *Auction) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	lots, markup := auctionMessage(world)
	return c.Send(lots, markup, tele.ModeHTML)
}

func auctionMessage(w *game.World) (string, *tele.ReplyMarkup) {
	encode := func(l *auction.Lot) string {
		c := auctionCallback{l.Key}
		return c.encode()
	}
	return format.Auction(w.Auction.List(), encode)
}

type AuctionBuy struct {
	Universe *game.Universe
}

func (h *AuctionBuy) Match(s string) bool {
	return callbackMatch(&auctionCallback{}, s)
}

func (h *AuctionBuy) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	var data auctionCallback
	data.decode(c.Callback().Data)
	lot, err := user.AuctionBuy(world, data.lotKey)
	if errors.Is(err, game.ErrNoKey) {
		return c.Send(format.NoLot)
	} else if errors.Is(err, game.ErrNoMoney) {
		return c.Send(format.NoMoney)
	} else if err != nil {
		return err
	}

	lots, markup := auctionMessage(world)
	editErr := c.Edit(lots, markup, tele.ModeHTML)
	sendErr := c.Send(format.AuctionBought(
		tu.Mention(c, user.TUID),
		tu.Mention(c, lot.SellerID),
		lot.Price(),
		lot.Item,
	), tele.ModeHTML)
	return errors.Join(editErr, sendErr)
}

const auctionCallbackFormat = "/auction/%d"

type auctionCallback struct {
	lotKey int
}

func (c *auctionCallback) encode() string {
	return fmt.Sprintf(auctionCallbackFormat, c.lotKey)
}

func (c *auctionCallback) decode(s string) error {
	_, err := fmt.Sscanf(s, auctionCallbackFormat, &c.lotKey)
	return err
}

type AuctionSell struct {
	Universe *game.Universe
}

var auctionSell = re("^!торг (.*)")

func (h *AuctionSell) Match(s string) bool {
	return auctionSell.MatchString(s)
}

func (h *AuctionSell) Handle(c tele.Context) error {
	error := func() error {
		n := format.NewConnector("\n")
		n.Add("🏦 <b>Помощь:</b>")
		n.Add("<code>!торг &lt;номер предмета&gt; &lt;начальная цена от 1000 ₴&gt;</code>")
		return c.Send(n.String(), tele.ModeHTML)
	}

	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	args := tu.NumArg(c, auctionSell, 1)
	if len(args) != 2 {
		return error()
	}
	key, price := args[0], args[1]

	item, ok := user.Inventory.ByKey(key)
	if !ok {
		return c.Send(format.BadKey(key), tele.ModeHTML)
	}
	if world.Auction.Full() {
		return c.Send(format.AuctionFull)
	}
	if err := user.AuctionSell(world, item, price); err != nil {
		return error()
	}
	return c.Send(format.AuctionSell)
}
