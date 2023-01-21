package danbooru

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type file struct {
	URL    string `json:"file_url"`
	Size   int    `json:"file_size"`
	Rating Rating `json:"rating"`
}

func get(url string) (*file, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code %s", http.StatusText(r.StatusCode))
	}

	var f file
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		return nil, err
	}
	return &f, nil
}

func getter(url string, files chan<- *file, errs chan<- error) {
	for {
		f, err := get(url)
		if err != nil {
			errs <- err
			time.Sleep(time.Second)
			continue
		}
		files <- f
	}
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

func downloader(in <-chan *file, out chan<- *Pic, errs chan<- error) {
	for f := range in {
		data, err := download(f.URL)
		if err != nil {
			errs <- err
			continue
		}
		out <- &Pic{data, f.Rating}
	}
}

func filter(in <-chan *file, out chan<- *file, keep func(*file) bool) {
	for f := range in {
		if keep(f) {
			out <- f
		}
	}
}

func (d *Danbooru) pipeline(pics chan<- *Pic, errs chan<- error, keep func(*file) bool) {
	const getters, downloaders = 2, 4

	files := make(chan *file)
	for i := 0; i < getters; i++ {
		go getter(d.URL, files, errs)
	}

	filtered := make(chan *file)
	go filter(files, filtered, keep)

	for i := 0; i < downloaders; i++ {
		go downloader(filtered, pics, errs)
	}
}
