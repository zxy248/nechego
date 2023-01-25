package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"math/rand"
	"nechego/avatar"
	"nechego/danbooru"
	"nechego/format"
	"nechego/game"
	"nechego/teleutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Infa struct{}

var infaRe = re("^!инфа ?(.*)")

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
	return c.Send(fmt.Sprintf(templates[rand.Intn(len(templates))],
		teleutil.Args(c, infaRe)[1],
		rand.Intn(101)))
}

type Who struct {
	Universe *game.Universe
}

var whoRe = re("^!кто ?(.*)")

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

	user := w.RandomUser()
	return c.Send(teleutil.Mention(c, user.TUID)+" "+
		html.EscapeString(teleutil.Args(c, whoRe)[1]), tele.ModeHTML)
}

type List struct {
	Universe *game.Universe
}

var listRe = re("^!список ?(.*)")

func (h *List) Match(s string) bool {
	return listRe.MatchString(s)
}

func (h *List) Handle(c tele.Context) error {
	world, _ := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	users := world.RandomUsers(3 + rand.Intn(3))
	arg := teleutil.Args(c, listRe)[1]
	s := []string{fmt.Sprintf("<b>📝 Список %s</b>", arg)}
	for _, u := range users {
		mention := teleutil.Mention(c, teleutil.Member(c, tele.ChatID(u.TUID)))
		s = append(s, fmt.Sprintf("<b>•</b> %s", mention))
	}
	return c.Send(strings.Join(s, "\n"), tele.ModeHTML)
}

type Top struct {
	Universe *game.Universe
}

var topRe = re("^!топ ?(.*)")

func (h *Top) Match(s string) bool {
	return topRe.MatchString(s)
}

func (h *Top) Handle(c tele.Context) error {
	world, _ := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	users := world.RandomUsers(3 + rand.Intn(3))
	arg := teleutil.Args(c, topRe)[1]
	s := []string{fmt.Sprintf("<b>🏆 Топ %s</b>", arg)}
	for i, u := range users {
		mention := teleutil.Mention(c, teleutil.Member(c, tele.ChatID(u.TUID)))
		s = append(s, fmt.Sprintf("<i>%d.</i> %s", i+1, mention))
	}
	return c.Send(strings.Join(s, "\n"), tele.ModeHTML)
}

type Mouse struct {
	Path string // Path to the video file.
}

var mouseRe = re("^!мыш")

func (h *Mouse) Match(s string) bool {
	return mouseRe.MatchString(s)
}

func (h *Mouse) Handle(c tele.Context) error {
	return c.Send(&tele.Video{File: tele.FromDisk(h.Path)})
}

type Tiktok struct {
	Path string // Path to the directory with WebMs.
}

var tiktokRe = re("^!тикток")

func (h *Tiktok) Match(s string) bool {
	return tiktokRe.MatchString(s)
}

func (h *Tiktok) Handle(c tele.Context) error {
	files, err := os.ReadDir(h.Path)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("empty directory %s", h.Path)
	}
	f := files[rand.Intn(len(files))]
	return c.Send(&tele.Video{File: tele.FromDisk(filepath.Join(h.Path, f.Name()))})
}

type Game struct{}

var gameRe = re("^!игр")

func (h *Game) Match(s string) bool {
	return gameRe.MatchString(s)
}

func (h *Game) Handle(c tele.Context) error {
	games := [...]*tele.Dice{tele.Dart, tele.Ball, tele.Goal, tele.Slot, tele.Bowl}
	return c.Send(games[rand.Intn(len(games))])
}

type Weather struct{}

var weatherRe = re("^!погода (.*)")

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
	city := url.PathEscape(teleutil.Args(c, weatherRe)[1])

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

var catRe = re("^!ко[тш]")

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

var animeRe = re("^!(аним|мульт)")

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

var furryRe = re("^!фур")

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

var flagRe = re("^!флаг")

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

var personRe = re("^!чел")

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

var horseRe = re("^!(лошад|конь)")

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

var artRe = re("^!арт")

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
	carRe    = re("^!(авто|машин|тачка)")
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
	dec := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(img))
	return c.Send(&tele.Photo{File: tele.FromReader(dec)})
}

type Soy struct{}

var soyRe = re("^!сой")

func (h *Soy) Match(s string) bool {
	return soyRe.MatchString(s)
}

func (h *Soy) Handle(c tele.Context) error {
	r, err := http.Get("https://booru.soy/random_image/download")
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(r.Body)})
}

type Danbooru struct{}

var danbooruRe = re("^!данб.ру")

func (h *Danbooru) Match(s string) bool {
	return danbooruRe.MatchString(s)
}

func (h *Danbooru) Handle(c tele.Context) error {
	pic, err := danbooru.Get()
	if err != nil {
		return err
	}
	photo := &tele.Photo{File: tele.FromReader(bytes.NewReader(pic.Data))}
	if pic.Rating == danbooru.Explicit {
		caps := [...]string{
			"🔞 Осторожно! Только для взрослых.",
			"<i>Содержимое предназначено для просмотра лицами старше 18 лет.</i>",
			"<b>ВНИМАНИЕ!</b> Вы увидите фотографии взрослых голых женщин. Будьте сдержанны.",
		}
		photo.Caption = caps[rand.Intn(len(caps))]
		photo.HasSpoiler = true
	}
	return c.Send(photo, tele.ModeHTML)
}

type Fap struct{}

var fapRe = re("^!(др[ао]ч|фап)")

func (h *Fap) Match(s string) bool {
	return fapRe.MatchString(s)
}

func (h *Fap) Handle(c tele.Context) error {
	pic, err := danbooru.GetNSFW()
	if err != nil {
		return err
	}
	photo := &tele.Photo{File: tele.FromReader(bytes.NewReader(pic.Data))}
	switch pic.Rating {
	case danbooru.Explicit:
		photo.Caption = "🔞"
	case danbooru.Questionable:
		photo.Caption = "❓"
	}
	photo.HasSpoiler = true
	return c.Send(photo, tele.ModeHTML)
}

type Masyunya struct{}

var masyunyaRe = re("^!ма[нс]ю[нс][а-я]*[пая]")

func (h *Masyunya) Match(s string) bool {
	return masyunyaRe.MatchString(s)
}

func (h *Masyunya) Handle(c tele.Context) error {
	set, err := c.Bot().StickerSet("masyunya_vk")
	if err != nil {
		return err
	}
	return c.Send(&set.Stickers[rand.Intn(len(set.Stickers))])
}

type Poppy struct{}

var poppyRe = re("^!паппи")

func (h *Poppy) Match(s string) bool {
	return poppyRe.MatchString(s)
}

func (h *Poppy) Handle(c tele.Context) error {
	names := []string{"pappy2_vk", "poppy_vk"}
	set, err := c.Bot().StickerSet(names[rand.Intn(len(names))])
	if err != nil {
		return err
	}
	return c.Send(&set.Stickers[rand.Intn(len(set.Stickers))])
}

type Sima struct{}

var simaRe = re("^!сима")

func (h *Sima) Match(s string) bool {
	return simaRe.MatchString(s)
}

func (h *Sima) Handle(c tele.Context) error {
	set, err := c.Bot().StickerSet("catsima_vk")
	if err != nil {
		return err
	}
	return c.Send(&set.Stickers[rand.Intn(len(set.Stickers))])
}

type Hello struct {
	Path  string
	cache []tele.Sticker
}

var helloRe = re("^!(п[рл]ив[а-я]*|хай|зд[ао]ров[а-я]*|ку|здрав[а-я]*)")

func (h *Hello) Match(s string) bool {
	return helloRe.MatchString(s)
}

func (h *Hello) Handle(c tele.Context) error {
	if h.cache == nil {
		f, err := os.Open(h.Path)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := json.NewDecoder(f).Decode(&h.cache); err != nil {
			return err
		}
	}
	return c.Send(&h.cache[rand.Intn(len(h.cache))])
}

type Basili struct {
	Path string
}

var basiliRe = re("^!(муся|марс|кот василия|кошка василия)")

func (h *Basili) Match(s string) bool {
	return basiliRe.MatchString(s)
}

func (h *Basili) Handle(c tele.Context) error {
	files, err := os.ReadDir(h.Path)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("empty directory %s", h.Path)
	}
	f := files[rand.Intn(len(files))]
	return c.Send(&tele.Photo{File: tele.FromDisk(filepath.Join(h.Path, f.Name()))})
}

type Casper struct {
	Path string
}

var casperRe = re("^!касп[ие]р")

func (h *Casper) Match(s string) bool {
	return casperRe.MatchString(s)
}

func (h *Casper) Handle(c tele.Context) error {
	files, err := os.ReadDir(h.Path)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("empty directory %s", h.Path)
	}
	f := files[rand.Intn(len(files))]
	return c.Send(&tele.Photo{File: tele.FromDisk(filepath.Join(h.Path, f.Name()))})
}

type Zeus struct {
	Path string
}

var zeusRe = re("^!зевс")

func (h *Zeus) Match(s string) bool {
	return zeusRe.MatchString(s)
}

func (h *Zeus) Handle(c tele.Context) error {
	files, err := os.ReadDir(h.Path)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("empty directory %s", h.Path)
	}
	f := files[rand.Intn(len(files))]
	return c.Send(&tele.Photo{File: tele.FromDisk(filepath.Join(h.Path, f.Name()))})
}

type Pic struct {
	Path string
}

var picRe = re("^!пик")

func (h *Pic) Match(s string) bool {
	return picRe.MatchString(s)
}

func (h *Pic) Handle(c tele.Context) error {
	dirs, err := os.ReadDir(h.Path)
	if err != nil {
		return err
	}
	if len(dirs) == 0 {
		return fmt.Errorf("empty directory %s", h.Path)
	}
	d := dirs[rand.Intn(len(dirs))]
	files, err := os.ReadDir(filepath.Join(h.Path, d.Name()))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("empty directory %s", h.Path)
	}
	f := files[rand.Intn(len(files))]
	return c.Send(&tele.Photo{File: tele.FromDisk(filepath.Join(h.Path, d.Name(), f.Name()))})
}

type Avatar struct {
	Avatars *avatar.Storage
}

var avatarRe = re("^!ава")

func (h *Avatar) Match(s string) bool {
	return avatarRe.MatchString(s)
}

func (h *Avatar) Handle(c tele.Context) error {
	if c.Message().Photo == nil {
		if avatar, ok := h.Avatars.Get(c.Sender().ID); ok {
			return c.Send(avatar)
		}
		return c.Send("📷 Прикрепите изображение.")
	}
	if err := h.Avatars.Set(c.Sender().ID, c.Message().Photo); errors.Is(err, avatar.ErrSize) {
		return c.Send("📷 Максимальный размер аватара %dx%d пикселей.",
			h.Avatars.MaxWidth, h.Avatars.MaxHeight)
	} else if err != nil {
		return err
	}
	return c.Send("📸 Аватар установлен.")
}

type TurnOn struct {
	Universe *game.Universe
}

var turnOnRe = re("^!(вкл|подкл|подруб)")

func (h *TurnOn) Match(s string) bool {
	return turnOnRe.MatchString(s)
}

func (h *TurnOn) Handle(c tele.Context) error {
	emoji := [...]string{"🔈", "🔔", "✅", "🆗", "▶️"}
	return c.Send(emoji[rand.Intn(len(emoji))])
}

type TurnOff struct {
	Universe *game.Universe
}

var turnOffRe = re("^!(выкл|откл)")

func (h *TurnOff) Match(s string) bool {
	return turnOffRe.MatchString(s)
}

func (h *TurnOff) Handle(c tele.Context) error {
	emoji := [...]string{"🔇", "🔕", "💤", "❌", "⛔️", "🚫", "⏹"}
	return c.Send(emoji[rand.Intn(len(emoji))])
}

type Ban struct {
	Universe *game.Universe
}

var banRe = re("^!бан")

func (h *Ban) Match(s string) bool {
	return banRe.MatchString(s)
}

func (h *Ban) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	if !user.Admin() {
		return c.Send(format.AdminsOnly)
	}
	reply, ok := teleutil.Reply(c)
	if !ok {
		return c.Send(format.RepostMessage)
	}
	world.UserByID(reply.ID).Banned = true
	return c.Send(format.UserBanned)
}

type Unban struct {
	Universe *game.Universe
}

var unbanRe = re("^!разбан")

func (h *Unban) Match(s string) bool {
	return unbanRe.MatchString(s)
}

func (h *Unban) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	if !user.Admin() {
		return c.Send(format.AdminsOnly)
	}
	reply, ok := teleutil.Reply(c)
	if !ok {
		return c.Send(format.RepostMessage)
	}
	world.UserByID(reply.ID).Banned = false
	return c.Send(format.UserUnbanned)
}
