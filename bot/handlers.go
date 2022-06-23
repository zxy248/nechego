package bot

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
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

	"golang.org/x/exp/slices"
	tele "gopkg.in/telebot.v3"
)

const dataPath = "data"

// handleProbability responds with the probability of the message.
func (b *Bot) handleProbability(c tele.Context) error {
	a := getMessage(c).Argument()
	return c.Send(probability(a))
}

// handleWho responds with the message appended to the random chat member.
func (b *Bot) handleWho(c tele.Context) error {
	argument := getMessage(c).Argument()
	gid := c.Chat().ID

	uid, err := b.users.Random(gid)
	if err != nil {
		return err
	}

	m, err := b.chatMember(gid, uid)
	if err != nil {
		return err
	}

	name := markdownEscaper.Replace(chatMemberName(m))
	text := markdownEscaper.Replace(argument)
	return c.Send(who(uid, name, text), tele.ModeMarkdownV2)
}

const catURL = "https://thiscatdoesnotexist.com/"

// handleCat sends a picture of a cat.
func (b *Bot) handleCat(c tele.Context) error {
	pic, err := fetchPicture(catURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

const (
	nameTooLong   = "Максимальная длина имени 16 символов"
	yourName      = "Ваше имя: *%s* 🔖"
	pleaseReEnter = "Для использования этой функции Вам необходимо перезайти в беседу"
	nameSet       = "Имя *%v* установлено ✅"
)

// handleTitle sets the admin title of the sender.
func (b *Bot) handleTitle(c tele.Context) error {
	group := c.Chat()
	gid := group.ID
	uid := c.Sender().ID
	title := getMessage(c).Argument()

	if utf8.RuneCountInString(title) > 16 {
		return c.Send(makeError(nameTooLong))
	}

	m, err := b.chatMember(gid, uid)
	if err != nil {
		return err
	}
	if title == "" {
		name := markdownEscaper.Replace(chatMemberName(m))
		return c.Send(fmt.Sprintf(yourName, name), tele.ModeMarkdownV2)
	}
	if err := c.Bot().SetAdminTitle(group, c.Sender(), title); err != nil {
		return c.Send(makeError(pleaseReEnter))
	}
	return c.Send(fmt.Sprintf(nameSet, markdownEscaper.Replace(title)), tele.ModeMarkdownV2)
}

const animeFormat = "https://thisanimedoesnotexist.ai/results/psi-%s/seed%s.png"

var animePsis = []string{"0.3", "0.4", "0.5", "0.6", "0.7", "0.8", "0.9", "1.0",
	"1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "2.0"}

// handleAnime sends an anime picture.
func (b *Bot) handleAnime(c tele.Context) error {
	psi := animePsis[rand.Intn(len(animePsis))]
	seed := randomNumbers(5)
	url := fmt.Sprintf(animeFormat, psi, seed)
	pic, err := fetchPicture(url)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

const furFormat = "https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%s.jpg"

// handleFurry sends a furry picture.
func (b *Bot) handleFurry(c tele.Context) error {
	seed := randomNumbers(5)
	url := fmt.Sprintf(furFormat, seed)
	pic, err := fetchPicture(url)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

const flagFormat = "https://thisflagdoesnotexist.com/images/%d.png"

// handleFlag sends a picture of a flag.
func (b *Bot) handleFlag(c tele.Context) error {
	seed := rand.Intn(5000)
	url := fmt.Sprintf(flagFormat, seed)
	pic, err := fetchPicture(url)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

const personURL = "https://thispersondoesnotexist.com/image"

// handlePerson sends a picture of a person.
func (b *Bot) handlePerson(c tele.Context) error {
	pic, err := fetchPicture(personURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

const horseURL = "https://thishorsedoesnotexist.com/"

// handleHorse sends a picture of a horse.
func (b *Bot) handleHorse(c tele.Context) error {
	pic, err := fetchPicture(horseURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

const artURL = "https://thisartworkdoesnotexist.com/"

// handleArt sends a picture of an art.
func (b *Bot) handleArt(c tele.Context) error {
	pic, err := fetchPicture(artURL)
	if err != nil {
		return err
	}
	return c.Send(pic)
}

const carURL = "https://www.thisautomobiledoesnotexist.com/"

var carImageRe = regexp.MustCompile(
	"<img id = \"vehicle\" src=\"data:image/png;base64,(.+)\" class=\"center\">")

// handleCar sends a picture of a car.
func (b *Bot) handleCar(c tele.Context) error {
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
	b64 := ss[1]
	img, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	return c.Send(photoFromBytes(img))
}

const pairOfTheDayFormat = "Пара дня ✨\n%s 💘 %s"

// handlePair sends the current pair of the day, randomly choosing a new pair if
// needed.
func (b *Bot) handlePair(c tele.Context) error {
	gid := c.Chat().ID

	uidx, uidy, err := b.pairs.Get(gid)
	if errors.Is(err, model.ErrNoPair) {
		x, err := b.users.Random(gid)
		if err != nil {
			return err
		}
		y, err := b.users.Random(gid)
		if err != nil {
			return err
		}

		if x == y {
			return c.Send("💔")
		}

		if err := b.pairs.Insert(gid, x, y); err != nil {
			return err
		}

		uidx = x
		uidy = y
	} else if err != nil {
		return err
	}

	mx, err := b.chatMember(gid, uidx)
	if err != nil {
		return err
	}
	my, err := b.chatMember(gid, uidy)
	if err != nil {
		return err
	}
	namex := markdownEscaper.Replace(chatMemberName(mx))
	namey := markdownEscaper.Replace(chatMemberName(my))
	return c.Send(fmt.Sprintf(pairOfTheDayFormat,
		mention(uidx, namex), mention(uidy, namey)), tele.ModeMarkdownV2)
}

const eblanOfTheDayFormat = "Еблан дня: %s 😸"

// handleEblan sends the current eblan of the day, randomly choosing a new one if needed.
func (b *Bot) handleEblan(c tele.Context) error {
	gid := c.Chat().ID
	uid, err := b.getDaily(gid, b.eblans.Get, b.eblans.Insert, model.ErrNoEblan)
	if err != nil {
		return err
	}
	m, err := b.chatMember(gid, uid)
	if err != nil {
		return err
	}
	name := markdownEscaper.Replace(chatMemberName(m))
	return c.Send(fmt.Sprintf(eblanOfTheDayFormat, mention(uid, name)), tele.ModeMarkdownV2)
}

const adminOfTheDayFormat = "Админ дня: %s 👑"

func (b *Bot) handleAdmin(c tele.Context) error {
	gid := c.Chat().ID
	uid, err := b.getDaily(gid, b.admins.GetDaily, b.admins.InsertDaily, model.ErrNoAdmin)
	if err != nil {
		return err
	}
	m, err := b.chatMember(gid, uid)
	if err != nil {
		return err
	}
	name := markdownEscaper.Replace(chatMemberName(m))
	return c.Send(fmt.Sprintf(adminOfTheDayFormat, mention(uid, name)), tele.ModeMarkdownV2)
}

const masyunyaStickersName = "masyunya_vk"

// handleMasyunya sends a random sticker of Masyunya.
func (b *Bot) handleMasyunya(c tele.Context) error {
	ss, err := c.Bot().StickerSet(masyunyaStickersName)
	if err != nil {
		return err
	}
	s := ss.Stickers[rand.Intn(len(ss.Stickers))]
	return c.Send(&s)
}

var poppyStickersNames = []string{"pappy2_vk", "poppy_vk"}

func (b *Bot) handlePoppy(c tele.Context) error {
	var stickers []tele.Sticker
	for _, sn := range poppyStickersNames {
		ss, err := c.Bot().StickerSet(sn)
		if err != nil {
			return err
		}
		stickers = append(stickers, ss.Stickers...)
	}
	s := stickers[rand.Intn(len(stickers))]
	return c.Send(&s)
}

const helloChance = 0.2

// handleHello sends a hello sticker
func (b *Bot) handleHello(c tele.Context) error {
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
func (b *Bot) handleMouse(c tele.Context) error {
	return c.Send(mouseVideo)
}

const weatherTimeout = 10 * time.Second
const weatherURL = "https://wttr.in/"
const weatherFormat = `?format=%l:+%c+%t+\nОщущается+как+%f\n\nВетер+—+%w\nВлажность+—+%h\nДавление+—+%P\nФаза+луны+—+%m\nУФ-индекс+—+%u\n`

// handleWeather sends the current weather for a given city
func (b *Bot) handleWeather(c tele.Context) error {
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

func (b *Bot) handleTikTok(c tele.Context) error {
	return c.Send(tikTokVideo)
}

const (
	listTemplate = `Список %s 📝
%s`
	listLength = 5
)

func (b *Bot) handleList(c tele.Context) error {
	argument := markdownEscaper.Replace(getMessage(c).Argument())
	gid := c.Chat().ID

	var uids []int64
	for i := 0; i < 5; i++ {
		uid, err := b.users.Random(gid)
		if err != nil {
			return err
		}
		if !slices.Contains(uids, uid) {
			uids = append(uids, uid)
		}
	}

	var list string
	for _, uid := range uids {
		m, err := b.chatMember(gid, uid)
		if err != nil {
			return err
		}
		name := markdownEscaper.Replace(chatMemberName(m))
		list = list + "— " + mention(uid, name) + "\n"
	}
	return c.Send(fmt.Sprintf(listTemplate, argument, list), tele.ModeMarkdownV2)
}

const (
	numberedTopTemplate = `Топ %d %s 🏆
%s`
	unnumberedTopTemplate = `Топ %s 🏆
%s`
	maxTopNumber = 5
)

func (b *Bot) handleTop(c tele.Context) error {
	gid := c.Chat().ID
	a, err := getMessage(c).DynamicArgument()
	if err != nil {
		return err
	}
	argument, ok := a.(input.TopArgument)
	if !ok {
		return errors.New("the argument is not a TopArgument")
	}

	uids, err := b.users.List(gid)
	if err != nil {
		return err
	}
	rand.Shuffle(len(uids), func(i, j int) {
		uids[i], uids[j] = uids[j], uids[i]
	})

	var n int
	if argument.NumberPresent {
		n = argument.Number
	} else {
		if len(uids) > maxTopNumber {
			n = rand.Intn(maxTopNumber) + 1
		} else {
			n = rand.Intn(len(uids)) + 1
		}
	}

	if n < 1 || n > len(uids) || n > maxTopNumber {
		return c.Send(errorSign())
	}
	uids = uids[:n]

	var list string
	for i, uid := range uids {
		m, err := b.chatMember(gid, uid)
		if err != nil {
			return err
		}
		name := markdownEscaper.Replace(chatMemberName(m))
		list = list + fmt.Sprintf("_%d\\._ %s\n", i+1, mention(uid, name))
	}

	s := markdownEscaper.Replace(argument.String)
	var result string
	if argument.NumberPresent {
		result = fmt.Sprintf(numberedTopTemplate, n, s, list)
	} else {
		result = fmt.Sprintf(unnumberedTopTemplate, s, list)
	}
	return c.Send(result, tele.ModeMarkdownV2)
}

var (
	albumsPath     = filepath.Join(dataPath, "vk.com-albums")
	basiliCatsPath = filepath.Join(albumsPath, "basili")
	casperPath     = filepath.Join(albumsPath, "casper")
	zeusPath       = filepath.Join(albumsPath, "zeus")
	picPath        = filepath.Join(albumsPath, "saved")
)

func (b *Bot) handleBasili(c tele.Context) error {
	path, err := randomFilename(basiliCatsPath)
	if err != nil {
		return err
	}
	return c.Send(&tele.Photo{File: tele.FromDisk(path)})
}

func (b *Bot) handleCasper(c tele.Context) error {
	path, err := randomFilename(casperPath)
	if err != nil {
		return err
	}
	return c.Send(&tele.Photo{File: tele.FromDisk(path)})
}

func (b *Bot) handleZeus(c tele.Context) error {
	path, err := randomFilename(zeusPath)
	if err != nil {
		return err
	}
	return c.Send(&tele.Photo{File: tele.FromDisk(path)})
}

func (b *Bot) handlePic(c tele.Context) error {
	path, err := randomFilename(picPath)
	if err != nil {
		return err
	}
	return c.Send(&tele.Photo{File: tele.FromDisk(path)})
}

func (b *Bot) handleDice(c tele.Context) error {
	return c.Send(tele.Cube)
}

func (b *Bot) handleGame(c tele.Context) error {
	games := []*tele.Dice{tele.Dart, tele.Ball, tele.Goal, tele.Slot, tele.Bowl}
	game := games[rand.Intn(len(games))]
	return c.Send(game)
}

const randomPhotoChance = 0.02

func (b *Bot) handleRandomPhoto(c tele.Context) error {
	r := rand.Float64()
	if r <= randomPhotoChance {
		// Alternative version: sends a large photo.
		//
		// ps, err := c.Bot().ProfilePhotosOf(c.Sender())
		// if err != nil {
		// 	return err
		// }
		// if len(ps) < 1 {
		// 	return nil
		// }
		// return c.Send(&ps[0])

		user, err := c.Bot().ChatByID(c.Sender().ID)
		if err != nil {
			return err
		}
		file, err := b.bot.FileByID(user.Photo.SmallFileID)
		if err != nil {
			return err
		}
		f, err := b.bot.File(&file)
		if err != nil {
			return err
		}
		defer f.Close()
		return c.Send(&tele.Photo{File: tele.FromReader(f)})
	}
	return nil
}

// handleKeyboardOpen opens the keyboard.
func (b *Bot) handleKeyboardOpen(c tele.Context) error {
	return c.Send("Клавиатура ⌨️", b.keyboard)
}

// handleKeyboardClose closes the keyboard.
func (b *Bot) handleKeyboardClose(c tele.Context) error {
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
func (b *Bot) handleTurnOn(c tele.Context) error {
	emoji := emojisActive[rand.Intn(len(emojisActive))]
	gid := c.Chat().ID
	b.status.Enable(gid)
	return c.Send(fmt.Sprintf(botTurnedOn, emoji))
}

// handleTurnOff turns the bot off.
func (b *Bot) handleTurnOff(c tele.Context) error {
	emoji := emojisInactive[rand.Intn(len(emojisInactive))]
	gid := c.Chat().ID
	b.status.Disable(gid)
	return c.Send(fmt.Sprintf(botTurnedOff, emoji), tele.RemoveKeyboard)
}

const (
	userBlocked          = "Пользователь заблокирован 🚫"
	userAlreadyBlocked   = "Пользователь уже заблокирован 🛑"
	userUnblocked        = "Пользователь разблокирован ✅"
	userAlreadyUnblocked = "Пользователь не заблокирован ❎"
)

// handleBan adds the user ID of the reply message's sender to the ban list.
func (b *Bot) handleBan(c tele.Context) error {
	uid := c.Message().ReplyTo.Sender.ID
	banned, err := b.bans.Banned(uid)
	if err != nil {
		return err
	}
	if banned {
		return c.Send(userAlreadyBlocked)
	}

	if err := b.bans.Ban(uid); err != nil {
		return err
	}
	return c.Send(userBlocked)
}

// handleUnban removes the user ID of the reply message's sender from the ban list.
func (b *Bot) handleUnban(c tele.Context) error {
	uid := c.Message().ReplyTo.Sender.ID
	banned, err := b.bans.Banned(uid)
	if err != nil {
		return err
	}
	if !banned {
		return c.Send(userAlreadyUnblocked)
	}

	if err := b.bans.Unban(uid); err != nil {
		return err
	}
	return c.Send(userUnblocked)
}

const infoTemplate = `ℹ️ *Информация* 📌

👤 _Администрация_
%s
🛑 _Черный список_
%s
🔒 _Запрещенные команды_
%s
`

// handleInfo sends a list of useful information.
func (b *Bot) handleInfo(c tele.Context) error {
	gid := c.Chat().ID
	l, err := b.admins.List(gid)
	if err != nil {
		return err
	}

	var admins string
	for _, uid := range l {
		m, err := b.chatMember(gid, uid)
		if err != nil {
			return err
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

	l, err = b.bans.List()
	if err != nil {
		return err
	}

	var banned string
	for _, uid := range l {
		m, err := b.chatMember(gid, uid)
		if err != nil {
			return err
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

	forbiddenCommands, err := b.forbid.List(gid)
	if err != nil {
		return err
	}
	var forbiddenList string
	for _, c := range forbiddenCommands {
		t := markdownEscaper.Replace(input.CommandText(c))
		forbiddenList += "— " + t + "\n"
	}
	if forbiddenList == "" {
		forbiddenList = "…\n"
	}

	list := fmt.Sprintf(infoTemplate, admins, banned, forbiddenList)
	return c.Send(list, tele.ModeMarkdownV2)
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

func (b *Bot) handleHelp(c tele.Context) error {
	return c.Send(help, tele.ModeMarkdownV2)
}

func (b *Bot) handleJoin(c tele.Context) error {
	group := c.Chat()
	gid := group.ID
	uid := c.Message().UserJoined.ID

	m, err := b.chatMember(gid, uid)
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
func (b *Bot) handleForbid(c tele.Context) error {
	return b.handleCommandAction(c, func(command input.Command) error {
		if isCommandForbidden(c) {
			return c.Send(commandAlreadyForbidden)
		}
		if err := b.forbid.Forbid(c.Chat().ID, command); err != nil {
			return err
		}
		return c.Send(commandForbidden)
	})
}

// handlePermit permits a command.
func (b *Bot) handlePermit(c tele.Context) error {
	return b.handleCommandAction(c, func(command input.Command) error {
		if !isCommandForbidden(c) {
			return c.Send(commandAlreadyPermitted)
		}
		if err := b.forbid.Permit(c.Chat().ID, command); err != nil {
			return err
		}
		return c.Send(commandPermitted)
	})
}

// handleCommandAction performs an action on a command.
func (b *Bot) handleCommandAction(c tele.Context, a func(input.Command) error) error {
	arg, err := getMessage(c).DynamicArgument()
	if err != nil {
		if errors.Is(err, input.ErrNoCommand) {
			return c.Send(makeError("Укажите команду"))
		}
		if errors.Is(err, input.ErrUnknownCommand) {
			return c.Send(makeError("Неизвестная команда"))
		}
		return err
	}
	command, ok := arg.(input.Command)
	if !ok {
		return errors.New("not a command")
	}
	return a(command)
}

func handleNothing(c tele.Context) error {
	return nil
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

func (b *Bot) isGroupMember(group tele.Recipient, user tele.Recipient) bool {
	member, err := b.bot.ChatMemberOf(group, user)
	if err != nil || member.Role == tele.Kicked || member.Role == tele.Left {
		return false
	}
	return true
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

var markdownEscaper = newMarkdownEscaper()

// newMarkdownEscaper creates a new Markdown replacer. The replacer
// escapes any character with the code between 1 and 126 inclusively
// with a preceding backslash.
func newMarkdownEscaper() *strings.Replacer {
	var table []string
	for i := 1; i <= 126; i++ {
		c := string(rune(i))
		table = append(table, c, "\\"+c)
	}
	return strings.NewReplacer(table...)
}

func randomFilename(path string) (string, error) {
	ds, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	d := ds[rand.Intn(len(ds))]
	return filepath.Join(path, d.Name()), nil
}

func (b *Bot) chatMember(gid, uid int64) (*tele.ChatMember, error) {
	group, err := b.bot.ChatByID(gid)
	if err != nil {
		return nil, err
	}
	member, err := b.bot.ChatMemberOf(group, tele.ChatID(uid))
	if err != nil {
		return nil, err
	}
	if !chatMemberPresent(member) {
		b.users.Delete(gid, uid)
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

type getter func(gid int64) (int64, error)
type inserter func(gid, uid int64) error

func (b *Bot) getDaily(gid int64, get getter, insert inserter, e error) (int64, error) {
	uid, err := get(gid)
	if errors.Is(err, e) {
		id, err := b.users.Random(gid)
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
