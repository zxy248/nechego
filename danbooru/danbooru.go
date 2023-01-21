package danbooru

import (
	"errors"
	"log"
	"time"
)

type Pic struct {
	Data []byte
	Rating
}

type Danbooru struct {
	URL     string
	MaxSize int
	Timeout time.Duration
	pics    chan *Pic
	nsfw    chan *Pic
}

func New(url string, maxSize int, timeout time.Duration) *Danbooru {
	d := &Danbooru{
		URL:     url,
		MaxSize: maxSize,
		Timeout: timeout,
		pics:    make(chan *Pic, 32),
		nsfw:    make(chan *Pic, 32),
	}
	test := func(f *file) bool { return f.URL != "" && f.Size < d.MaxSize }
	errs := make(chan error, 1)
	d.pipeline(d.pics, errs, test)
	d.pipeline(d.nsfw, errs, func(f *file) bool { return test(f) && f.Rating.NSFW() })
	go func() {
		for err := range errs {
			log.Printf("danbooru: %v", err)
		}
	}()
	return d
}

func (d *Danbooru) Get() (*Pic, error) {
	select {
	case p := <-d.pics:
		return p, nil
	case <-time.After(d.Timeout):
		return nil, errors.New("danbooru: Get timeout")
	}
}

func (d *Danbooru) GetNSFW() (*Pic, error) {
	select {
	case p := <-d.nsfw:
		return p, nil
	case <-time.After(d.Timeout):
		return nil, errors.New("danbooru: GetNSFW timeout")
	}
}
