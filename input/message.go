package input

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Message represents an input message.
type Message struct {
	Raw     string
	Command Command
}

// ParseMessage returns a new message.
func ParseMessage(s string) *Message {
	return &Message{s, ParseCommand(s)}
}

// Argument returns an argument probably contained in the message.
func (m *Message) Argument() string {
	switch m.Command {
	case CommandWeather:
		return weatherRe.FindStringSubmatch(m.Raw)[1]
	case CommandProbability:
		return probabilityRe.FindStringSubmatch(m.Raw)[1]
	case CommandWho:
		return whoRe.FindStringSubmatch(m.Raw)[1]
	case CommandList:
		return listRe.FindStringSubmatch(m.Raw)[1]
	case CommandTop:
		return topRe.FindStringSubmatch(m.Raw)[1]
	}
	_, s, _ := strings.Cut(m.Raw, " ")
	return s
}

// Dynamic returns an argument that is not just a string.
func (m *Message) Dynamic() (interface{}, error) {
	switch m.Command {
	case CommandTop:
		return m.topArgument()
	case CommandForbid, CommandPermit:
		return m.commandActionArgument()
	case CommandTransfer:
		return m.transferArgument()
	case CommandDice:
		return m.diceArgument()
	}
	return nil, fmt.Errorf("no dynamic argument for %v", m.Raw)
}

func (m *Message) topArgument() (TopArgument, error) {
	match := topRe.FindStringSubmatch(m.Raw)
	number := match[1]
	desc := match[2]

	maybe, err := strconv.ParseInt(number, 10, 32)
	if err != nil {
		return TopArgument{nil, desc}, nil
	}
	i := int(maybe)
	return TopArgument{&i, desc}, nil
}

func (m *Message) commandActionArgument() (Command, error) {
	s := m.Argument()
	if s == "" {
		return CommandUnknown, ErrNoCommand
	}
	if !strings.HasPrefix(s, "!") {
		s = "!" + s
	}
	c := ParseCommand(s)
	if c == CommandUnknown {
		return CommandUnknown, ErrUnknownCommand
	}
	return c, nil
}

var (
	ErrSpecifyAmount = errors.New("specify amount")
	ErrNotPositive   = errors.New("not positive")
)

func (m *Message) transferArgument() (int, error) {
	s := m.Argument()
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, ErrSpecifyAmount
	}
	if n <= 0 {
		return 0, ErrNotPositive
	}
	return int(n), nil
}

func (m *Message) diceArgument() (int, error) {
	s := m.Argument()
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, ErrSpecifyAmount
	}
	if n <= 0 {
		return 0, ErrNotPositive
	}
	return int(n), nil
}

// TopArgument represents an argument of the CommandTop.
type TopArgument struct {
	Number *int
	String string
}

var (
	ErrNoCommand      = errors.New("no command")
	ErrUnknownCommand = errors.New("unknown command")
)
