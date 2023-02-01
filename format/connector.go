package format

import (
	"strings"
)

// Connector joins added strings.
type Connector struct {
	sep string
	b   strings.Builder
}

// NewConnector returns a Connector that joins added strings with sep.
func NewConnector(sep string) *Connector {
	return &Connector{sep: sep}
}

// Add adds a string to the connector.
func (l *Connector) Add(s string) {
	if l.b.Len() > 0 {
		l.b.WriteString(l.sep)
	}
	l.b.WriteString(s)
}

// String returns the added strings joined by the separator.
func (l *Connector) String() string {
	return l.b.String()
}
