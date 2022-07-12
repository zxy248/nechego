package input

import "golang.org/x/exp/slices"

// Command is an integer in one of the intervals:
// Management		[001, 100)
// General		[100, 200)
// Neural network	[200, 300)
// Fun			[300, 400)
// Picture		[400, 500)
type Command int

const (
	CommandUnknown       Command = 0
	CommandHelp                  = 1
	CommandInfo                  = 2
	CommandBan                   = 3
	CommandUnban                 = 4
	CommandTurnOn                = 5
	CommandTurnOff               = 6
	CommandTitle                 = 9
	CommandForbid                = 10
	CommandPermit                = 11
	CommandProbability           = 100
	CommandWho                   = 101
	CommandPair                  = 102
	CommandEblan                 = 103
	CommandWeather               = 104
	CommandList                  = 105
	CommandTop                   = 106
	CommandAdmin                 = 107
	CommandFight                 = 108
	CommandBalance               = 109
	CommandTransfer              = 110
	CommandProfile               = 111
	CommandTopRich               = 112
	CommandTopPoor               = 113
	CommandCapital               = 114
	CommandStrength              = 115
	CommandEnergy                = 116
	CommandFishing               = 117
	CommandFishingRod            = 118
	CommandTopStrong             = 119
	CommandEatFish               = 120
	CommandDeposit               = 121
	CommandWithdraw              = 122
	CommandBank                  = 123
	CommandDebt                  = 124
	CommandRepay                 = 125
	CommandTopWeak               = 126
	CommandParliament            = 127
	CommandImpeachment           = 128
	CommandCat                   = 200
	CommandAnime                 = 201
	CommandFurry                 = 202
	CommandFlag                  = 203
	CommandPerson                = 204
	CommandHorse                 = 205
	CommandArt                   = 206
	CommandCar                   = 207
	CommandMasyunya              = 300
	CommandPoppy                 = 301
	CommandHello                 = 302
	CommandMouse                 = 303
	CommandTikTok                = 304
	CommandGame                  = 305
	CommandDice                  = 306
	CommandKeyboardOpen          = 307
	CommandKeyboardClose         = 308
	CommandBasili                = 400
	CommandCasper                = 401
	CommandZeus                  = 402
	CommandPic                   = 403
)

var commandText = map[Command]string{
	CommandUnknown:       "...",
	CommandHelp:          "!помощь",
	CommandInfo:          "!информация",
	CommandBan:           "!бан",
	CommandUnban:         "!разбан",
	CommandTurnOn:        "!включить",
	CommandTurnOff:       "!выключить",
	CommandTitle:         "!имя",
	CommandForbid:        "!запретить",
	CommandPermit:        "!разрешить",
	CommandProbability:   "!инфа",
	CommandWho:           "!кто",
	CommandPair:          "!пара дня",
	CommandEblan:         "!еблан дня",
	CommandWeather:       "!погода",
	CommandList:          "!список",
	CommandTop:           "!топ",
	CommandAdmin:         "!админ",
	CommandFight:         "!драка",
	CommandBalance:       "!баланс",
	CommandTransfer:      "!перевод",
	CommandProfile:       "!профиль",
	CommandTopRich:       "!топ богатых",
	CommandTopPoor:       "!топ бедных",
	CommandCapital:       "!капитал",
	CommandStrength:      "!сила",
	CommandEnergy:        "!энергия",
	CommandFishing:       "!рыбалка",
	CommandFishingRod:    "!удочка",
	CommandTopStrong:     "!топ сильных",
	CommandEatFish:       "!еда",
	CommandDeposit:       "!депозит",
	CommandWithdraw:      "!обнал",
	CommandBank:          "!банк",
	CommandDebt:          "!кредит",
	CommandRepay:         "!погасить",
	CommandTopWeak:       "!топ слабых",
	CommandParliament:    "!парламент",
	CommandImpeachment:   "!импичмент",
	CommandCat:           "!кот",
	CommandAnime:         "!аниме",
	CommandFurry:         "!фурри",
	CommandFlag:          "!флаг",
	CommandPerson:        "!чел",
	CommandHorse:         "!лошадь",
	CommandArt:           "!арт",
	CommandCar:           "!авто",
	CommandMasyunya:      "!масюня",
	CommandPoppy:         "!паппи",
	CommandHello:         "!привет",
	CommandMouse:         "!мыш",
	CommandTikTok:        "!тикток",
	CommandGame:          "!игра",
	CommandDice:          "!кости",
	CommandKeyboardOpen:  "!открыть",
	CommandKeyboardClose: "!закрыть",
	CommandBasili:        "!марсик",
	CommandCasper:        "!каспер",
	CommandZeus:          "!зевс",
	CommandPic:           "!пик",
}

func CommandText(c Command) string {
	return commandText[c]
}

func AllCommands() []Command {
	var l []Command
	for c := range commandText {
		l = append(l, c)
	}
	return l
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

var immuneCommands = []Command{
	CommandUnban,
	CommandTurnOn,
	CommandPermit,
	CommandParliament,
	CommandImpeachment,
	CommandAdmin,
}

func IsImmune(c Command) bool {
	return slices.Contains(immuneCommands, c)
}
