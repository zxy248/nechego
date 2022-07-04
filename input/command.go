package input

// Command is an integer in one of the intervals:
// Management		[001, 100)
// General		[100, 200)
// Neural network	[200, 300)
// Fun			[300, 400)
// Picture		[400, 500)
type Command int

// TODO: add CommandRating, CommandRatingIncrement, CommandRatingDecrement
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
