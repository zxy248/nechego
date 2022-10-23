package input

import "golang.org/x/exp/slices"

type Command int

const (
	CommandUnknown Command = 0

	// Management
	CommandHelp          = 1
	CommandInfo          = 2
	CommandTurnOn        = 3
	CommandTurnOff       = 4
	CommandKeyboardOpen  = 5
	CommandKeyboardClose = 6

	// General
	CommandProbability = 100
	CommandWho         = 101
	CommandPair        = 102
	CommandEblan       = 103
	CommandWeather     = 104
	CommandList        = 105
	CommandTop         = 106
	CommandTitle       = 107

	// Neural networks
	CommandCat    = 200
	CommandAnime  = 201
	CommandFurry  = 202
	CommandFlag   = 203
	CommandPerson = 204
	CommandHorse  = 205
	CommandArt    = 206
	CommandCar    = 207

	// Fun
	CommandMasyunya = 300
	CommandPoppy    = 301
	CommandHello    = 302
	CommandMouse    = 303
	CommandTikTok   = 304
	CommandGame     = 305
	CommandSima     = 306

	// Pictures
	CommandBasili = 400
	CommandCasper = 401
	CommandZeus   = 402
	CommandPic    = 403

	// Fishing
	CommandFishing      = 500
	CommandFishingRod   = 501
	CommandEatFish      = 502
	CommandFishList     = 503
	CommandFreezeFish   = 504
	CommandFreezer      = 505
	CommandUnfreezeFish = 506
	CommandSellFish     = 507

	// Economy
	CommandDice     = 600
	CommandBalance  = 601
	CommandTransfer = 602
	CommandTopRich  = 603
	CommandTopPoor  = 604
	CommandCapital  = 605
	CommandDeposit  = 606
	CommandWithdraw = 607
	CommandBank     = 608
	CommandDebt     = 609
	CommandRepay    = 610

	// Fight
	CommandProfile   = 700
	CommandTopStrong = 701
	CommandTopWeak   = 702
	CommandFight     = 703
	CommandStrength  = 704
	CommandEnergy    = 705
	CommandRating    = 706

	// Admin
	CommandAdmin       = 800
	CommandBan         = 801
	CommandUnban       = 802
	CommandForbid      = 803
	CommandPermit      = 804
	CommandParliament  = 805
	CommandImpeachment = 806

	// Pets
	CommandPet     = 900
	CommandBuyPet  = 901
	CommandNamePet = 902
	CommandDropPet = 903
)

var commandText = map[Command]string{
	CommandUnknown: "...",

	CommandHelp:          "!помощь",
	CommandInfo:          "!информация",
	CommandTurnOn:        "!включить",
	CommandTurnOff:       "!выключить",
	CommandKeyboardOpen:  "!открыть",
	CommandKeyboardClose: "!закрыть",

	CommandProbability: "!инфа",
	CommandWho:         "!кто",
	CommandPair:        "!пара дня",
	CommandEblan:       "!еблан дня",
	CommandWeather:     "!погода",
	CommandList:        "!список",
	CommandTop:         "!топ",
	CommandTitle:       "!имя",

	CommandCat:    "!кот",
	CommandAnime:  "!аниме",
	CommandFurry:  "!фурри",
	CommandFlag:   "!флаг",
	CommandPerson: "!чел",
	CommandHorse:  "!лошадь",
	CommandArt:    "!арт",
	CommandCar:    "!авто",

	CommandMasyunya: "!масюня",
	CommandPoppy:    "!паппи",
	CommandHello:    "!привет",
	CommandMouse:    "!мыш",
	CommandTikTok:   "!тикток",
	CommandGame:     "!игра",

	CommandBasili: "!марсик",
	CommandCasper: "!каспер",
	CommandZeus:   "!зевс",
	CommandPic:    "!пик",

	CommandFishing:      "!рыбалка",
	CommandFishingRod:   "!удочка",
	CommandEatFish:      "!еда",
	CommandFishList:     "!рыба",
	CommandFreezeFish:   "!заморозка",
	CommandFreezer:      "!холодильник",
	CommandUnfreezeFish: "!разморозка",
	CommandSellFish:     "!продажа",

	CommandDice:     "!кости",
	CommandBalance:  "!баланс",
	CommandTransfer: "!перевод",
	CommandTopRich:  "!топ богатых",
	CommandTopPoor:  "!топ бедных",
	CommandCapital:  "!капитал",
	CommandDeposit:  "!депозит",
	CommandWithdraw: "!обнал",
	CommandBank:     "!банк",
	CommandDebt:     "!кредит",
	CommandRepay:    "!погасить",

	CommandProfile:   "!профиль",
	CommandTopStrong: "!топ сильных",
	CommandTopWeak:   "!топ слабых",
	CommandFight:     "!драка",
	CommandStrength:  "!сила",
	CommandEnergy:    "!энергия",
	CommandRating:    "!рейтинг",

	CommandAdmin:       "!админ дня",
	CommandBan:         "!бан",
	CommandUnban:       "!разбан",
	CommandForbid:      "!запретить",
	CommandPermit:      "!разрешить",
	CommandParliament:  "!парламент",
	CommandImpeachment: "!импичмент",

	CommandPet:     "!питомец",
	CommandBuyPet:  "!взять",
	CommandNamePet: "!назвать",
	CommandDropPet: "!выкинуть",
}

func (c Command) String() string {
	return commandText[c]
}

func AllCommands() []Command {
	var l []Command
	for c := range commandText {
		l = append(l, c)
	}
	return l
}

func IsUnknownCommand(c Command) bool {
	return c == 0
}

func IsManagementCommand(c Command) bool {
	return c >= 1 && c < 100
}

func IsGeneralCommand(c Command) bool {
	return c >= 100 && c < 200
}

func IsNeuralNetworkCommand(c Command) bool {
	return c >= 200 && c < 300
}

func IsFunCommand(c Command) bool {
	return c >= 300 && c < 400
}

func IsPictureCommand(c Command) bool {
	return c >= 400 && c < 500
}

func IsFishingCommand(c Command) bool {
	return c >= 500 && c < 600
}

func IsEconomyCommand(c Command) bool {
	return c >= 600 && c < 700
}

func IsFightCommand(c Command) bool {
	return c >= 700 && c < 800
}

func IsAdminCommand(c Command) bool {
	return c >= 800 && c < 900
}

func IsPetCommand(c Command) bool {
	return c >= 900 && c < 1000
}

func IsImmune(c Command) bool {
	return IsAdminCommand(c) || slices.Contains(immuneCommands, c)
}

var immuneCommands = []Command{
	CommandTurnOn,
}
