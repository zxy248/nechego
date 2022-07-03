package app

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"nechego/input"
	"nechego/model"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	tele "gopkg.in/telebot.v3"
)

// handleProbability responds with a probability of the message.
func (a *App) handleProbability(c tele.Context) error {
	m := getMessage(c).Argument()
	return c.Send(probability(m))
}

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

// probability returns a probability of the message.
func probability(message string) string {
	t := probabilityTemplates[rand.Intn(len(probabilityTemplates))]
	p := rand.Intn(101)
	return fmt.Sprintf(t, message, p)
}

// handleWho responds with the message appended to the random chat member.
func (a *App) handleWho(c tele.Context) error {
	gid := c.Chat().ID
	uid, err := a.model.Users.Random(gid)
	if err != nil {
		return err
	}
	memb, err := a.chatMember(gid, uid)
	if err != nil {
		return err
	}
	name := markdownEscaper.Replace(chatMemberName(memb))
	message := markdownEscaper.Replace(getMessage(c).Argument())
	return c.Send(who(uid, name, message), tele.ModeMarkdownV2)
}

const catURL = "https://thiscatdoesnotexist.com/"

// handleCat sends a picture of a cat.
func (a *App) handleCat(c tele.Context) error {
	pic, err := fetchPicture(catURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

const (
	maxNameLength = 16
	nameTooLong   = "Максимальная длина имени 16 символов"
	yourName      = "Ваше имя: *%s* 🔖"
	pleaseReEnter = "Для использования этой функции Вам необходимо перезайти в беседу"
	nameSet       = "Имя *%v* установлено ✅"
)

// handleTitle sets the admin title of the sender.
func (a *App) handleTitle(c tele.Context) error {
	group := c.Chat()
	sender := c.Sender()
	gid := group.ID
	uid := sender.ID
	title := getMessage(c).Argument()

	if title == "" {
		m, err := a.chatMember(gid, uid)
		if err != nil {
			return err
		}
		name := markdownEscaper.Replace(chatMemberName(m))
		return c.Send(fmt.Sprintf(yourName, name), tele.ModeMarkdownV2)
	}
	if utf8.RuneCountInString(title) > maxNameLength {
		return c.Send(makeError(nameTooLong))
	}
	if err := c.Bot().SetAdminTitle(group, sender, title); err != nil {
		return c.Send(makeError(pleaseReEnter))
	}
	return c.Send(fmt.Sprintf(nameSet, markdownEscaper.Replace(title)), tele.ModeMarkdownV2)
}

const animeFormat = "https://thisanimedoesnotexist.ai/results/psi-%s/seed%s.png"

var animePsis = []string{"0.3", "0.4", "0.5", "0.6", "0.7", "0.8", "0.9", "1.0",
	"1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "2.0"}

// handleAnime sends an anime picture.
func (a *App) handleAnime(c tele.Context) error {
	psi := animePsis[rand.Intn(len(animePsis))]
	seed := randomNumbers(5)
	url := fmt.Sprintf(animeFormat, psi, seed)
	return a.fetchAndSend(c, url)
}

const furFormat = "https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%s.jpg"

// handleFurry sends a furry picture.
func (a *App) handleFurry(c tele.Context) error {
	seed := randomNumbers(5)
	url := fmt.Sprintf(furFormat, seed)
	return a.fetchAndSend(c, url)
}

const flagFormat = "https://thisflagdoesnotexist.com/images/%d.png"

// handleFlag sends a picture of a flag.
func (a *App) handleFlag(c tele.Context) error {
	seed := rand.Intn(5000)
	url := fmt.Sprintf(flagFormat, seed)
	return a.fetchAndSend(c, url)
}

const personURL = "https://thispersondoesnotexist.com/image"

// handlePerson sends a picture of a person.
func (a *App) handlePerson(c tele.Context) error {
	return a.fetchAndSend(c, personURL)
}

const horseURL = "https://thishorsedoesnotexist.com/"

// handleHorse sends a picture of a horse.
func (a *App) handleHorse(c tele.Context) error {
	return a.fetchAndSend(c, horseURL)
}

const artURL = "https://thisartworkdoesnotexist.com/"

// handleArt sends a picture of an art.
func (a *App) handleArt(c tele.Context) error {
	return a.fetchAndSend(c, artURL)
}

const carURL = "https://www.thisautomobiledoesnotexist.com/"

var carImageRe = regexp.MustCompile(
	"<img id = \"vehicle\" src=\"data:image/png;base64,(.+)\" class=\"center\">")

// handleCar sends a picture of a car.
func (a *App) handleCar(c tele.Context) error {
	r, err := http.Get(carURL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	b64 := carImageRe.FindStringSubmatch(string(data))[1]
	img, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	return c.Send(photoFromBytes(img))
}

const pairOfTheDayFormat = "Пара дня ✨\n%s 💘 %s"

// handlePair sends the current pair of the day, randomly choosing a new pair if needed.
func (a *App) handlePair(c tele.Context) error {
	gid := c.Chat().ID
	x, y, err := a.getDailyPair(gid)
	if err != nil {
		if errors.Is(err, model.ErrNoPair) {
			return c.Send("💔")
		}
		return err
	}

	mx, err := a.chatMember(gid, x)
	if err != nil {
		return err
	}
	my, err := a.chatMember(gid, y)
	if err != nil {
		return err
	}
	namex := markdownEscaper.Replace(chatMemberName(mx))
	namey := markdownEscaper.Replace(chatMemberName(my))
	return c.Send(fmt.Sprintf(pairOfTheDayFormat,
		mention(x, namex), mention(y, namey)), tele.ModeMarkdownV2)
}

const eblanOfTheDayFormat = "Еблан дня: %s 😸"

// handleEblan sends the current eblan of the day, randomly choosing a new one if needed.
func (a *App) handleEblan(c tele.Context) error {
	gid := c.Chat().ID
	uid, err := a.getDaily(gid, a.model.Eblans.Get, a.model.Eblans.Insert, model.ErrNoEblan)
	if err != nil {
		return err
	}
	m, err := a.chatMember(gid, uid)
	if err != nil {
		return err
	}
	name := markdownEscaper.Replace(chatMemberName(m))
	return c.Send(fmt.Sprintf(eblanOfTheDayFormat, mention(uid, name)), tele.ModeMarkdownV2)
}

const adminOfTheDayFormat = "Админ дня: %s 👑"

func (a *App) handleAdmin(c tele.Context) error {
	gid := c.Chat().ID
	uid, err := a.getDaily(gid, a.model.Admins.GetDaily, a.model.Admins.InsertDaily, model.ErrNoAdmin)
	if err != nil {
		return err
	}
	m, err := a.chatMember(gid, uid)
	if err != nil {
		return err
	}
	name := markdownEscaper.Replace(chatMemberName(m))
	return c.Send(fmt.Sprintf(adminOfTheDayFormat, mention(uid, name)), tele.ModeMarkdownV2)
}

const masyunyaStickersName = "masyunya_vk"

func (a *App) masyunyaHandler() tele.HandlerFunc {
	set, err := a.bot.StickerSet(masyunyaStickersName)
	if err != nil {
		log.Println("masyunyaHandler unavailable: ", err)
		return func(c tele.Context) error {
			return nil
		}
	}
	return func(c tele.Context) error {
		return c.Send(&set.Stickers[rand.Intn(len(set.Stickers))])
	}
}

var poppyStickersNames = []string{"pappy2_vk", "poppy_vk"}

func (a *App) poppyHandler() tele.HandlerFunc {
	var stickers []tele.Sticker
	for _, sn := range poppyStickersNames {
		set, err := a.bot.StickerSet(sn)
		if err != nil {
			log.Println("poppyHandler unavailable: ", err)
			return func(c tele.Context) error {
				return nil
			}
		}
		stickers = append(stickers, set.Stickers...)
	}
	return func(c tele.Context) error {
		return c.Send(&stickers[rand.Intn(len(stickers))])
	}
}

const helloChance = 0.2

// handleHello sends a hello sticker
func (a *App) handleHello(c tele.Context) error {
	if strings.HasPrefix(getMessage(c).Raw, "!") || rand.Float64() <= helloChance {
		return c.Send(helloSticker())
	}
	return nil
}

var (
	mouseVideoPath = filepath.Join(dataPath, "mouse.mp4")
	mouseVideo     = &tele.Video{File: tele.FromDisk(mouseVideoPath)}
)

// handleMouse sends the mouse video
func (a *App) handleMouse(c tele.Context) error {
	return c.Send(mouseVideo)
}

const weatherTimeout = 10 * time.Second
const weatherURL = "https://wttr.in/"
const weatherFormat = `?format=%l:+%c+%t+\nОщущается+как+%f\n\nВетер+—+%w\nВлажность+—+%h\nДавление+—+%P\nФаза+луны+—+%m\nУФ-индекс+—+%u\n`

// handleWeather sends the current weather for a given city
func (a *App) handleWeather(c tele.Context) error {
	place := getMessage(c).Argument()

	ctx, cancel := context.WithTimeout(context.Background(), weatherTimeout)
	defer cancel()

	l := weatherURL + place + weatherFormat
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, l, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept-Language", "ru")

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		if err.(*url.Error).Timeout() {
			return c.Send(makeError("Время запроса вышло ☔️"))
		}
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		if r.StatusCode == http.StatusNotFound {
			return c.Send(makeError("Такого места не существует ☔️"))
		}
		return c.Send(makeError("Неудачный запрос ☔️"))
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return c.Send(string(data))
}

var tikTokVideo = &tele.Video{File: tele.FromDisk("data/tiktok.mp4")}

func (a *App) handleTikTok(c tele.Context) error {
	return c.Send(tikTokVideo)
}

const (
	listTemplate  = "Список %s 📝\n%s"
	minListLength = 3
	maxListLength = 5
)

func (a *App) handleList(c tele.Context) error {
	gid := c.Chat().ID
	uids, err := a.model.Users.NRandom(gid, minListLength+rand.Intn(maxListLength-minListLength))
	if err != nil {
		return err
	}
	var list string
	for _, uid := range uids {
		m, err := a.chatMember(gid, uid)
		if err != nil {
			return err
		}
		name := markdownEscaper.Replace(chatMemberName(m))
		list = list + "— " + mention(uid, name) + "\n"
	}
	msg := markdownEscaper.Replace(getMessage(c).Argument())
	return c.Send(fmt.Sprintf(listTemplate, msg, list), tele.ModeMarkdownV2)
}

const (
	numberedTopTemplate   = "Топ %d %s 🏆\n%s"
	unnumberedTopTemplate = "Топ %s 🏆\n%s"
	maxTopNumber          = 5
)

func (a *App) handleTop(c tele.Context) error {
	gid := c.Chat().ID
	arg, err := getMessage(c).Dynamic()
	if err != nil {
		return err
	}
	argument := arg.(input.TopArgument)

	var number int
	if argument.Number != nil {
		number = *argument.Number
	} else {
		number = maxTopNumber
	}
	if number < 1 || number > maxTopNumber {
		return c.Send(errorSign())
	}
	uids, err := a.model.Users.NRandom(gid, number)
	if err != nil {
		return err
	}

	var top string
	for i, uid := range uids {
		m, err := a.chatMember(gid, uid)
		if err != nil {
			return err
		}
		name := markdownEscaper.Replace(chatMemberName(m))
		top = top + fmt.Sprintf("_%d\\._ %s\n", i+1, mention(uid, name))
	}

	s := markdownEscaper.Replace(argument.String)
	var result string
	if argument.Number != nil {
		result = fmt.Sprintf(numberedTopTemplate, number, s, top)
	} else {
		result = fmt.Sprintf(unnumberedTopTemplate, s, top)
	}
	return c.Send(result, tele.ModeMarkdownV2)
}

var (
	albumsPath     = filepath.Join(dataPath, "vk.com-albums")
	basiliCatsPath = filepath.Join(albumsPath, "basili")
	casperPath     = filepath.Join(albumsPath, "casper")
	zeusPath       = filepath.Join(albumsPath, "zeus")
	picPath        = filepath.Join(albumsPath, "pic")
)

// handleBasili sends a photo of the Basili's cat.
func (a *App) handleBasili(c tele.Context) error {
	return sendRandomFile(c, basiliCatsPath)
}

// handleBasili sends a photo of the Leonid's cat.
func (a *App) handleCasper(c tele.Context) error {
	return sendRandomFile(c, casperPath)
}

// handleZeus sends a photo of the Solar's cat.
func (a *App) handleZeus(c tele.Context) error {
	return sendRandomFile(c, zeusPath)
}

// handlePic sends a photo from a hierarchy of directories located at picPath.
func (a *App) handlePic(c tele.Context) error {
	return sendRandomFileWith(c, picPath, randomFileFromHierarchy)
}

// handleDice rolls a dice.
func (a *App) handleDice(c tele.Context) error {
	return c.Send(tele.Cube)
}

var games = []*tele.Dice{tele.Dart, tele.Ball, tele.Goal, tele.Slot, tele.Bowl}

func (a *App) handleGame(c tele.Context) error {
	game := games[rand.Intn(len(games))]
	return c.Send(game)
}

const handleFightTemplate = `
⚔️ Нападает %s, сила в бою ` + "`%.2f`" + `
🛡 Защищается %s, сила в бою ` + "`%.2f`" + `

🏆 %s выходит победителем и забирает ` + "`%v монет`" + `

⚡️ Энергии осталось: ` + "`%v`" + `
`

const fightEnergyUpdate = -1

// handleFight conducts a fight between two users.
func (a *App) handleFight(c tele.Context) error {
	gid := c.Chat().ID
	att := c.Sender().ID
	def := c.Message().ReplyTo.Sender.ID

	if att == def {
		return c.Send(makeError("Вы не можете напасть на самого себя"))
	}
	if c.Message().ReplyTo.Sender.IsBot {
		return c.Send(makeError("Можно напасть только на пользователя"))
	}
	exists, err := a.model.Users.Exists(gid, def)
	if err != nil {
		return err
	}
	if !exists {
		return c.Send(makeError("Неизвестный пользователь"))
	}
	energy0, err := a.model.Energy.Energy(gid, att)
	if err != nil {
		return err
	}
	if energy0 < 1 {
		return c.Send(makeError("Недостаточно энергии"))
	}

	attstr, err := a.userStrength(gid, att)
	if err != nil {
		return err
	}
	defstr, err := a.userStrength(gid, def)
	if err != nil {
		return err
	}

	attmem, err := a.chatMember(gid, att)
	if err != nil {
		return err
	}
	defmem, err := a.chatMember(gid, def)
	if err != nil {
		return err
	}
	attment := mention(att, markdownEscaper.Replace(chatMemberName(attmem)))
	defment := mention(def, markdownEscaper.Replace(chatMemberName(defmem)))

	var win int64
	var winment string
	if attstr > defstr {
		win = att
		winment = attment
	} else {
		win = def
		winment = defment
	}
	fmt.Println(win)

	// TODO: make use for money
	money := 0
	if err := a.model.Energy.Update(gid, att, fightEnergyUpdate); err != nil {
		return err
	}
	energy, err := a.model.Energy.Energy(gid, att)
	if err != nil {
		return err
	}
	s := fmt.Sprintf(handleFightTemplate,
		attment, attstr,
		defment, defstr,
		winment, money, energy)
	return c.Send(s, tele.ModeMarkdownV2)
}

// userStrength determines the strength of a user.
func (a *App) userStrength(gid, uid int64) (float64, error) {
	chance := rand.Float64()
	week := time.Hour * 24 * 7
	user, err := a.userMessageCount(gid, uid, week)
	if err != nil {
		return 0, err
	}
	total, err := a.totalMessageCount(gid, week)
	if err != nil {
		return 0, err
	}
	score := 1.0 / (float64(1+total) / float64(1+user))
	return chance * score, nil
}

// userMessageCount returns the number of messages sent by the user in the specified interval.
func (a *App) userMessageCount(gid, uid int64, interval time.Duration) (int, error) {
	c, err := a.model.Messages.UserCount(gid, uid, time.Now().Add(-interval))
	if err != nil {
		return 0, err
	}
	return c, nil
}

// totalMessageCount returns the number of messages sent in the specified interval.
func (a *App) totalMessageCount(gid int64, interval time.Duration) (int, error) {
	c, err := a.model.Messages.TotalCount(gid, time.Now().Add(-interval))
	if err != nil {
		return 0, err
	}
	return c, nil
}

const randomPhotoChance = 0.02

func (a *App) handleRandomPhoto(c tele.Context) error {
	if rand.Float64() <= randomPhotoChance {
		return sendSmallProfilePhoto(c)
	}
	return nil
}

// handleKeyboardOpen opens the keyboard.
func (a *App) handleKeyboardOpen(c tele.Context) error {
	return c.Send("Клавиатура ⌨️", keyboard)
}

// handleKeyboardClose closes the keyboard.
func (a *App) handleKeyboardClose(c tele.Context) error {
	return c.Send("Клавиатура закрыта 😣", tele.RemoveKeyboard)
}

var (
	emojisActive   = []string{"🔈", "🔔", "✅", "🆗", "▶️"}
	emojisInactive = []string{"🔇", "🔕", "💤", "❌", "⛔️", "🚫", "⏹"}
)

const (
	botTurnedOn         = "Бот включен %s"
	botAlreadyTurnedOn  = "Бот уже включен %s"
	botTurnedOff        = "Бот выключен %s"
	botAlreadyTurnedOff = "Бот уже выключен %s"
)

// handleTurnOn turns the bot on.
func (a *App) handleTurnOn(c tele.Context) error {
	emoji := emojisActive[rand.Intn(len(emojisActive))]
	gid := c.Chat().ID
	a.model.Status.Enable(gid)
	return c.Send(fmt.Sprintf(botTurnedOn, emoji))
}

// handleTurnOff turns the bot off.
func (a *App) handleTurnOff(c tele.Context) error {
	emoji := emojisInactive[rand.Intn(len(emojisInactive))]
	gid := c.Chat().ID
	a.model.Status.Disable(gid)
	return c.Send(fmt.Sprintf(botTurnedOff, emoji), tele.RemoveKeyboard)
}

const (
	userBlocked          = "Пользователь заблокирован 🚫"
	userAlreadyBlocked   = "Пользователь уже заблокирован 🛑"
	userUnblocked        = "Пользователь разблокирован ✅"
	userAlreadyUnblocked = "Пользователь не заблокирован ❎"
)

// handleBan adds the user ID of the reply message's sender to the ban list.
func (a *App) handleBan(c tele.Context) error {
	uid := c.Message().ReplyTo.Sender.ID
	banned, err := a.model.Bans.Banned(uid)
	if err != nil {
		return err
	}
	if banned {
		return c.Send(userAlreadyBlocked)
	}

	if err := a.model.Bans.Ban(uid); err != nil {
		return err
	}
	return c.Send(userBlocked)
}

// handleUnban removes the user ID of the reply message's sender from the ban list.
func (a *App) handleUnban(c tele.Context) error {
	uid := c.Message().ReplyTo.Sender.ID
	banned, err := a.model.Bans.Banned(uid)
	if err != nil {
		return err
	}
	if !banned {
		return c.Send(userAlreadyUnblocked)
	}

	if err := a.model.Bans.Unban(uid); err != nil {
		return err
	}
	return c.Send(userUnblocked)
}

const infoTemplate = "ℹ️ *Информация* 📌\n\n%s\n%s\n%s\n"

// handleInfo sends a few lists of useful information.
func (a *App) handleInfo(c tele.Context) error {
	gid := c.Chat().ID
	admins, err := a.adminList(gid)
	if err != nil {
		return err
	}
	bans, err := a.banList(gid)
	if err != nil {
		return err
	}
	commands, err := a.forbiddenCommandList(gid)
	if err != nil {
		return err
	}

	lists := fmt.Sprintf(infoTemplate, admins, bans, commands)
	return c.Send(lists, tele.ModeMarkdownV2)
}

const adminListTemplate = "👤 _Администрация_\n%s"

func (a *App) adminList(gid int64) (string, error) {
	l, err := a.model.Admins.List(gid)
	if err != nil {
		return "", err
	}
	var admins string
	for _, uid := range l {
		m, err := a.chatMember(gid, uid)
		if err != nil {
			return "", err
		}
		if !chatMemberPresent(m) {
			continue
		}
		name := markdownEscaper.Replace(chatMemberName(m))
		admins += "— " + mention(uid, name) + "\n"
	}
	if admins == "" {
		admins = "…\n"
	}
	return fmt.Sprintf(adminListTemplate, admins), nil
}

const banListTemplate = "🛑 _Черный список_\n%s"

func (a *App) banList(gid int64) (string, error) {
	l, err := a.model.Bans.List()
	if err != nil {
		return "", err
	}
	var banned string
	for _, uid := range l {
		m, err := a.chatMember(gid, uid)
		if err != nil {
			return "", err
		}
		if !chatMemberPresent(m) {
			continue
		}
		name := markdownEscaper.Replace(chatMemberName(m))
		banned += "— " + mention(uid, name) + "\n"
	}
	if banned == "" {
		banned = "…\n"
	}
	return fmt.Sprintf(banListTemplate, banned), nil
}

const forbiddenCommandListTemplate = "🔒 _Запрещенные команды_\n%s"

func (a *App) forbiddenCommandList(gid int64) (string, error) {
	l, err := a.model.Forbid.List(gid)
	if err != nil {
		return "", err
	}
	var commands string
	for _, c := range l {
		t := markdownEscaper.Replace(input.CommandText(c))
		commands += "— " + t + "\n"
	}
	if commands == "" {
		commands = "…\n"
	}
	return fmt.Sprintf(forbiddenCommandListTemplate, commands), nil
}

const help = `📖 *Команды* 📌

📄 _Базовые_
` +
	"— `!инфа\n`" +
	"— `!кто`\n" +
	"— `!список`\n" +
	"— `!топ`\n" +
	"— `!погода`\n" +
	"— `!пара дня`\n" +
	"— `!еблан дня`\n" +
	"— `!мыш`\n" +
	"— `!тикток`\n" +
	"— `!масюня` ||💖||\n" +
	"— `!паппи`\n" +
	"— `!игра`\n" +
	"— `!кости`\n" +
	`
🔮 _Нейросети_
` +
	"— `!кот`\n" +
	"— `!аниме`\n" +
	"— `!фурри`\n" +
	"— `!флаг`\n" +
	"— `!чел`\n" +
	"— `!лошадь`\n" +
	"— `!арт`\n" +
	"— `!авто`\n" +
	`
🐈 _Кошки_
` +
	"— `!марсик`\n" +
	"— `!муся`\n" +
	"— `!каспер`\n" +
	"— `!зевс`\n" +
	`
🔧 _Управление_
` +
	"— `!открыть`\n" +
	"— `!закрыть`\n" +
	"— `!включить`\n" +
	"— `!выключить`\n" +
	"— `!запретить`\n" +
	"— `!разрешить`\n" +
	"— `!бан`\n" +
	"— `!разбан`\n" +
	"— `!имя`\n" +
	"— `!информация`\n" +
	"— `!команды`\n"

func (a *App) handleHelp(c tele.Context) error {
	return c.Send(help, tele.ModeMarkdownV2)
}

func (a *App) handleJoin(c tele.Context) error {
	group := c.Chat()
	gid := group.ID
	uid := c.Message().UserJoined.ID

	m, err := a.chatMember(gid, uid)
	if err != nil {
		return err
	}

	if m.Role != tele.Administrator {
		m.Rights.CanBeEdited = true
		m.Rights.CanManageChat = true
		if err := c.Bot().Promote(group, m); err != nil {
			return err
		}
	}
	return c.Send(helloSticker())
}

const (
	commandForbidden        = "Команда запрещена 🚫"
	commandPermitted        = "Команда разрешена ✅"
	commandAlreadyForbidden = "Команда уже запрещена 🛑"
	commandAlreadyPermitted = "Команда уже разрешена ❎"
)

// handleForbid forbids a command.
func (a *App) handleForbid(c tele.Context) error {
	return a.handleCommandAction(c, func(command input.Command) error {
		gid := c.Chat().ID
		f, err := a.model.Forbid.Forbidden(gid, command)
		if err != nil {
			return err
		}
		if f {
			return c.Send(commandAlreadyForbidden)
		}
		if err := a.model.Forbid.Forbid(gid, command); err != nil {
			return err
		}
		return c.Send(commandForbidden)
	})
}

// handlePermit permits a command.
func (a *App) handlePermit(c tele.Context) error {
	return a.handleCommandAction(c, func(command input.Command) error {
		gid := c.Chat().ID
		f, err := a.model.Forbid.Forbidden(gid, command)
		if err != nil {
			return err
		}
		if !f {
			return c.Send(commandAlreadyPermitted)
		}
		if err := a.model.Forbid.Permit(gid, command); err != nil {
			return err
		}
		return c.Send(commandPermitted)
	})
}

// handleCommandAction performs an action on a command.
func (a *App) handleCommandAction(c tele.Context, action func(input.Command) error) error {
	arg, err := getMessage(c).Dynamic()
	if err != nil {
		if errors.Is(err, input.ErrNoCommand) {
			return c.Send(makeError("Укажите команду"))
		}
		if errors.Is(err, input.ErrUnknownCommand) {
			return c.Send(makeError("Неизвестная команда"))
		}
		return err
	}
	command := arg.(input.Command)
	return action(command)
}

// randomNumbers returns a string of random numbers of length c.
func randomNumbers(c int) string {
	var nums string
	for i := 0; i < c; i++ {
		n := rand.Intn(10)
		nums = nums + fmt.Sprint(n)
	}
	return nums
}

// displayedUsername returns the displayed user name.
func displayedUsername(chat *tele.Chat) string {
	return strings.TrimSpace(chat.FirstName + " " + chat.LastName)
}

func (a *App) isGroupMember(group tele.Recipient, user tele.Recipient) bool {
	member, err := a.bot.ChatMemberOf(group, user)
	if err != nil || member.Role == tele.Kicked || member.Role == tele.Left {
		return false
	}
	return true
}

// who returns the mention of the user prepended to the message.
func who(uid int64, name, message string) string {
	return mention(uid, name) + " " + message
}

// mention returns the mention of the user by the name.
func mention(uid int64, name string) string {
	return fmt.Sprintf("[%s](tg://user?id=%d)", name, uid)
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
	return photoFromBytes(body), nil
}

// photoFromBytes converts the image data to Photo.
func photoFromBytes(data []byte) *tele.Photo {
	return &tele.Photo{File: tele.FromReader(bytes.NewReader(data))}
}

// markdownEscaper escapes any character with the code between 1 and 126
// inclusively with a preceding backslash.
var markdownEscaper = func() *strings.Replacer {
	var table []string
	for i := 1; i <= 126; i++ {
		c := string(rune(i))
		table = append(table, c, "\\"+c)
	}
	return strings.NewReplacer(table...)
}()

func (a *App) chatMember(gid, uid int64) (*tele.ChatMember, error) {
	group, err := a.bot.ChatByID(gid)
	if err != nil {
		return nil, err
	}
	member, err := a.bot.ChatMemberOf(group, tele.ChatID(uid))
	if err != nil {
		return nil, err
	}
	if !chatMemberPresent(member) {
		a.model.Users.Delete(gid, uid)
	}
	return member, nil
}

func chatMemberPresent(m *tele.ChatMember) bool {
	if m.Role == tele.Kicked || m.Role == tele.Left {
		return false
	}
	return true
}

func chatMemberName(m *tele.ChatMember) string {
	name := m.Title
	if name == "" {
		name = m.User.FirstName + " " + m.User.LastName
	}
	return strings.TrimSpace(name)
}

func errorSign() string {
	errors := []string{"❌", "🚫", "⭕️", "🛑", "⛔️", "📛", "💢", "❗️", "‼️", "⚠️"}
	return errors[rand.Intn(len(errors))]
}

func makeError(s string) string {
	return errorSign() + " " + s
}

type dailyGet func(gid int64) (int64, error)
type dailyInsert func(gid, uid int64) error

func (a *App) getDaily(gid int64, get dailyGet, insert dailyInsert, e error) (int64, error) {
	uid, err := get(gid)
	if errors.Is(err, e) {
		id, err := a.model.Users.Random(gid)
		if err != nil {
			return 0, err
		}
		if err := insert(gid, id); err != nil {
			return 0, err
		}
		uid = id
	} else if err != nil {
		return 0, err
	}
	return uid, nil
}

func (a *App) getDailyPair(gid int64) (int64, int64, error) {
	x, y, err := a.model.Pairs.Get(gid)
	if errors.Is(err, model.ErrNoPair) {
		pair, err := a.model.Users.NRandom(gid, 2)
		if err != nil {
			return 0, 0, err
		}
		if len(pair) != 2 {
			return 0, 0, model.ErrNoPair
		}
		x = pair[0]
		y = pair[1]
		if err := a.model.Pairs.Insert(gid, x, y); err != nil {
			return 0, 0, err
		}
	} else if err != nil {
		return 0, 0, err
	}
	return x, y, nil
}

func (a *App) fetchAndSend(c tele.Context, url string) error {
	pic, err := fetchPicture(url)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

// sendRandomFile sends a random file from dir.
func sendRandomFile(c tele.Context, dir string) error {
	return sendRandomFileWith(c, dir, randomFile)
}

// sendRandomFileWith sends a random file chosen by f from dir.
func sendRandomFileWith(c tele.Context, dir string, f randomFileFunc) error {
	path, err := f(dir)
	if err != nil {
		return err
	}
	return sendFile(c, path)
}

// sendFile sends a file located at path.
func sendFile(c tele.Context, path string) error {
	return c.Send(&tele.Photo{File: tele.FromDisk(path)})
}

type randomFileFunc func(dir string) (string, error)

// randomFile returns a random filename from a directory.
func randomFile(dir string) (string, error) {
	fs, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	f := fs[rand.Intn(len(fs))]
	return filepath.Join(dir, f.Name()), nil
}

// randomFileFromHierarchy returns a random filename from a hierarchy of directories.
func randomFileFromHierarchy(root string) (string, error) {
	var filenames []string
	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() {
			filenames = append(filenames, path)
		}
		return nil
	}); err != nil {
		return "", err
	}
	return filenames[rand.Intn(len(filenames))], nil
}

func sendSmallProfilePhoto(c tele.Context) error {
	user, err := c.Bot().ChatByID(c.Sender().ID)
	if err != nil {
		return err
	}
	file, err := c.Bot().FileByID(user.Photo.SmallFileID)
	if err != nil {
		return err
	}
	f, err := c.Bot().File(&file)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.Send(&tele.Photo{File: tele.FromReader(f)})
}

func sendLargeProfilePhoto(c tele.Context) error {
	ps, err := c.Bot().ProfilePhotosOf(c.Sender())
	if err != nil {
		return err
	}
	if len(ps) < 1 {
		return nil
	}
	return c.Send(&ps[0])

}
