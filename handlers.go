package main

import (
	"bytes"
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

var emojisActive = []string{"🔈", "🔔", "✅", "🆗", "▶️"}
var emojisInactive = []string{"🔇", "🔕", "💤", "❌", "⛔️", "🚫", "⏹"}

const catURL = "https://thiscatdoesnotexist.com/"
const animeFormat = "https://thisanimedoesnotexist.ai/results/psi-%s/seed%s.png"
const furFormat = "https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%s.jpg"
const flagFormat = "https://thisflagdoesnotexist.com/images/%d.png"
const personURL = "https://thispersondoesnotexist.com/image"
const horseURL = "https://thishorsedoesnotexist.com/"
const artURL = "https://thisartworkdoesnotexist.com/"
const carURL = "https://www.thisautomobiledoesnotexist.com/"
const weatherFormat = "https://wttr.in/%s?format=3"

var animePsis = []string{"0.3", "0.4", "0.5", "0.6", "0.7", "0.8", "0.9", "1.0",
	"1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "2.0"}
var carImageRe = regexp.MustCompile(
	"<img id = \"vehicle\" src=\"data:image/png;base64,(.+)\" class=\"center\">")
var infaRe = regexp.MustCompile("^!инфа?(.*)")

const helloChance = 0.2

var markdownEscaper = newMarkdownEscaper()

var mouseVideo = &tele.Video{File: tele.FromDisk("mouse.mp4")}

// handleProbability responds with the probability of the message.
func (a *app) handleProbability(c tele.Context, m *message) error {
	s := infaRe.FindStringSubmatch(m.text)[1]
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

	b64img := carImageRe.FindStringSubmatch(string(data))[1]
	img, err := base64.StdEncoding.DecodeString(b64img)
	if err != nil {
		return err
	}
	return c.Send(byteSliceToPhoto(img))
}

// handlePair sends the current pair of the day, randomly choosing a new pair if
// needed.
func (a *app) handlePair(c tele.Context) error {
	groupID := c.Chat().ID

	p, err := a.store.getPair(groupID)
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

		p = pair{x, y}
		if err := a.store.insertPair(groupID, p); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	chatX, err := c.Bot().ChatByID(p.x)
	if err != nil {
		return err
	}
	chatY, err := c.Bot().ChatByID(p.y)
	if err != nil {
		return err
	}

	return c.Send(fmt.Sprintf("Пара дня ✨\n%s 💘 %s",
		mention(p.x, getUserNameEscaped(chatX)),
		mention(p.y, getUserNameEscaped(chatY))),
		tele.ModeMarkdownV2)
}

// handleEblan sends the current eblan of the day, randomly choosing a new one
// if needed.
func (a *app) handleEblan(c tele.Context) error {
	groupID := c.Chat().ID

	userID, err := a.store.getEblan(groupID)
	if errors.Is(err, errNoEblan) {
		e, err := a.getRandomGroupMember(groupID)
		if err != nil {
			return err
		}

		userID = e
		if err := a.store.insertEblan(groupID, e); err != nil {
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
	s := mention(userID, eblan)
	return c.Send(fmt.Sprintf("Еблан дня: %s 😸", s), tele.ModeMarkdownV2)
}

// handleMasyunya sends a random sticker of Masyunya.
func (a *app) handleMasyunya(c tele.Context) error {
	name := "masyunya_vk"
	ss, err := c.Bot().StickerSet(name)
	if err != nil {
		return err
	}
	s := ss.Stickers[rand.Intn(len(ss.Stickers))]
	return c.Send(&s)
}

// handleHello sends a hello sticker
func (a *app) handleHello(c tele.Context) error {
	n := rand.Float64()
	if n <= helloChance {
		s := helloStickers[rand.Intn(len(helloStickers))]
		return c.Send(s)
	}
	return nil
}

// handleMouse sends the mouse video
func (a *app) handleMouse(c tele.Context) error {
	return c.Send(mouseVideo)
}

// handleWeather sends the current weather for a given city
func (a *app) handleWeather(c tele.Context, m *message) error {
	client := &http.Client{Timeout: 3 * time.Second}

	r, err := client.Get(fmt.Sprintf(weatherFormat, m.argument()))
	if err != nil {
		if err.(*url.Error).Timeout() {
			return c.Send("Ошибка: время запроса вышло ☔️")
		}
		return err
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusNotFound {
		return c.Send("Ошибка: такого места не существует ☔️")
	} else if r.StatusCode != http.StatusOK {
		return c.Send("Ошибка: неудачный запрос ☔️")
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return c.Send(string(data))
}

func (a *app) handleKeyboardOpen(c tele.Context) error {
	return c.Send("Клавиатура ⌨️", a.keyboard)
}

func (a *app) handleKeyboardClose(c tele.Context) error {
	return c.Send("Клавиатура отключена 😣", tele.RemoveKeyboard)
}

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
	nums := []string{}
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

// getUserNameEscaped returns the displayed user name and escapes it for
// Markdown.
func getUserNameEscaped(chat *tele.Chat) string {
	return markdownEscaper.Replace(getUserName(chat))
}

// probability returns the probability of the message.
func probability(message string) string {
	p := rand.Intn(101)
	t := probabilityTemplates[rand.Intn(len(probabilityTemplates))]
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
	return byteSliceToPhoto(body), nil
}

// byteSliceToPhoto converts the byte slice of image data to Photo.
func byteSliceToPhoto(data []byte) *tele.Photo {
	return &tele.Photo{File: tele.FromReader(bytes.NewReader(data))}
}

// newMarkdownEscaper creates a new Markdown replacer. The replacer escapes any
// character with code between 1 and 126 inclusively with a preceding '\'.
func newMarkdownEscaper() *strings.Replacer {
	var table []string
	for i := 1; i <= 126; i++ {
		c := string(i)
		table = append(table, c, "\\"+c)
	}
	return strings.NewReplacer(table...)
}
