package app

import (
	"bytes"
	"fmt"
	"math/rand"
	"nechego/input"
	"nechego/model"
	"nechego/numbers"
	"nechego/pets"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

func formatWeight(n float64) string {
	return fmt.Sprintf("<code>%.2f кг ⚖️</code>", n)
}

func formatCommand(c input.Command) string {
	return fmt.Sprintf("<code>%s</code>", c)
}

func formatTitles(s ...string) string {
	if len(s) > 0 {
		s[0] = strings.Title(s[0])
	}
	titles := joinSpace(s...)
	if titles == "" {
		titles = "Пользователь"
	}
	return titles
}

func formatIcons(icon ...string) string {
	icons := strings.Join(icon, "·")
	return fmt.Sprintf("<code>%s</code>", icons)
}

func formatEnergyRemaining(n int) string {
	return fmt.Sprintf("<i>Энергии осталось: %s</i>", formatEnergy(n))
}

func formatEnergyCooldown(d time.Duration) string {
	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	return fmt.Sprintf("⏰ До восстановления энергии: <code>%d минут %d секунд</code>.",
		mins, secs)
}

const maxTopNumber = 5

// clampTopNumber returns x if x < maxTopNumber; otherwise returns maxTopNumber.
func clampTopNumber(x int) int {
	return numbers.Min(x, maxTopNumber)
}

var (
	emojisActive   = []string{"🔈", "🔔", "✅", "🆗", "▶️"}
	emojisInactive = []string{"🔇", "🔕", "💤", "❌", "⛔️", "🚫", "⏹"}
)

func activeEmoji() string {
	return emojisActive[rand.Intn(len(emojisActive))]
}

func inactiveEmoji() string {
	return emojisInactive[rand.Intn(len(emojisInactive))]
}
