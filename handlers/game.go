package handlers

import (
	"errors"
	"fmt"
	"html"
	"math/rand"
	"nechego/format"
	"nechego/game"
	"nechego/teleutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	tele "gopkg.in/telebot.v3"
)

type Save struct {
	Universe *game.Universe
}

var saveRe = regexp.MustCompile("^!сохран")

func (h *Save) Match(s string) bool {
	return saveRe.MatchString(s)
}

func (h *Save) Handle(c tele.Context) error {
	if err := h.Universe.SaveAll(); err != nil {
		return err
	}
	return c.Send("💾 Игра сохранена.")
}

type Name struct{}

var nameRe = regexp.MustCompile("!имя (.*)")

func (h *Name) Match(s string) bool {
	return nameRe.MatchString(s)
}

func (h *Name) Handle(c tele.Context) error {
	name := html.EscapeString(teleutil.Args(c, nameRe)[1])
	const maxlen = 16
	if utf8.RuneCountInString(name) > maxlen {
		return c.Send(fmt.Sprintf("⚠️ Максимальная длина имени %d символов.", maxlen))
	}
	if err := teleutil.Promote(c, teleutil.Member(c, c.Sender())); err != nil {
		return err
	}
	if err := c.Bot().SetAdminTitle(c.Chat(), c.Sender(), name); err != nil {
		return err
	}
	return c.Send(fmt.Sprintf("Имя <b>%s</b> установлено ✅", name), tele.ModeHTML)
}

type Inventory struct {
	Universe *game.Universe
}

var inventoryRe = regexp.MustCompile("^!(инвентарь|лут)")

func (h *Inventory) Match(s string) bool {
	return inventoryRe.MatchString(s)
}

func (h *Inventory) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	items := user.Inventory.List()
	warn := ""
	if user.Inventory.Count() > game.InventorySize {
		warn = " (!)"
	}
	head := fmt.Sprintf("<b>🗄 %s: Инвентарь <code>[%d/%d%s]</code></b>\n",
		teleutil.Mention(c, user), len(items), game.InventorySize, warn)
	list := format.Items(items)
	return c.Send(head+list, tele.ModeHTML)
}

type Catch struct {
	Universe *game.Universe
}

var catchRe = regexp.MustCompile("^!улов")

func (h *Catch) Match(s string) bool {
	return catchRe.MatchString(s)
}

func (h *Catch) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	head := fmt.Sprintf("<b>🐟 %s: Улов</b>\n", teleutil.Mention(c, user))
	list := format.Catch(user.Inventory.List())
	return c.Send(head+list, tele.ModeHTML)
}

type Drop struct {
	Universe *game.Universe
}

var dropRe = regexp.MustCompile("^!(выкинуть|выбросить|выложить|дроп|положить) (.*)")

func (h *Drop) Match(s string) bool {
	return dropRe.MatchString(s)
}

func (h *Drop) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	for _, key := range teleutil.NumArg(c, dropRe, 2) {
		item, ok := user.Inventory.ByKey(key)
		if !ok {
			return c.Send(fmt.Sprintf("🗄 Предмета %s нет в инвентаре.",
				format.Key(key)), tele.ModeHTML)
		}
		if !user.Inventory.Move(world.Floor, item) {
			return c.Send(fmt.Sprintf("♻ Вы не можете выбросить %s.",
				format.Item(item)), tele.ModeHTML)
		}
		c.Send(fmt.Sprintf("🚮 Вы выбросили %s.", format.Item(item)), tele.ModeHTML)
	}
	world.Floor.Retain(10)
	return nil
}

type Pick struct {
	Universe *game.Universe
}

var pickRe = regexp.MustCompile("^!(взять|подобрать|поднять) (.*)")

func (h *Pick) Match(s string) bool {
	return pickRe.MatchString(s)
}

func (h *Pick) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	if user.Inventory.Count() > game.InventoryCap {
		return c.Send(format.InventoryFull)
	}
	for _, key := range teleutil.NumArg(c, pickRe, 2) {
		item, ok := world.Floor.ByKey(key)
		if !ok {
			return c.Send(fmt.Sprintf("🗄 Предмета %s нет на полу.",
				format.Key(key)), tele.ModeHTML)
		}
		if !world.Floor.Move(user.Inventory, item) {
			return c.Send(fmt.Sprintf("♻ Вы не можете взять %s.",
				format.Item(item)), tele.ModeHTML)
		}
		c.Send(fmt.Sprintf("🫳 Вы взяли %s.", format.Item(item)), tele.ModeHTML)
	}
	return nil
}

type Floor struct {
	Universe *game.Universe
}

var floorRe = regexp.MustCompile("^!(пол|мусор|вещи|предметы)")

func (h *Floor) Match(s string) bool {
	return floorRe.MatchString(s)
}

func (h *Floor) Handle(c tele.Context) error {
	world, _ := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	head := "<b>🗃️ Предметы</b>\n"
	list := format.Items(world.Floor.List())
	return c.Send(head+list, tele.ModeHTML)
}

type Market struct {
	Universe *game.Universe
}

var marketRe = regexp.MustCompile("^!(магаз|шоп)")

func (h *Market) Match(s string) bool {
	return marketRe.MatchString(s)
}

func (h *Market) Handle(c tele.Context) error {
	world, _ := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	head := fmt.Sprintf("<b>%s</b>\n", world.Market)
	list := format.Products(world.Market.Products())
	return c.Send(head+list, tele.ModeHTML)
}

type NameMarket struct {
	Universe *game.Universe
}

var nameMarketRe = regexp.MustCompile("^!назвать магазин (.*)")

func (h *NameMarket) Match(s string) bool {
	return nameMarketRe.MatchString(s)
}

func (h *NameMarket) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()
	if !user.Admin() {
		return c.Send(format.AdminsOnly)
	}
	name := teleutil.Args(c, nameMarketRe)[1]
	if !world.Market.SetName(name) {
		return c.Send(format.BadMarketName)
	}
	return c.Send(format.MarketRenamed)
}

type Buy struct {
	Universe *game.Universe
}

var buyRe = regexp.MustCompile("^!купить (.*)")

func (h *Buy) Match(s string) bool {
	return buyRe.MatchString(s)
}

func (h *Buy) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	if user.Inventory.Count() > game.InventoryCap {
		return c.Send(format.InventoryFull)
	}
	for _, key := range teleutil.NumArg(c, buyRe, 1) {
		p, err := user.Buy(world.Market, key)
		if errors.Is(err, game.ErrNoKey) {
			return c.Send(format.BadKey(key), tele.ModeHTML)
		} else if err != nil {
			return c.Send(format.NoMoney, tele.ModeHTML)
		}
		c.Send(fmt.Sprintf("🛒 Вы приобрели %s за %s.",
			format.Item(p.Item), format.Money(p.Price)), tele.ModeHTML)
	}
	return nil
}

type Eat struct {
	Universe *game.Universe
}

var eatRe = regexp.MustCompile("^!(с[ъь]есть|еда) (.*)")

func (h *Eat) Match(s string) bool {
	return eatRe.MatchString(s)
}

func (h *Eat) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	ate := false
	defer func() {
		if ate {
			c.Send(format.EnergyRemaining(user.Energy), tele.ModeHTML)
		}
	}()
	for _, key := range teleutil.NumArg(c, eatRe, 2) {
		if user.Energy == game.EnergyCap {
			return c.Send(format.NotHungry)
		}
		item, ok := user.Inventory.ByKey(key)
		if !ok {
			return c.Send(format.BadKey(key), tele.ModeHTML)
		}
		if !user.Eat(item) {
			return c.Send("🤮")
		}
		ate = true
		c.Send(format.Eat(item), tele.ModeHTML)
	}
	return nil
}

type EatQuick struct {
	Universe *game.Universe
}

var eatQuickRe = regexp.MustCompile("^!еда")

func (h *EatQuick) Match(s string) bool {
	return eatQuickRe.MatchString(s)
}

func (h *EatQuick) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	if user.Energy == game.EnergyCap {
		return c.Send(format.NotHungry)
	}
	i, ok := user.EatQuick()
	if !ok {
		return c.Send(format.NoFood)
	}
	return c.Send(format.Eat(i)+"\n\n"+
		format.EnergyRemaining(user.Energy), tele.ModeHTML)
}

type Fish struct {
	Universe *game.Universe
}

var fishRe = regexp.MustCompile("^!(р[ыі]балка|ловля рыб)")

func (h *Fish) Match(s string) bool {
	return fishRe.MatchString(s)
}

func (h *Fish) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	if user.Inventory.Count() > game.InventoryCap {
		return c.Send(format.InventoryFull)
	}
	rod, ok := user.FishingRod()
	if !ok {
		return c.Send("🎣 Приобретите удочку в магазине, прежде чем рыбачить.")
	}
	if !user.SpendEnergy(20) {
		return c.Send(format.NoEnergy)
	}
	item := user.Fish(rod)
	if rod.Durability < 0 {
		c.Send("🎣 Ваша удочка сломалась.")
	}
	chance := rand.Float64() + (-0.02 + 0.04*user.Luck())
	if chance < 0.5 {
		outcomes := [...]string{
			"Вы не смогли выудить рыбу.",
			"Рыба сорвалась с крючка.",
			"Рыба сорвала леску.",
			"Рыба скрылась в водорослях.",
			"Рыба выскользнула из рук.",
			"Вы отпустили рыбу обратно в воду.",
		}
		return c.Send("🎣 " + outcomes[rand.Intn(len(outcomes))])
	}
	user.Inventory.Add(item)
	return c.Send(fmt.Sprintf("🎣 %s получает %s",
		teleutil.Mention(c, user), format.Item(item)), tele.ModeHTML)
}

type Status struct {
	Universe *game.Universe
}

var statusRe = regexp.MustCompile("^!статус (.*)")

func (h *Status) Match(s string) bool {
	return statusRe.MatchString(s)
}

func (h *Status) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	status := teleutil.Args(c, statusRe)[1]
	const maxlen = 120
	if utf8.RuneCountInString(status) > maxlen {
		return c.Send(fmt.Sprintf("💬 Максимальная длина статуса %d символов.", maxlen))
	}
	user.Status = status
	return c.Send("✅ Статус установлен.")
}

type Sell struct {
	Universe *game.Universe
}

var sellRe = regexp.MustCompile("^!прода(ть|жа) (.*)")

func (h *Sell) Match(s string) bool {
	return sellRe.MatchString(s)
}

func (h *Sell) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	items := teleutil.NumArg(c, sellRe, 2)
	for _, key := range items {
		item, ok := user.Inventory.ByKey(key)
		if !ok {
			return c.Send(format.BadKey(key), tele.ModeHTML)
		}
		profit, ok := user.Sell(item)
		if !ok {
			return c.Send(fmt.Sprintf("🏪 Вы не можете продать %s.",
				format.Item(item)), tele.ModeHTML)
		}
		c.Send(fmt.Sprintf("💵 Вы продали %s, заработав %s.",
			format.Item(item), format.Money(profit)), tele.ModeHTML)
	}
	return nil
}

type Stack struct {
	Universe *game.Universe
}

var stackRe = regexp.MustCompile("^!сложить")

func (h *Stack) Match(s string) bool {
	return stackRe.MatchString(s)
}

func (h *Stack) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	if user.Stack() {
		return c.Send("💵 Вы сложили деньги.")
	}
	return c.Send("✅")
}

type Fight struct {
	Universe *game.Universe
}

var fightRe = regexp.MustCompile("^!(драка|дуэль|поединок|атака|битва|схватка|сражение|бой|борьба)")

func (h *Fight) Match(s string) bool {
	return fightRe.MatchString(s)
}

func (h *Fight) Handle(c tele.Context) error {
	reply, ok := teleutil.Reply(c)
	if !ok {
		return c.Send(format.RepostMessage)
	}
	if c.Sender().ID == reply.ID {
		return c.Send(format.CannotAttackYourself)
	}
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	opnt := world.UserByID(reply.ID)
	if !user.SpendEnergy(25) {
		return c.Send(format.NoEnergy)
	}
	c.Send(fmt.Sprintf("⚔️ <b>%s</b> <code>[%.2f]</code> <b><i>vs.</i></b> <b>%s</b> <code>[%.2f]</code>",
		teleutil.Mention(c, user.TUID), user.Strength(),
		teleutil.Mention(c, opnt.TUID), opnt.Strength()),
		tele.ModeHTML)
	winner, loser, rating := user.Fight(opnt)
	winnerMent := teleutil.Mention(c, winner.TUID)
	if i, ok := loser.Inventory.Random(); ok && rand.Float64() < 1.0/8 {
		if i.Type != game.ItemTypeWallet && loser.Inventory.Move(world.Floor, i) {
			c.Send(fmt.Sprintf("🥊 %s выбивает %s из проигравшего.",
				winnerMent, format.Item(i)), tele.ModeHTML)
		}
	}
	if i, ok := user.Inventory.Random(); ok && rand.Float64() < 1.0/12 {
		if user.Inventory.Move(world.Floor, i) {
			c.Send(fmt.Sprintf("🌀 %s уронил %s во время драки.",
				teleutil.Mention(c, user.TUID), format.Item(i)), tele.ModeHTML)
		}
	}
	return c.Send(fmt.Sprintf("🏆 %s <code>(+%.1f)</code> выигрывает в поединке.",
		winnerMent, rating), tele.ModeHTML)
}

type Profile struct {
	Universe   *game.Universe
	AvatarPath string
}

var profileRe = regexp.MustCompile("^!(профиль|стат)")

func (h *Profile) Match(s string) bool {
	return profileRe.MatchString(s)
}

func (h *Profile) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	const profile = "📇 <b>%s %s</b>\n<code>%s  %s  %s  %s  %s</code>\n\n%s\n\n%s\n\n%s"
	mods := user.Modset().List()
	out := fmt.Sprintf(profile,
		format.ModifierTitles(mods),
		teleutil.Mention(c, c.Sender()),
		format.Energy(user.Energy),
		format.Rating(user.Rating),
		format.Strength(user.Strength()),
		format.Luck(user.Luck()),
		format.Messages(user.Messages),
		format.ModifierDescriptions(mods),
		format.ModifierEmojis(mods),
		format.Status(user.Status),
	)
	if a, ok := avatar(h.AvatarPath, c.Sender().ID); ok {
		a.Caption = out
		return c.Send(a, tele.ModeHTML)
	}
	return c.Send(out, tele.ModeHTML)
}

func avatar(dir string, id int64) (a *tele.Photo, ok bool) {
	_, err := os.Stat(dir)
	if err != nil {
		return nil, false
	}
	f := tele.FromDisk(filepath.Join(dir, strconv.FormatInt(id, 10)))
	if f.OnDisk() {
		return &tele.Photo{File: f}, true
	}
	return nil, false
}

type Dice struct {
	Universe *game.Universe
}

var diceRe = regexp.MustCompile("!кости (.*)")

func (h *Dice) Match(s string) bool {
	return diceRe.MatchString(s)
}

func (h *Dice) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	if _, ok := user.Dice(); !ok {
		return c.Send("🎲 У вас нет костей.")
	}
	args := teleutil.NumArg(c, diceRe, 1)
	if len(args) != 1 {
		return c.Send("💵 Сделайте ставку.")
	}
	bet := args[0]
	const minbet = 100
	if bet < minbet {
		return c.Send(fmt.Sprintf("💵 Минимальная ставка %s.",
			format.Money(minbet)), tele.ModeHTML)
	}
	if world.Casino.GameGoing() {
		return c.Send("🎲 Игра уже идет.")
	}
	if !user.SpendMoney(bet) {
		return c.Send("💵 Недостаточно средств.")
	}
	if err := world.Casino.PlayDice(
		user, bet,
		func() (int, error) {
			msg, err := tele.Cube.Send(c.Bot(), c.Chat(), nil)
			if err != nil {
				return 0, err
			}
			return msg.Dice.Value, nil
		},
		func() {
			c.Send(fmt.Sprintf("<i>Время вышло: вы потеряли %s</i>",
				format.Money(bet)), tele.ModeHTML)
		},
	); err != nil {
		return err
	}
	return c.Send(fmt.Sprintf("🎲 %s играет на %s\nУ вас <code>%d секунд</code> на то, чтобы кинуть кости!",
		teleutil.Mention(c, c.Sender()), format.Money(bet), world.Casino.Timeout/time.Second), tele.ModeHTML)
}

type Roll struct {
	Universe *game.Universe
}

func (h *Roll) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	game, ok := world.Casino.DiceGame()
	if !ok || game.Player != user {
		return nil
	}
	game.Finish()
	switch score := c.Message().Dice.Value; {
	case score > game.CasinoScore:
		win := game.Bet * 2
		game.Player.AddMoney(win)
		return c.Send(fmt.Sprintf("💥 Вы выиграли %s",
			format.Money(win)), tele.ModeHTML)
	case score == game.CasinoScore:
		draw := game.Bet
		game.Player.AddMoney(draw)
		return c.Send("🎲 Ничья.")
	}
	return c.Send("😵 Вы проиграли.")
}

type TopStrong struct {
	Universe *game.Universe
}

var topStrongRe = regexp.MustCompile("^!топ сил")

func (h *TopStrong) Match(s string) bool {
	return topStrongRe.MatchString(s)
}

func (h *TopStrong) Handle(c tele.Context) error {
	world, _ := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	users := world.SortedUsers(game.ByStrength)
	users = users[:min(len(users), 5)]
	list := []string{"🏋️‍♀️ <b>Самые сильные пользователи</b>"}
	for i, u := range users {
		list = append(list, fmt.Sprintf("<b><i>%d.</i></b> %s %s",
			i+1, teleutil.Mention(c, u.TUID), format.Strength(u.Strength())))
	}
	return c.Send(strings.Join(list, "\n"), tele.ModeHTML)
}

type TopRating struct {
	Universe *game.Universe
}

var topRating = regexp.MustCompile("^!(рейтинг|ммр|эло)")

func (h *TopRating) Match(s string) bool {
	return topRating.MatchString(s)
}

func (h *TopRating) Handle(c tele.Context) error {
	world, _ := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	users := world.SortedUsers(game.ByElo)
	users = users[:min(len(users), 5)]
	list := []string{"🏆 <b>Боевой рейтинг</b>"}
	for i, u := range users {
		list = append(list, fmt.Sprintf("<b><i>%d.</i></b> %s %s",
			i+1, teleutil.Mention(c, u.TUID), format.Rating(u.Rating)))
	}
	return c.Send(strings.Join(list, "\n"), tele.ModeHTML)
}

type TopRich struct {
	Universe *game.Universe
}

var topRich = regexp.MustCompile("^!топ бога[тч]")

func (h *TopRich) Match(s string) bool {
	return topRich.MatchString(s)
}

func (h *TopRich) Handle(c tele.Context) error {
	world, _ := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	users := world.SortedUsers(game.ByWealth)
	users = users[:min(len(users), 5)]
	list := []string{"💵 <b>Самые богатые пользователи</b>"}
	for i, u := range users {
		list = append(list, fmt.Sprintf("<b><i>%d.</i></b> %s %s",
			i+1, teleutil.Mention(c, u.TUID), format.Money(u.Total())))
	}
	return c.Send(strings.Join(list, "\n"), tele.ModeHTML)
}

type Capital struct {
	Universe *game.Universe
}

var capitalRe = regexp.MustCompile("^!капитал")

func (h *Capital) Match(s string) bool {
	return capitalRe.MatchString(s)
}

func (h *Capital) Handle(c tele.Context) error {
	world, _ := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	total, avg := world.Capital()
	users := world.SortedUsers(game.ByWealth)
	users = users[:min(len(users), 5)]
	rich := users[0]
	balance := rich.Total()
	list := []string{
		fmt.Sprintf("💸 Капитал беседы <b>%s</b>: %s\n",
			c.Chat().Title, format.Money(total)),
		fmt.Sprintf("<i>В среднем на счету: %s</i>\n",
			format.Money(avg)),
		fmt.Sprintf("<i>В руках магната %s %s,</i>",
			teleutil.Mention(c, users[0].TUID), format.Money(balance)),
		fmt.Sprintf("<i>или %s от общего количества средств.</i>\n",
			format.Percentage(float64(balance)/float64(total))),
	}
	return c.Send(strings.Join(list, "\n"), tele.ModeHTML)
}

type Balance struct {
	Universe *game.Universe
}

var balanceRe = regexp.MustCompile("^!(баланс|деньги)")

func (h *Balance) Match(s string) bool {
	return balanceRe.MatchString(s)
}

func (h *Balance) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()
	return c.Send(fmt.Sprintf("💵 Ваш баланс: %s", format.Money(user.Total())), tele.ModeHTML)
}

type Energy struct {
	Universe *game.Universe
}

var energyRe = regexp.MustCompile("^!энергия")

func (h *Energy) Match(s string) bool {
	return energyRe.MatchString(s)
}

func (h *Energy) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()
	return c.Send(fmt.Sprintf("%s Запас энергии: %s",
		tern(user.Energy < game.EnergyCap/2, "🪫", "🔋"),
		format.EnergyOutOf(user.Energy, game.EnergyCap)), tele.ModeHTML)
}

type NamePet struct {
	Universe *game.Universe
}

var namePetRe = regexp.MustCompile("^!назвать (.*)")

func (h *NamePet) Match(s string) bool {
	return namePetRe.MatchString(s)
}

func (h *NamePet) Handle(c tele.Context) error {
	world, user := teleutil.Lock(c, h.Universe)
	defer world.Unlock()

	name := teleutil.Args(c, namePetRe)[1]
	pet, ok := user.Pet()
	if !ok {
		return c.Send("🐈 У вас нет питомца.")
	}
	if pet.Name != "" {
		return c.Send("🐈 У вашего питомца уже есть имя.")
	}
	if !pet.SetName(name) {
		return c.Send("🐈 Такое имя не подходит для питомца.")
	}
	return c.Send(fmt.Sprintf("🐈 Вы назвали питомца <code>%s</code>.",
		name), tele.ModeHTML)
}
