package fun

import (
	"io"
	"nechego/handlers"
	"net/http"
	"net/url"

	tele "gopkg.in/telebot.v3"
)

type Weather struct{}

var weatherRe = handlers.NewRegexp("^!погода (.+)")

func (h *Weather) Match(c tele.Context) bool {
	_, ok := weatherCity(c.Text())
	return ok
}

func (h *Weather) Handle(c tele.Context) error {
	const addr = "https://wttr.in/"
	const query = `?format=%l:+%c+%t+\n` +
		`Ощущается+как+%f\n\n` +
		`Ветер+—+%w\n` +
		`Влажность+—+%h\n` +
		`Давление+—+%P\n` +
		`Фаза+луны+—+%m\n` +
		`УФ-индекс+—+%u\n`
	city, _ := weatherCity(c.Text())
	city = url.PathEscape(city)

	resp, err := http.Get(addr + city + query)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return c.Send("☔️ Такого места не существует.")
	}
	if resp.StatusCode != http.StatusOK {
		return c.Send("☔️ Неудачный запрос.")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return c.Send(string(data))
}

func weatherCity(s string) (city string, ok bool) {
	m := weatherRe.FindStringSubmatch(s)
	if m == nil {
		return "", false
	}
	return m[1], true
}
