package input

import (
	"errors"
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

var (
	ErrWrongCommand   = errors.New("wrong command")
	ErrSpecifyAmount  = errors.New("specify amount")
	ErrNotPositive    = errors.New("not positive")
	ErrNoCommand      = errors.New("no command")
	ErrUnknownCommand = errors.New("unknown command")
)

// TopArgument represents an argument of the CommandTop.
type TopArgument struct {
	Number *int
	String string
}

func (m *Message) TopArgument() (TopArgument, error) {
	a := TopArgument{}
	if m.Command != CommandTop {
		return a, ErrWrongCommand
	}
	match := topRe.FindStringSubmatch(m.Raw)
	number := match[1]
	a.String = match[2]
	maybe, err := strconv.ParseInt(number, 10, 32)
	if err != nil {
		return a, nil
	}
	n := int(maybe)
	a.Number = &n
	return a, nil
}

func (m *Message) CommandActionArgument() (Command, error) {
	c := CommandUnknown
	if m.Command != CommandForbid && m.Command != CommandPermit {
		return c, ErrWrongCommand
	}
	arg := m.Argument()
	if arg == "" {
		return c, ErrNoCommand
	}
	if !strings.HasPrefix(arg, "!") {
		arg = "!" + arg
	}
	c = ParseCommand(arg)
	if c == CommandUnknown {
		return c, ErrUnknownCommand
	}
	return c, nil
}

func (m *Message) MoneyArgument() (int, error) {
	if m.Command != CommandTransfer &&
		m.Command != CommandDice &&
		m.Command != CommandDeposit &&
		m.Command != CommandWithdraw &&
		m.Command != CommandDebt &&
		m.Command != CommandRepay {
		return 0, ErrWrongCommand
	}
	arg := m.Argument()
	n, err := strconv.ParseInt(arg, 10, 32)
	if err != nil {
		return 0, ErrSpecifyAmount
	}
	if n <= 0 {
		return 0, ErrNotPositive
	}
	return int(n), nil
}
