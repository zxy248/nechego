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
	CommandPoppy
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
	CommandDice
	CommandGame
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
	eblanRe       = regexp.MustCompile("(?i)^![ие][б6п*]?лан[А-я]* дня")
	masyunyaRe    = regexp.MustCompile("(?i)^(!ма[нс]ю[нс][а-я]*[пая])")
	helloRe       = regexp.MustCompile(constructHelloRe("п[рл]ивет", "хай", "зд[ао]ров", "ку", "здрав"))
	weatherRe     = regexp.MustCompile("(?i)^!погода ([-А-я]+)")
	probabilityRe = regexp.MustCompile("(?i)^!инфа *(.*)")
	whoRe         = regexp.MustCompile("(?i)^!кто *(.*)")
	listRe        = regexp.MustCompile("(?i)^!список *(.*)")
	topRe         = regexp.MustCompile("(?i)^!топ[- ]*(\\d*) *(.*)")
)

// recognizeCommand returns the command contained in the input string.
func recognizeCommand(s string) Command {
	switch {
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
	case masyunyaRe.MatchString(s) || startsWith(s, "Масюня 🎀"):
		return CommandMasyunya
	case startsWith(s, "!паппи", "Паппи 🦊"):
		return CommandPoppy
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
	case startsWith(s, "!касп", "!кот касп"):
		return CommandCasper
	case startsWith(s, "!зевс"):
		return CommandZeus
	case startsWith(s, "!кот", "!кош"):
		return CommandCat
	case startsWith(s, "!пик"):
		return CommandPic
	case startsWith(s, "!кости"):
		return CommandDice
	case startsWith(s, "!игр"):
		return CommandGame
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
	s = strings.ToLower(s)
	for _, p := range prefix {
		p = strings.ToLower(p)
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

const (
	helloPrefix = "((^|[^а-я])"
	helloSuffix = "[а-я]*([^а-я]|$))"
)

func constructHelloRe(hello ...string) string {
	var l []string
	for _, h := range hello {
		l = append(l, helloPrefix+h+helloSuffix)
	}
	return "(?i)" + strings.Join(l, "|")
}
