package input

import (
	"strings"
)

// Message represents an input Message.
type Message struct {
	Raw     string
	Command Command
}

// Parse returns a new message.
func Parse(s string) *Message {
	return &Message{s, recognizeCommand(s)}
}

// Argument returns the Argument probably contained in the message.
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
