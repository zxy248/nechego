package input

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Message represents an input Message.
type Message struct {
	Raw     string
	Command Command
}

// ParseMessage returns a new message.
func ParseMessage(s string) *Message {
	return &Message{s, ParseCommand(s)}
}

// Argument returns the argument probably contained in the message.
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

func (m *Message) DynamicArgument() (interface{}, error) {
	switch m.Command {
	case CommandTop:
		m := topRe.FindStringSubmatch(m.Raw)
		n := m[1]
		s := m[2]

		i, err := strconv.ParseInt(n, 10, 32)
		if err != nil {
			return TopArgument{false, 0, s}, nil
		}
		return TopArgument{true, int(i), s}, nil
	case CommandForbid, CommandPermit:
		s := m.Argument()
		if s == "" {
			return nil, ErrNoCommand
		}
		if !strings.HasPrefix(s, "!") {
			s = "!" + s
		}
		c := ParseCommand(s)
		if c == CommandUnknown {
			return nil, ErrUnknownCommand
		}
		return c, nil
	}
	return nil, fmt.Errorf("no dynamic argument for %v", m.Raw)
}

// TopArgument represents an argument of the CommandTop.
type TopArgument struct {
	NumberPresent bool
	Number        int
	String        string
}

var (
	ErrNoCommand      = errors.New("no command")
	ErrUnknownCommand = errors.New("unknown command")
)
