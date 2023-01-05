package handlers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html"
	"io"
	"math/rand"
	"nechego/game"
	"nechego/teleutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Infa struct{}

var infaRe = regexp.MustCompile("!инфа (.*)")

func (h *Infa) Match(s string) bool {
	return infaRe.MatchString(s)
}

func (h *Infa) Handle(c tele.Context) error {
	templates := [...]string{
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
	tmpl := templates[rand.Intn(len(templates))]
	arg := infaRe.FindStringSubmatch(c.Message().Text)[1]
	return c.Send(fmt.Sprintf(tmpl, arg, rand.Intn(101)))
}

type Who struct {
	Universe *game.Universe
}

var whoRe = regexp.MustCompile("!кто (.*)")

func (h *Who) Match(s string) bool {
	return whoRe.MatchString(s)
}

func (h *Who) Handle(c tele.Context) error {
	w, err := h.Universe.World(c.Chat().ID)
	if err != nil {
		return err
	}
	w.Lock()
	defer w.Unlock()

	u := w.RandomUser()
	arg := whoRe.FindStringSubmatch(c.Message().Text)[1]
	msg := fmt.Sprintf("%s %s", teleutil.Mention(c, w.TGID, u.TUID), html.EscapeString(arg))
	return c.Send(msg, tele.ModeHTML)
}

type Mouse struct {
	Path string // path to video file
}

var mouseRe = regexp.MustCompile("^!мыш")

func (h *Mouse) Match(s string) bool {
	return mouseRe.MatchString(s)
}

func (h *Mouse) Handle(c tele.Context) error {
	return c.Send(&tele.Video{File: tele.FromDisk(h.Path)})
}

type Tiktok struct {
	Path string // path to directory with webms
}

var tiktokRe = regexp.MustCompile("^!тикток")

func (h *Tiktok) Match(s string) bool {
	return tiktokRe.MatchString(s)
}

func (h *Tiktok) Handle(c tele.Context) error {
	files, err := os.ReadDir(h.Path)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("empty directory")
	}
	f := files[rand.Intn(len(files))]
	return c.Send(&tele.Video{File: tele.FromDisk(filepath.Join(h.Path, f.Name()))})
}

type Game struct{}

var gameRe = regexp.MustCompile("^!игр")

func (h *Game) Match(s string) bool {
	return gameRe.MatchString(s)
}

func (h *Game) Handle(c tele.Context) error {
	games := [...]*tele.Dice{tele.Dart, tele.Ball, tele.Goal, tele.Slot, tele.Bowl}
	return c.Send(games[rand.Intn(len(games))])
}

type Weather struct{}

var weatherRe = regexp.MustCompile("^!погода (.*)")

func (h *Weather) Match(s string) bool {
	return weatherRe.MatchString(s)
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
	city := teleutil.Args(c, weatherRe)[1]

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

type Cat struct{}

var catRe = regexp.MustCompile("^!ко[тш]")

func (h *Cat) Match(s string) bool {
	return catRe.MatchString(s)
}

func (h *Cat) Handle(c tele.Context) error {
	addr := "https://thiscatdoesnotexist.com/"
	r, err := http.Get(addr)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(r.Body)})
}

type Anime struct{}

var animeRe = regexp.MustCompile("^!(аним|мульт)")

func (h *Anime) Match(s string) bool {
	return animeRe.MatchString(s)
}

func (h *Anime) Handle(c tele.Context) error {
	const format = "https://thisanimedoesnotexist.ai/results/psi-%s/seed%05d.png"
	psis := [...]string{"0.3", "0.4", "0.5", "0.6", "0.7", "0.8", "0.9", "1.0",
		"1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "2.0"}
	psi := psis[rand.Intn(len(psis))]
	addr := fmt.Sprintf(format, psi, rand.Intn(100_000))
	r, err := http.Get(addr)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(r.Body)})
}

type Furry struct{}

var furryRe = regexp.MustCompile("^!фур")

func (h *Furry) Match(s string) bool {
	return furryRe.MatchString(s)
}

func (h *Furry) Handle(c tele.Context) error {
	const format = "https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%05d.jpg"
	addr := fmt.Sprintf(format, rand.Intn(100_000))
	r, err := http.Get(addr)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(r.Body)})
}

type Flag struct{}

var flagRe = regexp.MustCompile("^!флаг")

func (h *Flag) Match(s string) bool {
	return flagRe.MatchString(s)
}

func (h *Flag) Handle(c tele.Context) error {
	const format = "https://thisflagdoesnotexist.com/images/%d.png"
	addr := fmt.Sprintf(format, rand.Intn(5000))
	r, err := http.Get(addr)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(r.Body)})
}

type Person struct{}

var personRe = regexp.MustCompile("^!чел")

func (h *Person) Match(s string) bool {
	return personRe.MatchString(s)
}

func (h *Person) Handle(c tele.Context) error {
	const addr = "https://thispersondoesnotexist.com/image"
	r, err := http.Get(addr)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(r.Body)})
}

type Horse struct{}

var horseRe = regexp.MustCompile("^!(лошад|конь)")

func (h *Horse) Match(s string) bool {
	return horseRe.MatchString(s)
}

func (h *Horse) Handle(c tele.Context) error {
	const addr = "https://thishorsedoesnotexist.com/"
	r, err := http.Get(addr)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(r.Body)})
}

type Art struct{}

var artRe = regexp.MustCompile("^!арт")

func (h *Art) Match(s string) bool {
	return artRe.MatchString(s)
}

func (h *Art) Handle(c tele.Context) error {
	const addr = "https://thisartworkdoesnotexist.com/"
	r, err := http.Get(addr)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(r.Body)})
}

type Car struct{}

var (
	carRe    = regexp.MustCompile("^!(авто|машин|тачка)")
	carImgRe = regexp.MustCompile(`<img id = "vehicle" src="data:image/png;base64,(.+)" class="center">`)
)

func (h *Car) Match(s string) bool {
	return carRe.MatchString(s)
}

func (h *Car) Handle(c tele.Context) error {
	const addr = "https://www.thisautomobiledoesnotexist.com/"
	r, err := http.Get(addr)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	img := carImgRe.FindSubmatch(data)[1]
	buf := bytes.NewBuffer(img)
	dec := base64.NewDecoder(base64.StdEncoding, buf)
	return c.Send(&tele.Photo{File: tele.FromReader(dec)})
}
