package main

import (
	"log"
	"regexp"
	"strings"

	tele "gopkg.in/telebot.v3"
)

var eblanRe = regexp.MustCompile("![ие][б6п*]?лан дня")

// handleMessage processes the input message and handles it to the corresponding
// function. Ignores the message if the chat type is not a group. Caches the
// group ID and the user ID.
func (a *app) handleMessage(c tele.Context) error {
	if !chatTypeIsGroup(c.Chat().Type) {
		return nil
	}

	message := c.Text()
	groupID := c.Chat().ID
	userID := c.Sender().ID

	log.Printf("%d@%d: %s\n", userID, groupID, message)

	a.cacheGroupMember(groupID, userID)

	if !a.status.active() {
		if strings.HasPrefix(message, "!вкл") {
			return a.handleTurnOn(c)
		}
		return nil
	}
	if strings.HasPrefix(message, "!выкл") {
		return a.handleTurnOff(c)
	}

	switch {
	case strings.HasPrefix(message, "!инфа"):
		arg := getCommandArgument(message, "!инфа")
		return a.handleProbability(c, arg)
	case strings.HasPrefix(message, "!кто"):
		arg := getCommandArgument(message, "!кто")
		return a.handleWho(c, arg)
	case strings.HasPrefix(message, "!кот") || strings.HasPrefix(message, "!кош"):
		return a.handleCat(c)
	case strings.HasPrefix(message, "!имя"):
		arg := getCommandArgument(message, "!имя")
		return a.handleTitle(c, arg)
	case strings.HasPrefix(message, "!аним") || strings.HasPrefix(message, "!мульт"):
		return a.handleAnime(c)
	case strings.HasPrefix(message, "!фур"):
		return a.handleFurry(c)
	case strings.HasPrefix(message, "!флаг"):
		return a.handleFlag(c)
	case strings.HasPrefix(message, "!чел"):
		return a.handlePerson(c)
	case strings.HasPrefix(message, "!лошадь") || strings.HasPrefix(message, "!конь"):
		return a.handleHorse(c)
	case strings.HasPrefix(message, "!арт") || strings.HasPrefix(message, "!пик"):
		return a.handleArt(c)
	case strings.HasPrefix(message, "!авто") || strings.HasPrefix(message, "!тачк") || strings.HasPrefix(message, "!машин"):
		return a.handleCar(c)
	case strings.HasPrefix(message, "!пара дня"):
		return a.handlePair(c)
	case eblanRe.MatchString(message):
		return a.handleEblan(c)
	}
	return nil
}

// cacheGroupMember adds the user to the group if he is not already there
func (a *app) cacheGroupMember(groupID int64, userID int64) error {
	userIDs, err := a.store.getUserIDs(groupID)
	if err != nil {
		return err
	}
	for _, id := range userIDs {
		if id == userID {
			return nil
		}
	}
	if err := a.store.insertUserID(groupID, userID); err != nil {
		return err
	}
	return nil
}

// getCommandArgument returns the part of the message after the prefix
func getCommandArgument(message, prefix string) string {
	return strings.TrimSpace(strings.TrimPrefix(message, prefix))
}

// chatTypeIsGroup returns true if the chat type is one of the group types
func chatTypeIsGroup(t tele.ChatType) bool {
	return t == tele.ChatGroup || t == tele.ChatSuperGroup
}
