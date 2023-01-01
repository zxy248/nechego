package fishing

import "math/rand"

type Outcome int

const (
	Lost Outcome = iota
	Off
	Tear
	Seagrass
	Slip
	Release
	Collect
)

func (o Outcome) Success() bool {
	return o == goodOutcome()
}

func (o Outcome) String() string {
	return outcomeDescriptions[o]
}

func goodOutcome() Outcome {
	return Collect
}

var badOutcomes = []Outcome{
	Lost,
	Off,
	Tear,
	Seagrass,
	Slip,
	Release,
}

func badOutcome() Outcome {
	return badOutcomes[rand.Intn(len(badOutcomes))]
}

const SuccessChance = 0.5

type Session struct {
	Outcome
	Fish
}

func Cast() Session {
	return CastProbability(SuccessChance)
}

func CastProbability(success float64) Session {
	r := rand.Float64()
	var outcome Outcome
	if r <= success {
		outcome = goodOutcome()
	} else {
		outcome = badOutcome()
	}
	fish := RandomFish()
	return Session{outcome, fish}
}
