package teleutil

import (
	"nechego/format"
	"regexp"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func Name(m *tele.ChatMember) string {
	name := m.Title
	if name == "" {
		return strings.TrimSpace(m.User.FirstName + " " + m.User.LastName)
	}
	return name
}

func Mention(c tele.Context, m *tele.ChatMember) string {
	return format.Mention(c.Chat().ID, Name(m))
}

func Args(c tele.Context, re *regexp.Regexp) []string {
	return re.FindStringSubmatch(c.Message().Text)
}

func Member(c tele.Context, user tele.Recipient) *tele.ChatMember {
	m, err := c.Bot().ChatMemberOf(c.Chat(), user)
	if err != nil {
		panic("can't get chat member")
	}
	return m
}

func Promote(c tele.Context, m *tele.ChatMember) error {
	if Admin(m) {
		return nil
	}
	m.Rights.CanBeEdited = true
	m.Rights.CanManageChat = true
	return c.Bot().Promote(c.Chat(), m)
}

func Admin(m *tele.ChatMember) bool {
	return m.Role == tele.Administrator || m.Role == tele.Creator
}

func NumArg(c tele.Context, re *regexp.Regexp, n int) []int {
	s := Args(c, re)[n]
	nums := []int{}
	for _, x := range strings.Fields(s) {
		n, err := strconv.Atoi(x)
		if err != nil {
			continue
		}
		nums = append(nums, n)
	}
	return nums
}
