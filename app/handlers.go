package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"nechego/input"
	"nechego/model"
	"net/http"
	"net/url"
	"path/filepath"
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
	u, err := a.model.RandomUser(getGroup(c))
	if err != nil {
		return internalError(c, err)
	}
	message := markdownEscaper.Replace(getMessage(c).Argument())
	return c.Send(fmt.Sprintf("%s %s", a.mustMentionUser(u), message), tele.ModeMarkdownV2)
}

const (
	maxNameLength = 16
	nameTooLong   = "Максимальная длина имени 16 символов."
	yourNameIs    = "Ваше имя: *%s* 🔖"
	pleaseReEnter = "Перезайдите в беседу чтобы использовать эту функцию."
	nameSet       = "Имя *%s* установлено ✅"
)

// handleTitle sets the admin title of the sender.
func (a *App) handleTitle(c tele.Context) error {
	user := getUser(c)
	newName := getMessage(c).Argument()
	if newName == "" {
		return c.Send(fmt.Sprintf(yourNameIs, a.mustMentionUser(user)), tele.ModeMarkdownV2)
	}
	if utf8.RuneCountInString(newName) > maxNameLength {
		return userError(c, nameTooLong)
	}
	if err := setName(c, user, newName); err != nil {
		return userError(c, pleaseReEnter)
	}
	return c.Send(fmt.Sprintf(nameSet, markdownEscaper.Replace(newName)), tele.ModeMarkdownV2)
}

func setName(c tele.Context, u model.User, newName string) error {
	group := c.Chat()
	sender := c.Sender()
	return c.Bot().SetAdminTitle(group, sender, newName)
}

var (
	mouseVideoPath = filepath.Join(dataPath, "mouse.mp4")
	mouseVideo     = &tele.Video{File: tele.FromDisk(mouseVideoPath)}
)

// handleMouse sends the mouse video
func (a *App) handleMouse(c tele.Context) error {
	return c.Send(mouseVideo)
}

const (
	weatherTimeout      = 10 * time.Second
	weatherTimeoutError = "Время запроса вышло ☔️"
	placeNotExists      = "Такого места не существует ☔️"
	weatherBadRequest   = "Неудачный запрос ☔️"
	weatherURL          = "https://wttr.in/"
	weatherFormat       = "?format=%l:+%c+%t+\nОщущается+как+%f\n\nВетер+—+%w\nВлажность+—+%h\nДавление+—+%P\nФаза+луны+—+%m\nУФ-индекс+—+%u\n"
)

// handleWeather sends the current weather for a given city
func (a *App) handleWeather(c tele.Context) error {
	place := getMessage(c).Argument()
	r, err := fetchWeather(place)
	if err != nil {
		if err.(*url.Error).Timeout() {
			return userError(c, weatherTimeoutError)
		}
		return internalError(c, err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		if r.StatusCode == http.StatusNotFound {
			return userError(c, placeNotExists)
		}
		return userError(c, weatherBadRequest)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return internalError(c, err)
	}
	return c.Send(string(data))
}

func fetchWeather(place string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), weatherTimeout)
	defer cancel()

	url := weatherURL + place + weatherFormat
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept-Language", "ru")
	return http.DefaultClient.Do(req)
}

var tikTokVideo = &tele.Video{File: tele.FromDisk("data/tiktok.mp4")}

// !тикток
func (a *App) handleTikTok(c tele.Context) error {
	return c.Send(tikTokVideo)
}

const (
	handleListTemplate = "Список %s 📝\n%s"
	minListLength      = 3
	maxListLength      = 5
)

// !список
func (a *App) handleList(c tele.Context) error {
	n := randInRange(minListLength, maxListLength)
	users, err := a.model.RandomUsers(getGroup(c), n)
	if err != nil {
		return internalError(c, err)
	}
	what := markdownEscaper.Replace(getMessage(c).Argument())
	out := fmt.Sprintf(handleListTemplate, what, a.formatUnorderedList(users))
	return c.Send(out, tele.ModeMarkdownV2)
}

const (
	numberedTopTemplate   = "Топ %d %s 🏆\n%s"
	unnumberedTopTemplate = "Топ %s 🏆\n%s"
	maxTopNumber          = 5
	badTopNumber          = "Некорректное число."
)

// !топ
func (a *App) handleTop(c tele.Context) error {
	argument, err := getMessage(c).TopArgument()
	if err != nil {
		return internalError(c, err)
	}

	var number int
	if argument.Number != nil {
		number = *argument.Number
	} else {
		number = maxTopNumber
	}
	if number <= 0 || number > maxTopNumber {
		return userError(c, badTopNumber)
	}
	users, err := a.model.RandomUsers(getGroup(c), number)
	if err != nil {
		return internalError(c, err)
	}
	top := a.formatOrderedList(users)
	what := markdownEscaper.Replace(argument.String)
	var out string
	if argument.Number != nil {
		out = fmt.Sprintf(numberedTopTemplate, number, what, top)
	} else {
		out = fmt.Sprintf(unnumberedTopTemplate, what, top)
	}
	return c.Send(out, tele.ModeMarkdownV2)
}

var games = []*tele.Dice{tele.Dart, tele.Ball, tele.Goal, tele.Slot, tele.Bowl}

// !игра
func (a *App) handleGame(c tele.Context) error {
	game := games[rand.Intn(len(games))]
	return c.Send(game)
}

const randomPhotoChance = 0.02

func (a *App) handleRandomPhoto(c tele.Context) error {
	if rand.Float64() <= randomPhotoChance {
		return sendSmallProfilePhoto(c)
	}
	return nil
}

// !открыть
func (a *App) handleKeyboardOpen(c tele.Context) error {
	return c.Send("Клавиатура ⌨️", keyboard)
}

// !закрыть
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

// !включить
func (a *App) handleTurnOn(c tele.Context) error {
	emoji := emojisActive[rand.Intn(len(emojisActive))]
	ok := a.model.EnableGroup(getGroup(c))
	if !ok {
		return c.Send(fmt.Sprintf(botAlreadyTurnedOn, emoji))
	}
	return c.Send(fmt.Sprintf(botTurnedOn, emoji))
}

// !выключить
func (a *App) handleTurnOff(c tele.Context) error {
	emoji := emojisInactive[rand.Intn(len(emojisInactive))]
	ok := a.model.DisableGroup(getGroup(c))
	if !ok {
		return c.Send(fmt.Sprintf(botAlreadyTurnedOff, emoji))
	}
	return c.Send(fmt.Sprintf(botTurnedOff, emoji), tele.RemoveKeyboard)
}

const (
	userBlocked          = "Пользователь заблокирован 🚫"
	userAlreadyBlocked   = "Пользователь уже заблокирован 🛑"
	userUnblocked        = "Пользователь разблокирован ✅"
	userAlreadyUnblocked = "Пользователь не заблокирован ❎"
)

// !бан
func (a *App) handleBan(c tele.Context) error {
	user := getReplyUser(c)
	if user.Banned {
		return c.Send(userAlreadyBlocked)
	}
	a.model.BanUser(user)
	return c.Send(userBlocked)
}

// handleUnban removes the user ID of the reply message's sender from the ban list.
func (a *App) handleUnban(c tele.Context) error {
	user := getReplyUser(c)
	if !user.Banned {
		return c.Send(userAlreadyUnblocked)
	}
	a.model.UnbanUser(user)
	return c.Send(userUnblocked)
}

const infoTemplate = "ℹ️ *Информация* 📌\n\n%s\n%s\n%s\n"

// handleInfo sends a few lists of useful information.
func (a *App) handleInfo(c tele.Context) error {
	group := getGroup(c)
	admins, err := a.adminList(group)
	if err != nil {
		return internalError(c, err)
	}
	bans, err := a.banList(group)
	if err != nil {
		return internalError(c, err)
	}
	commands, err := a.forbiddenCommandList(group)
	if err != nil {
		return internalError(c, err)
	}
	lists := fmt.Sprintf(infoTemplate, admins, bans, commands)
	return c.Send(lists, tele.ModeMarkdownV2)
}

const adminListTemplate = "👤 _Администрация_\n%s"

func (a *App) adminList(g model.Group) (string, error) {
	users, err := a.model.ListUsers(g)
	if err != nil {
		return "", err
	}
	admins := []model.User{}
	for _, u := range users {
		if u.Admin {
			admins = append(admins, u)
		}
	}
	return fmt.Sprintf(adminListTemplate, a.formatUnorderedList(admins)), nil
}

const banListTemplate = "🛑 _Черный список_\n%s"

func (a *App) banList(g model.Group) (string, error) {
	users, err := a.model.ListUsers(g)
	if err != nil {
		return "", err
	}
	banned := []model.User{}
	for _, u := range users {
		if u.Banned {
			banned = append(banned, u)
		}
	}
	return fmt.Sprintf(banListTemplate, a.formatUnorderedList(banned)), nil
}

const forbiddenCommandListTemplate = "🔒 _Запрещенные команды_\n%s"

func (a *App) forbiddenCommandList(g model.Group) (string, error) {
	commands, err := a.model.ForbiddenCommands(g)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(forbiddenCommandListTemplate, formatCommandList(commands)), nil
}

var categories = map[string]string{
	"кредиты": `
С помощью кредитов вы можете взять деньги в долг, чтобы потратить их на то, чего вам так не хватает: например, закинуть побольше денег в кости или купить удочку\.

Непогашенный кредит вас ограничивает:
_\(1\)_ Вы не можете списывать деньги с банковского счета\.
_\(2\)_ На вашу силу накладывается дебаф\.

Максимальное количество средств, которые вы можете взять в кредит — _кредитный лимит_ — это самая большая сумма денег, которой вы когда либо владели\.

Если эти условия вас устраивают, кредит можно взять так:
` + "`" + `!кредит <сумма>` + "`" + `

Чтобы погасить кредит, сначала пополните свой банковский счет:
_1\._ ` + "`" + `!депозит <сумма>` + "`" + `
_2\._ ` + "`" + `!погасить <сумма>` + "`" + `

Не обязательно отдавать кредит целиком — его можно отдавать порциями\.`,
	"казино": `
Если хотите попытать удачу в казино, вы можете рискнуть бросить кости\.

Используйте команду ` + "`" + `!кости <ставка>` + "`" + `\. Робот\-прислужник бросит кости, затем бросайте свои\. Если вы выбьете больше, чем он, вы получите выигрыш\. В случае ничьей, вам возвращается ставка\. Если робот выбьет больше, чем вы — он заберет ваши деньги\.

С некоторым шансом вы можете получить бонус\. Минимальная ставка для получения бонуса — ` + formatMoney(diceBetForBonus) + `\.
`,
	"драка": `
Вы можете сражаться с другими пользователями\. За победу вы получите монеты, за проигрыш — потеряете\. Чтобы вызвать пользователя на бой, используйте команду ` + "`" + `!драка` + "`" + ` в ответ на его сообщение\.

На исход поединка влияет сила\. Проверить ее можно с помощью команды ` + "`" + `!сила` + "`" + `\.

Список самых сильных пользователей: ` + "`" + `!топ сильных` + "`" + `
`,
	"рыбалка": `
Вы можете ловить рыбу: ` + "`" + `!рыбалка` + "`" + `\. Для этого вам потребуется удочка\. Приобрести ее можно за ` + formatMoney(fishingRodPrice) + ` с помощью команды ` + "`" + `!удочка` + "`" + `\.

Если вы оставили рыбу себе, ее можно съесть, чтобы восстановить немного энергии: ` + "`" + `!еда` + "`" + `\.
`,
	"банк": `
Чтобы безопасно хранить деньги — так, чтобы у вас не могли их отнять — воспользуйтесь услугами банка\.

Информация о счете, коммиссиях, и лимитах: ` + "`" + `!банк` + "`" + `
Проверить баланс: ` + "`" + `!баланс` + "`" + `
Пополнить счет: ` + "`" + `!депозит <сумма>` + "`" + `

Вы не можете тратить средства из банка напрямую\. Можно использовать только то, что есть у вас в кошельке\.

Снять средства: ` + "`" + `!обнал <сумма>` + "`" + `
`,
	"экономика": `
Вы можете зарабатывать, тратить, передавать и хранить деньги — монеты 💰\.

Передать деньги другому пользователю можно с помощью команды ` + "`" + `!перевод <сумма>` + "`" + `\.

Подзаработать монет можно с помощью драки, игры в кости и рыбалки\. Также, вы можете взять деньги в кредит\.

Вы можете узнать статистику об экономическом состоянии беседы:
— ` + "`" + `!капитал` + "`" + `
— ` + "`" + `!профиль` + "`" + `
— ` + "`" + `!топ богатых` + "`" + `
— ` + "`" + `!топ бедных` + "`" + `

_См\. также: ` + "`" + `банк` + "`" + `, ` + "`" + `кредиты` + "`" + `\._
`,
	"нейросети": `
Картинки, сгенерированные компьютерными машинами\.

· ` + "`" + `!кот` + "`" + `
· ` + "`" + `!аниме` + "`" + `
· ` + "`" + `!фурри` + "`" + `
· ` + "`" + `!флаг` + "`" + `
· ` + "`" + `!чел` + "`" + `
· ` + "`" + `!лошадь` + "`" + `
· ` + "`" + `!арт` + "`" + `
· ` + "`" + `!авто` + "`" + `
`,
	"управление": `
Сервисные или информационные команды\.

Документация: ` + "`" + `!информация` + "`" + ` _или_ ` + "`" + `!команды` + "`" + `\.
Клавиатура: ` + "`" + `!открыть` + "`" + `, ` + "`" + `!закрыть` + "`" + `\.
Смена имени: ` + "`" + `!имя <новое имя>` + "`" + `\.

Если робот\-прислужник «мы за фашизм» вам надоел, вы можете воспользоваться командой ` + "`" + `!выключить` + "`" + `\.
Чтобы реактивировать его, используйте команду ` + "`" + `!включить` + "`" + `\.
`,
	"кошки": `
Кошки пользователей беседы «нечего»\.

· ` + "`" + `!марсик` + "`" + `
· ` + "`" + `!муся` + "`" + `
· ` + "`" + `!каспер` + "`" + `
· ` + "`" + `!зевс` + "`" + `
`,
	"администрирование": `
Если вы администратор, вы можете управлять командами и пользователями\.

Заблокировать/разблокировать команду:
` + "`" + `!запретить <команда>` + "`" + `
` + "`" + `!разрешить <команда>` + "`" + `

Забанить/разбанить пользователя:
` + "`" + `!бан` + "`" + ` или ` + "`" + `!разбан` + "`" + ` в ответ на сообщение пользователя\.
`,
}

const help = `
· ` + "`" + `\!инфа` + "`" + `
· ` + "`" + `\!кто` + "`" + `
· ` + "`" + `\!список` + "`" + `
· ` + "`" + `\!топ` + "`" + `
· ` + "`" + `\!погода` + "`" + `
· ` + "`" + `\!пара дня` + "`" + `
· ` + "`" + `\!еблан дня` + "`" + `
· ` + "`" + `\!админ дня` + "`" + `
· ` + "`" + `\!мыш` + "`" + `
· ` + "`" + `\!тикток` + "`" + `
· ` + "`" + `\!масюня` + "`" + `
· ` + "`" + `\!паппи` + "`" + `
· ` + "`" + `\!игра` + "`" + `

Для просмотра документации по конкретному вопросу, используйте команду ` + "`" + `\!помощь <тема>` + "`" + `\.

Доступные темы: %s\.
`

func (a *App) handleHelp(c tele.Context) error {
	category := getMessage(c).Argument()
	if category != "" {
		out, ok := categories[category]
		if !ok {
			return userError(c, "Такой категории нет.")
		}
		return c.Send(prependHelpHeader(category, out), tele.ModeMarkdownV2)
	}
	out := prependHelpHeader("базовые команды", fmt.Sprintf(help, availableCategories()))
	return c.Send(out, tele.ModeMarkdownV2)
}

func availableCategories() string {
	cats := []string{}
	for c := range categories {
		cats = append(cats, fmt.Sprintf("`%s`", c))
	}
	return strings.Join(cats, ", ")
}

func prependHelpHeader(category, help string) string {
	return fmt.Sprintf("🗂 Помощь: *%s*\\.\n%s", category, help)
}

func (a *App) handleJoin(c tele.Context) error {
	u := c.Message().UserJoined
	m, err := c.Bot().ChatMemberOf(c.Chat(), u)
	if err != nil {
		return err
	}
	if err := promoteIfNotAdmin(c, m); err != nil {
		return err
	}
	return c.Send(helloSticker())
}

func promoteIfNotAdmin(c tele.Context, m *tele.ChatMember) error {
	if m.Role != tele.Administrator && m.Role != tele.Creator {
		m.Rights.CanBeEdited = true
		m.Rights.CanManageChat = true
		return c.Bot().Promote(c.Chat(), m)
	}
	return nil
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
		ok := a.model.ForbidCommand(getGroup(c), command)
		if !ok {
			return c.Send(commandAlreadyForbidden)
		}
		return c.Send(commandForbidden)
	})
}

// handlePermit permits a command.
func (a *App) handlePermit(c tele.Context) error {
	return a.handleCommandAction(c, func(command input.Command) error {
		ok := a.model.PermitCommand(getGroup(c), command)
		if !ok {
			return c.Send(commandAlreadyPermitted)
		}
		return c.Send(commandPermitted)
	})
}

const (
	specifyCommand = "Укажите команду."
	unknownCommand = "Неизвестная команда."
)

// handleCommandAction performs an action on a command.
func (a *App) handleCommandAction(c tele.Context, action func(input.Command) error) error {
	command, err := getMessage(c).CommandActionArgument()
	if err != nil {
		if errors.Is(err, input.ErrNoCommand) {
			return userError(c, specifyCommand)
		}
		if errors.Is(err, input.ErrUnknownCommand) {
			return userError(c, unknownCommand)
		}
		return internalError(c, err)
	}
	return action(command)
}
