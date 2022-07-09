package app

import (
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"

	tele "gopkg.in/telebot.v3"
)

const catURL = "https://thiscatdoesnotexist.com/"

// handleCat sends a picture of a cat.
func (a *App) handleCat(c tele.Context) error {
	pic, err := fetchPicture(catURL)
	if err != nil {
		return internalError(c, err)
	}
	return c.Send(pic)
}

const animeFormat = "https://thisanimedoesnotexist.ai/results/psi-%s/seed%s.png"

var animePsis = []string{"0.3", "0.4", "0.5", "0.6", "0.7", "0.8", "0.9", "1.0",
	"1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "2.0"}

// handleAnime sends an anime picture.
func (a *App) handleAnime(c tele.Context) error {
	psi := animePsis[rand.Intn(len(animePsis))]
	seed := fmt.Sprintf("%05d", rand.Intn(100000))
	url := fmt.Sprintf(animeFormat, psi, seed)
	return a.fetchAndSend(c, url)
}

const furFormat = "https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%s.jpg"

// handleFurry sends a furry picture.
func (a *App) handleFurry(c tele.Context) error {
	seed := fmt.Sprintf("%05d", rand.Intn(100000))
	url := fmt.Sprintf(furFormat, seed)
	return a.fetchAndSend(c, url)
}

const flagFormat = "https://thisflagdoesnotexist.com/images/%d.png"

// handleFlag sends a picture of a flag.
func (a *App) handleFlag(c tele.Context) error {
	seed := rand.Intn(5000)
	url := fmt.Sprintf(flagFormat, seed)
	return a.fetchAndSend(c, url)
}

const personURL = "https://thispersondoesnotexist.com/image"

// handlePerson sends a picture of a person.
func (a *App) handlePerson(c tele.Context) error {
	return a.fetchAndSend(c, personURL)
}

const horseURL = "https://thishorsedoesnotexist.com/"

// handleHorse sends a picture of a horse.
func (a *App) handleHorse(c tele.Context) error {
	return a.fetchAndSend(c, horseURL)
}

const artURL = "https://thisartworkdoesnotexist.com/"

// handleArt sends a picture of an art.
func (a *App) handleArt(c tele.Context) error {
	return a.fetchAndSend(c, artURL)
}

const carURL = "https://www.thisautomobiledoesnotexist.com/"

var carImageRe = regexp.MustCompile(
	"<img id = \"vehicle\" src=\"data:image/png;base64,(.+)\" class=\"center\">")

// handleCar sends a picture of a car.
func (a *App) handleCar(c tele.Context) error {
	r, err := http.Get(carURL)
	if err != nil {
		return internalError(c, err)
	}
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return internalError(c, err)
	}
	car, err := decodeCarImage(data)
	if err != nil {
		return internalError(c, err)
	}
	return c.Send(car)
}

func decodeCarImage(data []byte) (*tele.Photo, error) {
	b64 := carImageRe.FindStringSubmatch(string(data))[1]
	img, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	return photoFromBytes(img), nil
}

// fetchPicture returns a picture located at the specified URL.
func fetchPicture(url string) (*tele.Photo, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return photoFromBytes(body), nil
}

// fetchAndSend downloads a picture at url and sends it.
func (a *App) fetchAndSend(c tele.Context, url string) error {
	pic, err := fetchPicture(url)
	if err != nil {
		return internalError(c, err)
	}
	return c.Send(pic)
}
