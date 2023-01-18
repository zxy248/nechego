package danbooru

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

var errZeroRetries = errors.New("danbooru: zero retries left")

type Pic struct {
	Data []byte
	Rating
}

type Danbooru struct {
	URL     string
	MaxSize int
	Timeout time.Duration
	pics    <-chan Pic
	nsfw    <-chan Pic
}

func New(url string, maxSize int, timeout time.Duration) *Danbooru {
	d := &Danbooru{URL: url, MaxSize: maxSize, Timeout: timeout}
	const size, workers = 16, 4
	pics, errPics := channel(d.get, 4, size, workers)
	nsfw, errNSFW := channel(d.getNSFW, 32, size, workers)
	d.pics = pics
	d.nsfw = nsfw
	go func() {
		for {
			select {
			case err := <-errPics:
				log.Printf("danbooru (pics): %v", err)
			case err := <-errNSFW:
				log.Printf("danbooru (nsfw): %v", err)
			}
		}
	}()
	return d
}

func (d *Danbooru) Get() (Pic, error) {
	select {
	case p := <-d.pics:
		return p, nil
	case <-time.After(d.Timeout):
		return Pic{}, errors.New("timeout")
	}
}

func (d *Danbooru) GetNSFW() (Pic, error) {
	select {
	case p := <-d.nsfw:
		return p, nil
	case <-time.After(d.Timeout):
		return Pic{}, errors.New("timeout")
	}
}

func (d *Danbooru) get(retries int) (Pic, error) {
	for i := 0; i < retries; i++ {
		r, err := d.random(retries)
		if err != nil {
			return Pic{}, err
		}
		data, err := download(r.FileURL)
		if err != nil {
			return Pic{}, err
		}
		return Pic{data, rate(r.Rating)}, nil
	}
	return Pic{}, errZeroRetries
}

func (d *Danbooru) getNSFW(retries int) (Pic, error) {
	for i := 0; i < retries; i++ {
		r, err := d.random(retries)
		if err != nil {
			return Pic{}, err
		}
		rating := rate(r.Rating)
		if !rating.NSFW() {
			continue
		}
		data, err := download(r.FileURL)
		if err != nil {
			return Pic{}, err
		}
		return Pic{data, rating}, nil
	}
	return Pic{}, errZeroRetries
}

type response struct {
	FileURL  string `json:"file_url"`
	FileSize int    `json:"file_size"`
	Rating   string `json:"rating"`
}

func (d *Danbooru) random(retries int) (*response, error) {
	if retries <= 0 {
		return nil, errZeroRetries
	}
	r, err := http.Get(d.URL)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var x response
	if err := json.NewDecoder(r.Body).Decode(&x); err != nil {
		return nil, err
	}
	if x.FileURL == "" || x.FileSize > d.MaxSize {
		return d.random(retries - 1)
	}
	return &x, nil
}

func download(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type getter func(retries int) (Pic, error)

func channel(f getter, retries, size, workers int) (<-chan Pic, <-chan error) {
	pics := make(chan Pic, size)
	errors := make(chan error, 1)
	for i := 0; i < workers; i++ {
		go func() {
			for {
				pic, err := f(retries)
				if err != nil {
					errors <- err
					continue
				}
				pics <- pic
			}
		}()
	}
	return pics, errors
}
