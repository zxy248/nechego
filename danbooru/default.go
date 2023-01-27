package danbooru

import "time"

var Default = New("https://danbooru.donmai.us/posts/random.json", 3*time.Second)

func Get(typ int) (*Pic, error) {
	return Default.Get(typ)
}
