package pictures

import (
	"fmt"
	"math/rand/v2"
)

type Anime struct{}

func (*Anime) URL() string {
	const format = "https://thisanimedoesnotexist.ai/results/psi-%s/seed%05d.png"
	psi := [...]string{"0.3", "0.4", "0.5", "0.6", "0.7", "0.8", "0.9", "1.0",
		"1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "2.0"}
	r1 := psi[rand.N(len(psi))]
	r2 := rand.N(100000)
	return fmt.Sprintf(format, r1, r2)
}

type Furry struct{}

func (*Furry) URL() string {
	const format = "https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%05d.jpg"
	return fmt.Sprintf(format, rand.N(100000))
}

type Flag struct{}

func (*Flag) URL() string {
	const format = "https://thisflagdoesnotexist.com/images/%d.png"
	return fmt.Sprintf(format, rand.N(5000))
}

type Soy struct{}

func (*Soy) URL() string {
	return "https://booru.soy/random_image/download"
}

type Cat struct{}

func (*Cat) URL() string {
	const format = "https://d2ph5fj80uercy.cloudfront.net/%02d/cat%d.jpg"
	return fmt.Sprintf(format, 1+rand.N(6), rand.N(5000))
}
