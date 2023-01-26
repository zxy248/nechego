package danbooru

import (
	"errors"
	"log"
	"net/http"
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

// New starts workers filling buffer channels and returns Danbooru.
// maxSize is a maximum size of a picture after which it is filtered.
// timeout specifies the time interval after which Get functions will timeout.
func New(url string, maxSize int, timeout time.Duration) *Danbooru {
	d := &Danbooru{
		URL:     url,
		MaxSize: maxSize,
		Timeout: timeout,
		pics:    make(chan *Pic, 32),
		nsfw:    make(chan *Pic, 32),
	}
	errs := make(chan error, 1)
	test := func(f *file) bool { return f.URL != "" && f.Size < d.MaxSize }
	testNSFW := func(f *file) bool { return test(f) && f.Rating.NSFW() }
	d.pipeline(d.pics, errs, test)
	d.pipeline(d.nsfw, errs, testNSFW)
	go func() {
		for err := range errs {
			var code errStatusCode
			if errors.As(err, &code) && code == http.StatusTooManyRequests {
				// StatusTooManyRequests arises frequently when
				// multiple workers are filling the picture
				// channels during application start. Skip it.
				continue
			}
			log.Printf("danbooru: %s", err)
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
