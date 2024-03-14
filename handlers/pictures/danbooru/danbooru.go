package danbooru

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const URL = "https://danbooru.donmai.us/posts/random.json"

type Picture struct {
	URL    string `json:"file_url"`
	Size   int    `json:"file_size"`
	Rating Rating `json:"rating"`
	Score  int    `json:"score"`
}

func Get() (*Picture, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, ErrBadStatus(resp.StatusCode)
	}

	var f Picture
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, err
	}
	return &f, nil
}

// ErrBadStatus is an error containing the status code of a response.
type ErrBadStatus int

func (e ErrBadStatus) Error() string {
	code := int(e)
	text := http.StatusText(code)
	return fmt.Sprintf("bad status code: %s (%d)", text, code)
}
