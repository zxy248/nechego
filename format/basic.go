package format

import (
	"fmt"
	"strings"
)

func Code(s string) string {
	return "<code>" + s + "</code>"
}

func Bold(s string) string {
	return "<b>" + s + "</b>"
}

func Italic(s string) string {
	return "<i>" + s + "</i>"
}

func Lines(s ...string) string {
	return strings.Join(s, "\n")
}

func Words(s ...string) string {
	return strings.Join(s, " ")
}

func Value(a any) string {
	return fmt.Sprint(a)
}
