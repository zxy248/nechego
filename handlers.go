package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"

	tele "gopkg.in/telebot.v3"
)

// infaTemplates regexp: "^.*%s.*%d%%\""
var infaTemplates = []string{
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

const kotURL = "https://thiscatdoesnotexist.com/"

const animeFormat = "https://thisanimedoesnotexist.ai/results/psi-%s/seed%s.png"

var animePsis = []string{"0.3", "0.4", "0.5", "0.6", "0.7", "0.8",
	"0.9", "1.0", "1.1", "1.2", "1.3", "1.4", "1.5",
	"1.6", "1.7", "1.8", "2.0"}

const furFormat = "https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%s.jpg"

const flagFormat = "https://thisflagdoesnotexist.com/images/%d.png"

const chelURL = "https://thispersondoesnotexist.com/image"

// handleInfa responds with the probability of message happening
func (a *app) handleInfa(c tele.Context, message string) error {
	return c.Send(infa(message))
}

// handleKto responds with message appended to the random chat member username
func (a *app) handleKto(c tele.Context, message string) error {
	userID, err := a.getRandomGroupMember(c.Chat().ID)
	if err != nil {
		return err
	}
	chat, err := c.Bot().ChatByID(userID)
	if err != nil {
		return err
	}
	name := getUserName(chat)
	return c.Send(kto(userID, name, message), tele.ModeMarkdownV2)
}

// handleKot responds with a cat picture
func (a *app) handleKot(c tele.Context) error {
	pic, err := fetchPicture(kotURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

// handleImya sets the user's admin title
func (a *app) handleImya(c tele.Context, title string) error {
	if len(title) > 16 {
		return c.Send("Ошибка: максимальная длина имени 16 символов")
	}
	if err := c.Bot().SetAdminTitle(c.Chat(), c.Sender(), title); err != nil {
		return c.Send("Ошибка")
	}
	return nil
}

// handleAnime responds with an anime picture
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

// handleFur responds with a furry picture
func (a *app) handleFur(c tele.Context) error {
	seed := getRandomNumbers(5)
	url := fmt.Sprintf(furFormat, seed)

	pic, err := fetchPicture(url)
	if err != nil {
		return err
	}

	return c.Send(pic)
}

// handleFlag respons with a flag picture
func (a *app) handleFlag(c tele.Context) error {
	seed := rand.Intn(5000)
	url := fmt.Sprintf(flagFormat, seed)

	pic, err := fetchPicture(url)
	if err != nil {
		return err
	}

	return c.Send(pic)
}

// handleChel respons with a human picture
func (a *app) handleChel(c tele.Context) error {
	pic, err := fetchPicture(chelURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

// getRandomGroupMember returns the random member's ID from the group
func (a *app) getRandomGroupMember(groupID int64) (int64, error) {
	userIDs, err := a.store.getUserIDs(groupID)
	if err != nil {
		return 0, err
	}
	return userIDs[rand.Intn(len(userIDs))], nil
}

// getRandomNumbers generates a string of random numbers of length c
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

// infa returns the probability of message happening
func infa(message string) string {
	p := rand.Intn(101)
	i := rand.Intn(len(infaTemplates))
	t := infaTemplates[i]
	return fmt.Sprintf(t, message, p)
}

// kto returns the mention for user together with message
func kto(userID int64, name, message string) string {
	return fmt.Sprintf("[%s](tg://user?id=%d) %s", name, userID, message)
}

// fetchPicture returns the picture located at url
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
	return &tele.Photo{File: tele.FromReader(bytes.NewReader(body))}, nil
}
