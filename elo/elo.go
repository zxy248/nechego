package elo

import "math"

const (
	ScoreWin  = 1.0
	ScoreDraw = 0.5
	ScoreLose = 0.0
	KDefault  = 20.0
)

func eloExpectedScore(a, b float64) float64 {
	return 1.0 / (1.0 + math.Pow(10.0, (b-a)/400.0))
}

func EloDelta(eloA, eloB, k, score float64) float64 {
	return k * (score - eloExpectedScore(eloA, eloB))
}
