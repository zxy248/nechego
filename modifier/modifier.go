package modifier

// Moder is implemented by any value that should have a corresponding modifier.
type Moder interface {
	Mod() (m *Mod, ok bool)
}

type Mod struct {
	Emoji       string
	Multiplier  float64
	Description string
}

func (x *Mod) Mod() (m *Mod, ok bool) {
	return x, true
}

var (
	Heavy        = &Mod{"🪨", -0.35, "Ваш инвентарь переполнен."}
	RatingFirst  = &Mod{"🥇", +0.03, "Вы на 1-м месте в рейтинге."}
	RatingSecond = &Mod{"🥈", +0.02, "Вы на 2-м месте в рейтинге."}
	RatingThird  = &Mod{"🥉", +0.01, "Вы на 3-м месте в рейтинге."}
	SMS          = &Mod{"📩", 0.0, "У вас есть непрочитанные сообщения."}
)

type Set map[*Mod]bool

func (s Set) Active(m *Mod) bool {
	return s[m]
}

func (s Set) Add(m *Mod) {
	s[m] = true
}

func (s Set) List() []*Mod {
	r := []*Mod{}
	for m := range s {
		r = append(r, m)
	}
	return r
}

func (s Set) Sum() float64 {
	r := 0.0
	for m := range s {
		r += m.Multiplier
	}
	return r
}
