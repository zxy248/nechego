package main

import (
	"fmt"
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

// handleInfa responds with the probability of message happening
func (a *app) handleInfa(c tele.Context, message string) error {
	p := rand.Intn(101)
	i := rand.Intn(len(infaTemplates))
	t := infaTemplates[i]
	return c.Send(fmt.Sprintf(t, message, p))
}

// handleKto responds with message appended to the random chat member username
func (a *app) handleKto(c tele.Context, message string) error {
	userID, err := a.getRandomGroupMember(c.Chat().ID)
	if err != nil {
		return err
	}
	user, err := c.Bot().ChatByID(userID)
	if err != nil {
		return err
	}
	name := getUserName(user)
	return c.Send(fmt.Sprintf("[%s](tg://user?id=%d) %s", name, userID, message), tele.ModeMarkdownV2)
}

// handleKot responds with a cat picture
func (a *app) handleKot(c tele.Context) error {
	l := "https://thiscatdoesnotexist.com/"
	r, err := http.Get(l)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	photo := &tele.Photo{File: tele.FromReader(r.Body)}
	return c.Send(photo)
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
	l := "https://thisanimedoesnotexist.ai/results/psi-1.0/seed%s.png"
	r, err := http.Get(fmt.Sprintf(l, getRandomNumbers(5)))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	photo := &tele.Photo{File: tele.FromReader(r.Body)}
	return c.Send(photo)
}

// handleFur responds with a furry picture
func (a *app) handleFur(c tele.Context) error {
	l := "https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%s.jpg"
	r, err := http.Get(fmt.Sprintf(l, getRandomNumbers(5)))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	photo := &tele.Photo{File: tele.FromReader(r.Body)}
	return c.Send(photo)
}

// getRandomGroupMember returns the random member's ID from the group
func (a *app) getRandomGroupMember(groupID int64) (int64, error) {
	userIDs, err := a.store.getUserIDs(groupID)
	if err != nil {
		return 0, err
	}
	i := rand.Intn(len(userIDs))
	return userIDs[i], nil
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
