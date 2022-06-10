package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	catURL        = "https://thiscatdoesnotexist.com/"
	animeFormat   = "https://thisanimedoesnotexist.ai/results/psi-%s/seed%s.png"
	furFormat     = "https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%s.jpg"
	flagFormat    = "https://thisflagdoesnotexist.com/images/%d.png"
	personURL     = "https://thispersondoesnotexist.com/image"
	horseURL      = "https://thishorsedoesnotexist.com/"
	artURL        = "https://thisartworkdoesnotexist.com/"
	carURL        = "https://www.thisautomobiledoesnotexist.com/"
	weatherURL    = "https://wttr.in/"
	weatherFormat = `?format=%l:+%c+%t+\nОщущается+как+%f\n\nВетер+—+%w\nВлажность+—+%h\nДавление+—+%P\nФаза+луны+—+%m\nУФ-индекс+—+%u\n`
)

var infaRe = regexp.MustCompile("^!инфа?(.*)")

// handleProbability responds with the probability of the message.
func (a *app) handleProbability(c tele.Context, m *message) error {
	ss := infaRe.FindStringSubmatch(m.text)
	if len(ss) < 2 {
		return nil
	}
	s := ss[1]
	return c.Send(probability(s))
}

// handleWho responds with the message appended to the random chat member.
func (a *app) handleWho(c tele.Context, m *message) error {
	userID, err := a.getRandomGroupMember(c.Chat().ID)
	if err != nil {
		return err
	}

	chat, err := c.Bot().ChatByID(userID)
	if err != nil {
		return err
	}

	name := getUserNameEscaped(chat)
	text := m.argumentEscaped()
	return c.Send(who(userID, name, text), tele.ModeMarkdownV2)
}

// handleCat sends a picture of a cat.
func (a *app) handleCat(c tele.Context) error {
	pic, err := fetchPicture(catURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

// handleTitle sets the admin title of the sender.
func (a *app) handleTitle(c tele.Context, m *message) error {
	title := m.argument()
	if len(title) > 16 {
		return c.Send("Ошибка: максимальная длина имени 16 символов")
	}
	if err := c.Bot().SetAdminTitle(c.Chat(), c.Sender(), title); err != nil {
		return c.Send("Ошибка")
	}
	return nil
}

var animePsis = []string{"0.3", "0.4", "0.5", "0.6", "0.7", "0.8", "0.9", "1.0",
	"1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "2.0"}

// handleAnime sends an anime picture.
func (a *app) handleAnime(c tele.Context) error {
	psi := animePsis[rand.Intn(len(animePsis))]
	seed := getRandomNumbers(5)
	url := fmt.Sprintf(animeFormat, psi, seed)

	pic, err := fetchPicture(url)
	if err != nil {
		return err
	}

	return c.Send(pic)
}

// handleFurry sends a furry picture.
func (a *app) handleFurry(c tele.Context) error {
	seed := getRandomNumbers(5)
	url := fmt.Sprintf(furFormat, seed)

	pic, err := fetchPicture(url)
	if err != nil {
		return err
	}

	return c.Send(pic)
}

// handleFlag sends a picture of a flag.
func (a *app) handleFlag(c tele.Context) error {
	seed := rand.Intn(5000)
	url := fmt.Sprintf(flagFormat, seed)

	pic, err := fetchPicture(url)
	if err != nil {
		return err
	}

	return c.Send(pic)
}

// handlePerson sends a picture of a person.
func (a *app) handlePerson(c tele.Context) error {
	pic, err := fetchPicture(personURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

// handleHorse sends a picture of a horse.
func (a *app) handleHorse(c tele.Context) error {
	pic, err := fetchPicture(horseURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

// handleArt sends a picture of an art.
func (a *app) handleArt(c tele.Context) error {
	pic, err := fetchPicture(artURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

var carImageRe = regexp.MustCompile(
	"<img id = \"vehicle\" src=\"data:image/png;base64,(.+)\" class=\"center\">")

// handleCar sends a picture of a car.
func (a *app) handleCar(c tele.Context) error {
	r, err := http.Get(carURL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	ss := carImageRe.FindStringSubmatch(string(data))
	if len(ss) < 2 {
		return nil
	}
	b64 := ss[1]
	img, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	return c.Send(dataToPhoto(img))
}

const pairOfTheDayFormat = "Пара дня ✨\n%s 💘 %s"

// handlePair sends the current pair of the day, randomly choosing a new pair if
// needed.
func (a *app) handlePair(c tele.Context) error {
	groupID := c.Chat().ID

	pair, err := a.store.getPair(groupID)
	if errors.Is(err, errNoPair) {
		x, err := a.getRandomGroupMember(groupID)
		if err != nil {
			return err
		}
		y, err := a.getRandomGroupMember(groupID)
		if err != nil {
			return err
		}

		if x == y {
			return c.Send("💔")
		}

		pair = pairOfTheDay{x, y}
		if err := a.store.insertPair(groupID, pair); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	chatX, err := c.Bot().ChatByID(pair.userIDx)
	if err != nil {
		return err
	}
	chatY, err := c.Bot().ChatByID(pair.userIDy)
	if err != nil {
		return err
	}

	return c.Send(fmt.Sprintf(pairOfTheDayFormat,
		mention(pair.userIDx, getUserNameEscaped(chatX)),
		mention(pair.userIDy, getUserNameEscaped(chatY))),
		tele.ModeMarkdownV2)
}

const eblanOfTheDayFormat = "Еблан дня: %s 😸"

// handleEblan sends the current eblan of the day, randomly choosing a new one
// if needed.
func (a *app) handleEblan(c tele.Context) error {
	groupID := c.Chat().ID

	userID, err := a.store.getEblan(groupID)
	if errors.Is(err, errNoEblan) {
		userID, err = a.getRandomGroupMember(groupID)
		if err != nil {
			return err
		}

		if err := a.store.insertEblan(groupID, userID); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	chat, err := c.Bot().ChatByID(userID)
	if err != nil {
		return err
	}

	eblan := getUserNameEscaped(chat)
	return c.Send(fmt.Sprintf(eblanOfTheDayFormat, mention(userID, eblan)),
		tele.ModeMarkdownV2)
}

const masyunyaStickersName = "masyunya_vk"

// handleMasyunya sends a random sticker of Masyunya.
func (a *app) handleMasyunya(c tele.Context) error {
	ss, err := c.Bot().StickerSet(masyunyaStickersName)
	if err != nil {
		return err
	}
	s := ss.Stickers[rand.Intn(len(ss.Stickers))]
	return c.Send(&s)
}

const helloChance = 0.2

// handleHello sends a hello sticker
func (a *app) handleHello(c tele.Context) error {
	n := rand.Float64()
	if n <= helloChance {
		s := helloStickers[rand.Intn(len(helloStickers))]
		return c.Send(s)
	}
	return nil
}

var mouseVideo = &tele.Video{File: tele.FromDisk("data/mouse.mp4")}

// handleMouse sends the mouse video
func (a *app) handleMouse(c tele.Context) error {
	return c.Send(mouseVideo)
}

const weatherTimeout = 8 * time.Second

// handleWeather sends the current weather for a given city
func (a *app) handleWeather(c tele.Context, m *message) error {
	ss := weatherRe.FindStringSubmatch(m.text)
	if len(ss) < 2 {
		return nil
	}
	loc := ss[1]

	ctx, cancel := context.WithTimeout(context.Background(), weatherTimeout)
	defer cancel()

	l := weatherURL + loc + weatherFormat
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, l, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept-Language", "ru")

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		if err.(*url.Error).Timeout() {
			return c.Send("Ошибка: время запроса вышло ☔️")
		}
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		if r.StatusCode == http.StatusNotFound {
			return c.Send("Ошибка: такого места не существует ☔️")
		}
		return c.Send("Ошибка: неудачный запрос ☔️")
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return c.Send(string(data))
}

// handleKeyboardOpen opens the keyboard.
func (a *app) handleKeyboardOpen(c tele.Context) error {
	return c.Send("Клавиатура ⌨️", a.keyboard)
}

// handleKeyboardClose closes the keyboard.
func (a *app) handleKeyboardClose(c tele.Context) error {
	return c.Send("Клавиатура закрыта 😣", tele.RemoveKeyboard)
}

var (
	emojisActive   = []string{"🔈", "🔔", "✅", "🆗", "▶️"}
	emojisInactive = []string{"🔇", "🔕", "💤", "❌", "⛔️", "🚫", "⏹"}
)

// handleTurnOn turns the bot on.
func (a *app) handleTurnOn(c tele.Context) error {
	emoji := emojisActive[rand.Intn(len(emojisActive))]
	groupID := c.Chat().ID
	a.status.turnOnLocal(groupID)
	return c.Send(fmt.Sprintf("Бот включен %s", emoji))
}

// handleTurnOff turns the bot off.
func (a *app) handleTurnOff(c tele.Context) error {
	emoji := emojisInactive[rand.Intn(len(emojisInactive))]
	groupID := c.Chat().ID
	a.status.turnOffLocal(groupID)
	return c.Send(fmt.Sprintf("Бот выключен %s", emoji), tele.RemoveKeyboard)
}

const accessRestricted = "Доступ ограничен 🔒"
const userBlocked = "Пользователь заблокирован 🚫"
const userAlreadyBlocked = "Пользователь уже заблокирован 🛑"
const userUnblocked = "Пользователь разблокирован ✅"
const userAlreadyUnblocked = "Пользователь не заблокирован ❎"

// handleBan adds the user ID of the reply message's sender to the ban list.
func (a *app) handleBan(c tele.Context) error {
	if !a.owners.check(c.Sender().ID) {
		return c.Send(accessRestricted)
	}
	if !c.Message().IsReply() {
		return nil
	}
	id := c.Message().ReplyTo.Sender.ID
	_, exist := a.bans.LoadOrStore(id, struct{}{})
	if exist {
		return c.Send(userAlreadyBlocked)
	}
	return c.Send(userBlocked)
}

// handleUnban removes the user ID of the reply message's sender from the ban list.
func (a *app) handleUnban(c tele.Context) error {
	if !a.owners.check(c.Sender().ID) {
		return c.Send(accessRestricted)
	}
	if !c.Message().IsReply() {
		return nil
	}
	id := c.Message().ReplyTo.Sender.ID
	_, exist := a.bans.LoadAndDelete(id)
	if !exist {
		c.Send(userAlreadyUnblocked)
	}
	return c.Send(userUnblocked)
}

const listTemplate = `📝 *Информация* 📌

👤 _Администрация_
%s
🛑 _Черный список_
%s`

// handleList sends a list of useful information.
func (a *app) handleList(c tele.Context) error {
	var admins string
	l := a.owners.list()
	fmt.Println(l)
	if len(l) == 0 {
		admins = "…\n"
	} else {
		for _, id := range l {
			user, err := c.Bot().ChatByID(id)
			if err != nil {
				continue
			}
			member, err := c.Bot().ChatMemberOf(c.Chat(), user)
			if err != nil {
				continue
			}
			if member.Role == tele.Left {
				continue
			}
			admins += "— " + mention(id, getUserNameEscaped(user)) + "\n"
		}
	}

	var banned string
	a.bans.Range(func(key, value any) bool {
		id := key.(int64)
		user, err := c.Bot().ChatByID(id)
		if err != nil {
			return true
		}
		member, err := c.Bot().ChatMemberOf(c.Chat(), user)
		if err != nil {
			return true
		}
		if member.Role == tele.Left {
			return true
		}
		banned += "— " + mention(id, getUserNameEscaped(user)) + "\n"
		return true
	})
	if banned == "" {
		banned = "…\n"
	}
	list := fmt.Sprintf(listTemplate, admins, banned)
	return c.Send(list, tele.ModeMarkdownV2)
}

// getRandomGroupMember returns the ID of the random group member.
func (a *app) getRandomGroupMember(groupID int64) (int64, error) {
	userIDs, err := a.store.getUserIDs(groupID)
	if err != nil {
		return 0, err
	}
	return userIDs[rand.Intn(len(userIDs))], nil
}

// getRandomNumbers returns a string of random numbers of length c.
func getRandomNumbers(c int) string {
	var nums []string
	for i := 0; i < c; i++ {
		n := rand.Intn(10)
		nums = append(nums, fmt.Sprintf("%d", n))
	}
	return strings.Join(nums, "")
}

// getUserName returns the displayed user name.
func getUserName(chat *tele.Chat) string {
	return strings.TrimSpace(chat.FirstName + " " + chat.LastName)
}

// getUserNameEscaped returns the displayed user name and escapes it for Markdown.
func getUserNameEscaped(chat *tele.Chat) string {
	return markdownEscaper.Replace(getUserName(chat))
}

// probabilityTemplates regexp: "^.*%s.*%d%%\"".
var probabilityTemplates = []string{
	"Здравый смысл говорит мне о том, что %s с вероятностью %d%%",
	"Благодаря чувственному опыту я определил, что %s с вероятностью %d%%",
	"Я думаю, что %s с вероятностью %d%%",
	"Используя диалектическую логику, я пришел к выводу, что %s с вероятностью %d%%",
	"Проведя некие изыскания, я высяснил, что %s с вероятностью %d%%",
	"Я провел мысленный экперимент и выяснил, что %s с вероятностью %d%%",
	"Мои интеллектуальные потуги привели меня к тому, что %s с вероятностью %d%%",
	"С помощью фактов и логики я доказал, что %s с вероятностью %d%%",
	"Как показывает практика, %s с вероятностью %d%%",
	"Прикинув раз на раз, я определился с тем, что %s с вероятностью %d%%",
	"Уверяю вас в том, что %s с вероятностью %d%%",
}

// probability returns the probability of the message.
func probability(message string) string {
	t := probabilityTemplates[rand.Intn(len(probabilityTemplates))]
	p := rand.Intn(101)
	return fmt.Sprintf(t, message, p)
}

// who returns the mention of the user prepended to the message.
func who(userID int64, name, message string) string {
	return fmt.Sprintf("%s %s", mention(userID, name), message)
}

// mention returns the mention of the user by the name.
func mention(userID int64, name string) string {
	return fmt.Sprintf("[%s](tg://user?id=%d)", name, userID)
}

// fetchPicture returns a picture located at the specified URL.
func fetchPicture(url string) (*tele.Photo, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return dataToPhoto(body), nil
}

// dataToPhoto converts the image data to Photo.
func dataToPhoto(data []byte) *tele.Photo {
	return &tele.Photo{File: tele.FromReader(bytes.NewReader(data))}
}

var markdownEscaper = newMarkdownEscaper()

// newMarkdownEscaper creates a new Markdown replacer. The replacer
// escapes any character with the code between 1 and 126 inclusively
// with a preceding backslash.
func newMarkdownEscaper() *strings.Replacer {
	var table []string
	for i := 1; i <= 126; i++ {
		c := string(i)
		table = append(table, c, "\\"+c)
	}
	return strings.NewReplacer(table...)
}
