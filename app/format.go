package app

import (
	"bytes"
	"math/rand"
	"strings"

	tele "gopkg.in/telebot.v3"
)

// photoFromBytes converts the image data to Photo.
func photoFromBytes(data []byte) *tele.Photo {
	return &tele.Photo{File: tele.FromReader(bytes.NewReader(data))}
}

// markdownEscaper escapes any character with the code between 1 and 126
// inclusively with a preceding backslash.
var markdownEscaper = func() *strings.Replacer {
	var table []string
	for i := 1; i <= 126; i++ {
		c := string(rune(i))
		table = append(table, c, "\\"+c)
	}
	return strings.NewReplacer(table...)
}()

var errorSigns = []string{"❌", "🚫", "⭕️", "🛑", "⛔️", "📛", "💢", "❗️", "‼️", "⚠️"}

// errorSign returns a random error sign.
func errorSign() string {
	return errorSigns[rand.Intn(len(errorSigns))]
}

// makeError formats the error message.
func makeError(s string) string {
	return errorSign() + " " + s
}

// randInRange returns a random value in range [min, max].
func randInRange(min, max int) int {
	return min + rand.Intn(max-min+1)
}
