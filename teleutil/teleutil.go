package teleutil

import (
	"fmt"
	"html"
	"nechego/game"
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

func Mention(c tele.Context, user any) string {
	var member *tele.ChatMember
	switch x := user.(type) {
	case *tele.ChatMember:
		member = x
	case tele.Recipient:
		member = Member(c, x)
	case int64:
		member = Member(c, tele.ChatID(x))
	case *game.User:
		member = Member(c, tele.ChatID(x.TUID))
	default:
		panic(fmt.Sprintf("unexpected type %T", x))
	}
	const format = `<a href="tg://user?id=%d">%s</a>`
	return fmt.Sprintf(format, member.User.ID, html.EscapeString(Name(member)))
}

func Args(c tele.Context, re *regexp.Regexp) []string {
	return re.FindStringSubmatch(c.Text())
}

func Member(c tele.Context, user tele.Recipient) *tele.ChatMember {
	m, err := c.Bot().ChatMemberOf(c.Chat(), user)
	if err != nil {
		panic("cannot get chat member")
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

func Left(m *tele.ChatMember) bool {
	return m.Role == tele.Kicked || m.Role == tele.Left
}

func NumArgAll(c tele.Context, re *regexp.Regexp, n int) []int {
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

func NumArg(c tele.Context, re *regexp.Regexp, n int) []int {
	const limit = 7
	nums := NumArgAll(c, re, n)
	if len(nums) > limit {
		nums = nums[:limit]
	}
	return nums
}

func Reply(c tele.Context) (u *tele.User, ok bool) {
	if !c.Message().IsReply() || c.Message().ReplyTo.Sender.IsBot {
		return nil, false
	}
	return c.Message().ReplyTo.Sender, true
}

func Lock(c tele.Context, u *game.Universe) (*game.World, *game.User) {
	world, err := u.World(c.Chat().ID)
	if err != nil {
		panic(fmt.Sprintf("cannot get world: %s", err))
	}
	world.Lock()
	user := world.UserByID(c.Sender().ID)
	return world, user
}
