package danbooru

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const URL = "https://danbooru.donmai.us/posts/random.json"

// MaxSize is a maximum size of a photo allowed by Telegram.
const MaxSize = 5 << 20 // 5 MB

// Picture type.
// Used to access a corresponding filter function.
const (
	All int = iota
	NSFW

	typN
)

// keep returns true if a file has a URL and its size is small enough
// for a Telegram photo.
var keep = func(f *file) bool { return f.URL != "" && f.Size < MaxSize }

// filters corresponding to each picture type.
var filters = [...]func(*file) bool{
	All:  keep,
	NSFW: func(f *file) bool { return keep(f) && f.Rating.NSFW() },
}

// Pic contains the picture's bytes and its rating.
type Pic struct {
	Data []byte
	Rating
}

// Danbooru structure represents the Danbooru API and holds the
// bufferized picture channels of each type.
type Danbooru struct {
	URL      string
	Timeout  time.Duration
	channels [typN]chan *Pic
}

// New starts workers filling buffer channels and returns Danbooru.
// Get function will return error after the specified timeout.
func New(url string, timeout time.Duration, cache int) *Danbooru {
	d := &Danbooru{
		URL:     url,
		Timeout: timeout,
	}
	errs := make(chan error, 1)
	for i := 0; i < typN; i++ {
		ch := make(chan *Pic, cache)
		d.pipeline(ch, errs, filters[i])
		d.channels[i] = ch
	}
	go func() {
		for err := range errs {
			var code errStatusCode
			if errors.As(err, &code) && code == http.StatusTooManyRequests {
				// StatusTooManyRequests arises frequently when
				// multiple workers are filling the picture
				// channels during application start.
				// This type of error can be safely ignored.
				continue
			}
			log.Printf("danbooru: %s", err)
		}
	}()
	return d
}

// Get returns a picture of the specified type from one of the channels.
func (d *Danbooru) Get(typ int) (*Pic, error) {
	if typ < 0 || typ >= typN {
		return nil, fmt.Errorf("unexpected type %d", typ)
	}
	select {
	case p := <-d.channels[typ]:
		return p, nil
	case <-time.After(d.Timeout):
		return nil, errors.New("danbooru: Get timeout")
	}
}
