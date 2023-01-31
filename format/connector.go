package format

import (
	"strings"
)

type Connector struct {
	c string
	b strings.Builder
}

func NewConnector(c string) *Connector {
	return &Connector{c: c}
}

func (l *Connector) Add(s string) {
	if l.b.Len() > 0 {
		l.b.WriteString(l.c)
	}
	l.b.WriteString(s)
}

func (l *Connector) String() string {
	return l.b.String()
}
