package app

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	weatherTimeout      = 10 * time.Second
	weatherTimeoutError = UserError("Время запроса вышло ☔️")
	placeNotExists      = "Такого места не существует ☔️"
	weatherBadRequest   = "Неудачный запрос ☔️"
	weatherURL          = "https://wttr.in/"
	weatherFormat       = `?format=%l:+%c+%t+\nОщущается+как+%f\n\nВетер+—+%w\nВлажность+—+%h\nДавление+—+%P\nФаза+луны+—+%m\nУФ-индекс+—+%u\n`
)

// handleWeather sends the current weather for a given city
func (a *App) handleWeather(c tele.Context) error {
	place := getMessage(c).Argument()
	r, err := fetchWeather(place)
	if err != nil {
		if err.(*url.Error).Timeout() {
			return respondUserError(c, weatherTimeoutError)
		}
		return respondInternalError(c, err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		if r.StatusCode == http.StatusNotFound {
			return respondUserError(c, placeNotExists)
		}
		return respondUserError(c, weatherBadRequest)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return respondInternalError(c, err)
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
