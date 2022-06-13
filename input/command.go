package input

import (
	"regexp"
	"strings"
)

type Command int

const (
	CommandUnknown Command = iota
	CommandProbability
	CommandWho
	CommandCat
	CommandTitle
	CommandAnime
	CommandFurry
	CommandFlag
	CommandPerson
	CommandHorse
	CommandArt
	CommandCar
	CommandPair
	CommandEblan
	CommandMasyunya
	CommandHello
	CommandMouse
	CommandWeather
	CommandTikTok
	CommandList
	CommandTop
	CommandBasili
	CommandCasper
	CommandZeus
	CommandPic
	CommandKeyboardOpen
	CommandKeyboardClose
	CommandTurnOn
	CommandTurnOff
	CommandBan
	CommandUnban
	CommandInfo
	CommandHelp
)

var (
	eblanRe    = regexp.MustCompile("(?i)^![ие][б6п*]?лан[А-я]* дня")
	masyunyaRe = regexp.MustCompile("(?i)^(!ма[нс]ю[нс][а-я]*[пая])|(🎀 Масюня 🎀)")
	helloRe    = regexp.MustCompile("(?i)((^|[^а-я])п[рл]ивет[а-я]*([^а-я]|$))" +
		"|((^|[^а-я])хай[а-я]*([^а-я]|$))" +
		"|((^|[^а-я])зд[ао]ров[а-я]*([^а-я]|$))" +
		"|((^|[^а-я])ку[а-я]*([^а-я]|$))")
	weatherRe     = regexp.MustCompile("(?i)^!погода ([-А-я]+)$")
	probabilityRe = regexp.MustCompile("(?i)^!инфа(.*)")
	whoRe         = regexp.MustCompile("(?i)^!кто(.*)")
	listRe        = regexp.MustCompile("(?i)^!список *(.*)")
	topRe         = regexp.MustCompile("(?i)^!топ *(-?\\d*) *(.*)")
)

// recognizeCommand returns the command contained in the input string.
func recognizeCommand(s string) Command {
	switch s = strings.ToLower(s); {
	case probabilityRe.MatchString(s):
		return CommandProbability
	case whoRe.MatchString(s):
		return CommandWho
	case startsWith(s, "!имя"):
		return CommandTitle
	case startsWith(s, "!аним", "!мульт"):
		return CommandAnime
	case startsWith(s, "!фур"):
		return CommandFurry
	case startsWith(s, "!флаг"):
		return CommandFlag
	case startsWith(s, "!чел"):
		return CommandPerson
	case startsWith(s, "!лошадь", "!конь"):
		return CommandHorse
	case startsWith(s, "!арт"):
		return CommandArt
	case startsWith(s, "!авто", "!тачк", "!машин"):
		return CommandCar
	case startsWith(s, "!пара дня"):
		return CommandPair
	case eblanRe.MatchString(s):
		return CommandEblan
	case masyunyaRe.MatchString(s):
		return CommandMasyunya
	case helloRe.MatchString(s):
		return CommandHello
	case startsWith(s, "!мыш"):
		return CommandMouse
	case weatherRe.MatchString(s):
		return CommandWeather
	case startsWith(s, "!тикток"):
		return CommandTikTok
	case listRe.MatchString(s):
		return CommandList
	case topRe.MatchString(s):
		return CommandTop
	case startsWith(s, "!кот василия", "!кошка василия", "!марс", "!муся"):
		return CommandBasili
	case startsWith(s, "!каспер"):
		return CommandCasper
	case startsWith(s, "!зевс"):
		return CommandZeus
	case startsWith(s, "!кот", "!кош"):
		return CommandCat
	case startsWith(s, "!пик"):
		return CommandPic
	case startsWith(s, "!клав", "!открыт"):
		return CommandKeyboardOpen
	case startsWith(s, "!закрыт", "!скрыт"):
		return CommandKeyboardClose
	case startsWith(s, "!вкл"):
		return CommandTurnOn
	case startsWith(s, "!выкл"):
		return CommandTurnOff
	case startsWith(s, "!бан"):
		return CommandBan
	case startsWith(s, "!разбан"):
		return CommandUnban
	case startsWith(s, "!инфо"):
		return CommandInfo
	case startsWith(s, "!помощь", "!команды"):
		return CommandHelp
	}
	return CommandUnknown
}

// startsWith returns true if the input string starts with one of the specified prefixes; false otherwise.
func startsWith(s string, prefix ...string) bool {
	for _, p := range prefix {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
