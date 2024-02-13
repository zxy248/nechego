package markov

import (
	"math/rand/v2"
	"strings"
)

const start = "__start__"
const end = "__end__"

type Chain struct {
	chain map[string][]string
}

func New(samples []string) *Chain {
	var tokens []string
	for _, s := range samples {
		words := strings.Split(s, " ")
		tokens = append(tokens, start)
		tokens = append(tokens, words...)
		tokens = append(tokens, end)
	}

	chain := map[string][]string{}
	for i, t := range tokens {
		if t == end {
			continue
		}
		k := hash(t)
		chain[k] = append(chain[k], tokens[i+1])
	}
	return &Chain{chain}
}

func (g *Chain) Generate() []string {
	words := []string{start}
	for {
		last := words[len(words)-1]
		choices := g.chain[hash(last)]
		next := choices[rand.N(len(choices))]
		if next == end {
			break
		}
		words = append(words, next)
	}
	return words[1:]
}

func hash(s string) string {
	return strings.ToLower(s)
}
