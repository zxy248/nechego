package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"

	tele "gopkg.in/telebot.v3"
)

// probabilityTemplates regexp: "^.*%s.*%d%%\""
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

const catURL = "https://thiscatdoesnotexist.com/"
const animeFormat = "https://thisanimedoesnotexist.ai/results/psi-%s/seed%s.png"
const furFormat = "https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%s.jpg"
const flagFormat = "https://thisflagdoesnotexist.com/images/%d.png"
const personURL = "https://thispersondoesnotexist.com/image"
const horseURL = "https://thishorsedoesnotexist.com/"
const artURL = "https://thisartworkdoesnotexist.com/"
const carURL = "https://www.thisautomobiledoesnotexist.com/"

var animePsis = []string{"0.3", "0.4", "0.5", "0.6", "0.7", "0.8", "0.9", "1.0",
	"1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "2.0"}
var carImageRe = regexp.MustCompile(
	"<img id = \"vehicle\" src=\"data:image/png;base64,(.+)\" class=\"center\">")

// handleProbability responds with the probability of the message
func (a *app) handleProbability(c tele.Context, message string) error {
	return c.Send(probability(message))
}

// handleWho responds with the message appended to the random chat member
func (a *app) handleWho(c tele.Context, message string) error {
	userID, err := a.getRandomGroupMember(c.Chat().ID)
	if err != nil {
		return err
	}
	chat, err := c.Bot().ChatByID(userID)
	if err != nil {
		return err
	}
	name := getUserName(chat)
	message = escapeMarkdown(message)
	return c.Send(who(userID, name, message), tele.ModeMarkdownV2)
}

// handleCat sends a picture of a cat
func (a *app) handleCat(c tele.Context) error {
	pic, err := fetchPicture(catURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

// handleTitle sets the user's admin title
func (a *app) handleTitle(c tele.Context, title string) error {
	if len(title) > 16 {
		return c.Send("Ошибка: максимальная длина имени 16 символов")
	}
	if err := c.Bot().SetAdminTitle(c.Chat(), c.Sender(), title); err != nil {
		return c.Send("Ошибка")
	}
	return nil
}

// handleAnime sends an anime picture
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

// handleFurry sends a furry picture
func (a *app) handleFurry(c tele.Context) error {
	seed := getRandomNumbers(5)
	url := fmt.Sprintf(furFormat, seed)

	pic, err := fetchPicture(url)
	if err != nil {
		return err
	}

	return c.Send(pic)
}

// handleFlag sends a picture of a flag
func (a *app) handleFlag(c tele.Context) error {
	seed := rand.Intn(5000)
	url := fmt.Sprintf(flagFormat, seed)

	pic, err := fetchPicture(url)
	if err != nil {
		return err
	}

	return c.Send(pic)
}

// handlePerson sends a picture of a person
func (a *app) handlePerson(c tele.Context) error {
	pic, err := fetchPicture(personURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

// handleHorse sends a picture of a horse
func (a *app) handleHorse(c tele.Context) error {
	pic, err := fetchPicture(horseURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

// handleArt sends a picture of an art
func (a *app) handleArt(c tele.Context) error {
	pic, err := fetchPicture(artURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

// handleCar sends a picture of a car
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

// handlePair sends a pair of the day
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
		mention(p.x, getUserName(chatX)),
		mention(p.y, getUserName(chatY))),
		tele.ModeMarkdownV2)
}

// handleEblan sends a eblan of the day
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

	return c.Send(fmt.Sprintf("Еблан дня: %s 😸", mention(userID, getUserName(chat))), tele.ModeMarkdownV2)
}

// getRandomGroupMember returns the random group member's ID
func (a *app) getRandomGroupMember(groupID int64) (int64, error) {
	userIDs, err := a.store.getUserIDs(groupID)
	if err != nil {
		return 0, err
	}
	return userIDs[rand.Intn(len(userIDs))], nil
}

// getRandomNumbers returns a string of random numbers of length c
func getRandomNumbers(c int) string {
	nums := []string{}
	for i := 0; i < c; i++ {
		n := rand.Intn(10)
		nums = append(nums, fmt.Sprintf("%d", n))
	}
	return strings.Join(nums, "")
}

// getUserName returns the displayed user name
func getUserName(chat *tele.Chat) string {
	return strings.TrimSpace(strings.Join([]string{chat.FirstName, chat.LastName}, " "))
}

// probability returns the probability of the message
func probability(message string) string {
	p := rand.Intn(101)
	i := rand.Intn(len(probabilityTemplates))
	t := probabilityTemplates[i]
	return fmt.Sprintf(t, message, p)
}

// who returns the mention of the user prepended to the message
func who(userID int64, name, message string) string {
	return fmt.Sprintf("%s %s", mention(userID, name), message)
}

// mention returns the mention to the user under the name
func mention(userID int64, name string) string {
	return fmt.Sprintf("[%s](tg://user?id=%d)", name, userID)
}

// fetchPicture returns a picture located at url
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

// byteSliceToPhoto converts the byte slice of image data to Photo
func byteSliceToPhoto(data []byte) *tele.Photo {
	return &tele.Photo{File: tele.FromReader(bytes.NewReader(data))}
}

// escapeMarkdown escapes the message for its use in Markdown
func escapeMarkdown(message string) string {
	chars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">",
		"#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, c := range chars {
		message = strings.ReplaceAll(message, c, fmt.Sprintf("\\%s", c))
	}
	return message
}
