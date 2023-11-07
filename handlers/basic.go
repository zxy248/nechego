package handlers

import (
	"errors"
	"fmt"
	"html"
	"io"
	"math/rand"
	"nechego/avatar"
	"nechego/format"
	"nechego/game"
	"nechego/handlers/parse"
	tu "nechego/teleutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Help struct{}

var helpRe = Regexp("^!(помощь|команды|документ)")

func (h *Help) Match(s string) bool {
	return helpRe.MatchString(s)
}

func (h *Help) Handle(c tele.Context) error {
	return c.Send("📖 <b>Документация:</b> nechego.pages.dev.", tele.ModeHTML)
}

type Infa struct{}

func (h *Infa) Match(s string) bool {
	_, ok := infaCommand(s)
	return ok
}

func (h *Infa) Handle(c tele.Context) error {
	s, _ := infaCommand(c.Text())
	return c.Send(infa(s, rand.Intn(101)))
}

func infaCommand(s string) (text string, ok bool) {
	re := Regexp("^!инфа ?(.*)")
	m := re.FindStringSubmatch(s)
	if m == nil {
		return "", false
	}
	return m[1], true
}

func infa(s string, p int) string {
	pres := []string{
		"Здравый смысл говорит мне о том, что",
		"Благодаря чувственному опыту я определил, что",
		"Я думаю, что",
		"Используя диалектическую логику, я пришел к выводу, что",
		"Проведя некие изыскания, я выяснил, что",
		"Я провёл мысленный эксперимент и выяснил, что",
		"Мои интеллектуальные потуги привели меня к тому, что",
		"С помощью фактов и логики я доказал, что",
		"Как показывает практика,",
		"Прикинув раз на раз, я определился с тем, что",
		"Уверяю вас в том, что",
	}
	pre := pres[rand.Intn(len(pres))]
	suf := "с вероятностью"
	prob := fmt.Sprintf("%d%%", p)
	return strings.Join([]string{pre, s, suf, prob}, " ")
}

type Who struct {
	Universe *game.Universe
}

func (h *Who) Match(s string) bool {
	_, ok := whoCommand(s)
	return ok
}

func (h *Who) Handle(c tele.Context) error {
	text, _ := whoCommand(c.Text())
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	m := tu.Link(c, world.RandomUserID())
	s := html.EscapeString(text)
	return c.Send(m+" "+s, tele.ModeHTML)
}

func whoCommand(s string) (text string, ok bool) {
	return textCommand(parse.Prefix("!кто"), s)
}

type List struct {
	Universe *game.Universe
}

var listRe = Regexp("^!список ?(.*)")

func (h *List) Match(s string) bool {
	return listRe.MatchString(s)
}

func (h *List) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	us := world.RandomUserIDs(3 + rand.Intn(3))
	arg := tu.Args(c, listRe)[1]
	s := []string{fmt.Sprintf("<b>📝 Список %s</b>", arg)}
	for _, u := range us {
		who := tu.Link(c, u)
		s = append(s, fmt.Sprintf("<b>•</b> %s", who))
	}
	return c.Send(strings.Join(s, "\n"), tele.ModeHTML)
}

type Top struct {
	Universe *game.Universe
}

func (h *Top) Match(s string) bool {
	_, ok := topCommand(s)
	return ok
}

func (h *Top) Handle(c tele.Context) error {
	text, _ := topCommand(c.Text())
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	us := world.RandomUserIDs(3 + rand.Intn(3))
	s := []string{fmt.Sprintf("<b>🏆 Топ %s</b>", text)}
	for i, u := range us {
		s = append(s, fmt.Sprintf("<i>%d.</i> %s", 1+i, tu.Link(c, u)))
	}
	return c.Send(strings.Join(s, "\n"), tele.ModeHTML)
}

func topCommand(s string) (text string, ok bool) {
	return textCommand(parse.Match("!топ"), s)
}

type Game struct{}

var gameRe = Regexp("^!игр")

func (h *Game) Match(s string) bool {
	return gameRe.MatchString(s)
}

func (h *Game) Handle(c tele.Context) error {
	games := [...]*tele.Dice{tele.Dart, tele.Ball, tele.Goal, tele.Slot, tele.Bowl}
	return c.Send(games[rand.Intn(len(games))])
}

type Weather struct{}

func (h *Weather) Match(s string) bool {
	_, ok := weatherCommand(s)
	return ok
}

func (h *Weather) Handle(c tele.Context) error {
	const addr = "https://wttr.in/"
	const format = `?format=%l:+%c+%t+\n` +
		`Ощущается+как+%f\n\n` +
		`Ветер+—+%w\n` +
		`Влажность+—+%h\n` +
		`Давление+—+%P\n` +
		`Фаза+луны+—+%m\n` +
		`УФ-индекс+—+%u\n`
	city, _ := weatherCommand(c.Text())
	city = url.PathEscape(city)

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest(http.MethodGet, addr+city+format, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept-Language", "ru")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return c.Send("☔️ Такого места не существует.")
	} else if resp.StatusCode != http.StatusOK {
		return c.Send("☔️ Неудачный запрос.")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return c.Send(string(data))
}

func weatherCommand(s string) (city string, ok bool) {
	return textCommand(parse.Prefix("!погода"), s)
}

type Avatar struct {
	Universe *game.Universe
	Avatars  *avatar.Storage
}

var avatarRe = Regexp("^!ава")

func (h *Avatar) Match(s string) bool {
	return avatarRe.MatchString(s)
}

func (h *Avatar) Handle(c tele.Context) error {
	target := c.Sender().ID
	photo := c.Message().Photo
	if reply, ok := tu.Reply(c); ok {
		// If the user has admin rights, they can set an
		// avatar for other users.
		world, user := tu.Lock(c, h.Universe)
		admin := user.Admin()
		world.Unlock()
		if !admin {
			return c.Send("📷 Нельзя установить аватар другому пользователю.")
		}
		target = reply.ID
	}

	if photo == nil {
		if avatar, ok := h.Avatars.Get(target); ok {
			return c.Send(avatar)
		}
		return c.Send("📷 Прикрепите изображение.")
	}
	if err := h.Avatars.Set(target, photo); errors.Is(err, avatar.ErrSize) {
		return c.Send("📷 Максимальный размер аватара %dx%d пикселей.",
			h.Avatars.MaxWidth, h.Avatars.MaxHeight)
	} else if err != nil {
		return err
	}
	return c.Send("📸 Аватар установлен.")
}

type Ban struct {
	Universe   *game.Universe
	DurationHr int // Ban duration in hours.
}

var banRe = Regexp("^!бан")

func (h *Ban) Match(s string) bool {
	return banRe.MatchString(s)
}

func (h *Ban) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if !user.Admin() {
		return c.Send(format.AdminsOnly)
	}
	reply, ok := tu.Reply(c)
	if !ok {
		return c.Send(format.RepostMessage)
	}
	target := world.User(reply.ID)
	duration := time.Hour * time.Duration(h.DurationHr)
	target.BannedUntil = time.Now().Add(duration)
	return c.Send(format.UserBanned(h.DurationHr))
}

type Unban struct {
	Universe *game.Universe
}

var unbanRe = Regexp("^!разбан")

func (h *Unban) Match(s string) bool {
	return unbanRe.MatchString(s)
}

func (h *Unban) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if !user.Admin() {
		return c.Send(format.AdminsOnly)
	}
	reply, ok := tu.Reply(c)
	if !ok {
		return c.Send(format.RepostMessage)
	}
	world.User(reply.ID).BannedUntil = time.Time{}
	return c.Send(format.UserUnbanned)
}
