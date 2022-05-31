package main

import (
	"log"
	"strings"

	tele "gopkg.in/telebot.v3"
)

// handleMessage processes the input message and handles it to the corresponding
// function. Ignores the message if the chat type is not a group. Caches the
// group ID and the user ID.
func (a *app) handleMessage(c tele.Context) error {
	text := c.Text()
	log.Printf("Msg: %s\n", text)

	if !chatTypeIsGroup(c.Chat().Type) {
		return nil
	}

	a.cacheGroupMember(c.Chat().ID, c.Sender().ID)

	switch {
	case strings.HasPrefix(text, "!инфа"):
		message := getCommandArgument(text, "!инфа")
		return a.handleProbability(c, message)
	case strings.HasPrefix(text, "!кто"):
		message := getCommandArgument(text, "!кто")
		return a.handleWho(c, message)
	case strings.HasPrefix(text, "!кот") || strings.HasPrefix(text, "!кош"):
		return a.handleCat(c)
	case strings.HasPrefix(text, "!имя"):
		message := getCommandArgument(text, "!имя")
		return a.handleTitle(c, message)
	case strings.HasPrefix(text, "!аним"):
		return a.handleAnime(c)
	case strings.HasPrefix(text, "!фур"):
		return a.handleFurry(c)
	case strings.HasPrefix(text, "!флаг"):
		return a.handleFlag(c)
	case strings.HasPrefix(text, "!чел"):
		return a.handlePerson(c)
	case strings.HasPrefix(text, "!лошадь") || strings.HasPrefix(text, "!конь"):
		return a.handleHorse(c)
	case strings.HasPrefix(text, "!арт") || strings.HasPrefix(text, "!пик"):
		return a.handleArt(c)
	case strings.HasPrefix(text, "!авто") || strings.HasPrefix(text, "!тачк"):
		return a.handleCar(c)
	case strings.HasPrefix(text, "!пара дня"):
		return a.handlePair(c)
	}

	return nil
}

// cacheGroupMember adds groupID and userID to the database
func (a *app) cacheGroupMember(groupID int64, userID int64) error {
	userIDs, err := a.store.getUserIDs(groupID)
	if err != nil {
		return err
	}
	for _, uid := range userIDs {
		if uid == userID {
			return nil
		}
	}
	if err := a.store.insertUserID(groupID, userID); err != nil {
		return err
	}
	return nil
}

// getCommandArgument returns the part of msg after cmd
func getCommandArgument(msg, cmd string) string {
	return strings.TrimSpace(strings.TrimPrefix(msg, cmd))
}

// chatTypeIsGroup returns true if the chat type is one of the group types
func chatTypeIsGroup(t tele.ChatType) bool {
	return t == tele.ChatGroup || t == tele.ChatSuperGroup
}
