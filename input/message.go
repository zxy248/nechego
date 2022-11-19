package input

import (
	"errors"
	"html"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
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

func (m *Message) RawArgument() string {
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
	_, arg, _ := strings.Cut(m.Raw, " ")
	return arg
}

// Argument returns an argument probably contained in the message.
func (m *Message) Argument() string {
	return Sanitize(m.RawArgument())
}

var (
	noArg = errors.New("this command has no argument")
)

// TopArgument represents an argument of the CommandTop.
type TopArgument struct {
	Number *int
	String string
}

func (m *Message) TopArgument() (TopArgument, error) {
	if m.Command != CommandTop {
		panic(noArg)
	}
	match := topRe.FindStringSubmatch(m.Raw)
	s := Sanitize(match[2])
	maybe, err := strconv.ParseInt(match[1], 10, 32)
	if err != nil {
		return TopArgument{String: s}, nil
	}
	n := int(maybe)
	return TopArgument{Number: &n, String: s}, nil
}

var commandCommands = []Command{
	CommandForbid,
	CommandPermit,
}

func hasCommandArgument(c Command) bool {
	return slices.Contains(commandCommands, c)
}

var (
	ErrNoCommand      = errors.New("no command")
	ErrUnknownCommand = errors.New("unknown command")
)

func (m *Message) CommandActionArgument() (Command, error) {
	c := CommandUnknown
	if !hasCommandArgument(m.Command) {
		panic(noArg)
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

var moneyCommands = []Command{
	CommandTransfer,
	CommandDice,
	CommandDeposit,
	CommandWithdraw,
	CommandDebt,
	CommandRepay,
}

func hasMoneyArgument(c Command) bool {
	return slices.Contains(moneyCommands, c)
}

var (
	ErrSpecifyAmount = errors.New("specify amount")
	ErrAllIn         = errors.New("all in")
)

func (m *Message) MoneyArgument() (int, error) {
	if !hasMoneyArgument(m.Command) {
		panic(noArg)
	}
	arg := m.Argument()
	if arg == "все" || arg == "всё" {
		return 0, ErrAllIn
	}
	n, err := strconv.ParseInt(arg, 10, 32)
	if err != nil || n <= 0 {
		return 0, ErrSpecifyAmount
	}
	return int(n), nil
}

func Sanitize(s string) string {
	return html.EscapeString(s)
}
