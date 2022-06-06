package main

import (
	"regexp"
	"strings"
)

var eblanRe = regexp.MustCompile("^![ие][б6п*]?лан дня")
var masyunyaRe = regexp.MustCompile("^!ма[нс]юн")

type command int

const (
	commandUnknown command = iota
	commandProbability
	commandWho
	commandCat
	commandTitle
	commandAnime
	commandFurry
	commandFlag
	commandPerson
	commandHorse
	commandArt
	commandCar
	commandPair
	commandEblan
	commandMasyunya
	commandTurnOn
	commandTurnOff
)

// recognizeCommand returns the command contained in the input string.
func recognizeCommand(s string) command {
	switch {
	case startsWith(s, "!инф"):
		return commandProbability
	case startsWith(s, "!кто"):
		return commandWho
	case startsWith(s, "!кот", "!кош"):
		return commandCat
	case startsWith(s, "!имя"):
		return commandTitle
	case startsWith(s, "!аним", "!мульт"):
		return commandAnime
	case startsWith(s, "!фур"):
		return commandFurry
	case startsWith(s, "!флаг"):
		return commandFlag
	case startsWith(s, "!чел"):
		return commandPerson
	case startsWith(s, "!лошадь", "!конь"):
		return commandHorse
	case startsWith(s, "!арт", "!пик"):
		return commandArt
	case startsWith(s, "!авто", "!тачк", "!машин"):
		return commandCar
	case startsWith(s, "!пара дня"):
		return commandPair
	case eblanRe.MatchString(s):
		return commandEblan
	case masyunyaRe.MatchString(s):
		return commandMasyunya
	case startsWith(s, "!вкл"):
		return commandTurnOn
	case startsWith(s, "!выкл"):
		return commandTurnOff
	}
	return commandUnknown
}

// message represents a message in a group.
type message struct {
	text    string
	command command
}

// newMessage creates a new message from the text.
func newMessage(text string) *message {
	return &message{text, recognizeCommand(text)}
}

// argument returns the argument probably contained in the message.
func (m *message) argument() string {
	_, s, _ := strings.Cut(m.text, " ")
	return s
}

// startsWith returns true if the input string starts with one of the specified
// prefixes; false otherwise.
func startsWith(s string, prefix ...string) bool {
	for _, p := range prefix {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
