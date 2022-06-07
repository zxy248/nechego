package main

import (
	"log"

	tele "gopkg.in/telebot.v3"
)

// processInput processes the input message and handles it to the corresponding
// function. Ignores the message if the chat type is not a group. Caches the
// group ID and the user ID.
func (a *app) processInput(c tele.Context) error {
	if !chatTypeIsGroup(c.Chat().Type) {
		return nil
	}

	text := c.Text()
	groupID := c.Chat().ID
	userID := c.Sender().ID
	message := newMessage(text)

	if !a.whitelist.allow(groupID) {
		return nil
	}

	log.Printf("%d@%d: %s\n", userID, groupID, text)
	a.cacheGroupMember(groupID, userID)

	if !a.status.activeGlobal() {
		return nil
	}

	if !a.status.activeLocal(groupID) {
		if message.command == commandTurnOn {
			return a.handleTurnOn(c)
		}
		return nil
	}

	return a.routeMessage(c, message)
}

// routeMessage routes the input message to the appropriate handler.
func (a *app) routeMessage(c tele.Context, m *message) error {
	switch m.command {
	case commandProbability:
		return a.handleProbability(c, m)
	case commandWho:
		return a.handleWho(c, m)
	case commandCat:
		return a.handleCat(c)
	case commandTitle:
		return a.handleTitle(c, m)
	case commandAnime:
		return a.handleAnime(c)
	case commandFurry:
		return a.handleFurry(c)
	case commandFlag:
		return a.handleFlag(c)
	case commandPerson:
		return a.handlePerson(c)
	case commandHorse:
		return a.handleHorse(c)
	case commandArt:
		return a.handleArt(c)
	case commandCar:
		return a.handleCar(c)
	case commandPair:
		return a.handlePair(c)
	case commandEblan:
		return a.handleEblan(c)
	case commandMasyunya:
		return a.handleMasyunya(c)
	case commandKeyboardOpen:
		return a.handleKeyboardOpen(c)
	case commandKeyboardClose:
		return a.handleKeyboardClose(c)
	case commandTurnOff:
		return a.handleTurnOff(c)
	}
	return nil
}

// cacheGroupMember adds the user to the group if he is not already there.
func (a *app) cacheGroupMember(groupID int64, userID int64) error {
	ids, err := a.store.getUserIDs(groupID)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if id == userID {
			return nil
		}
	}
	if err := a.store.insertUserID(groupID, userID); err != nil {
		return err
	}
	return nil
}

// chatTypeIsGroup returns true if the chat type is one of the group types;
// false otherwise.
func chatTypeIsGroup(t tele.ChatType) bool {
	return t == tele.ChatGroup || t == tele.ChatSuperGroup
}
