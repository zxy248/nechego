package danbooru

import "time"

var Default = New("https://danbooru.donmai.us/posts/random.json", 5<<20, 3*time.Second)

func Get() (*Pic, error) {
	return Default.Get()
}

func GetNSFW() (*Pic, error) {
	return Default.GetNSFW()
}
