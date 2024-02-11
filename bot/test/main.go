package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mb-14/gomarkov"
)

func main() {
	//Create a chain of order 2
	chain := gomarkov.NewChain(1)

	//Feed in training data
	chain.Add(strings.Split("I want a cheese burger", " "))
	chain.Add(strings.Split("I want a chilled sprite", " "))
	chain.Add(strings.Split("I want to go to the movies", " "))

	//Get transition probability of a sequence
	prob, _ := chain.Generate([]string{"I"})
	fmt.Println(prob)
	//Output: 0.6666666666666666

	//You can even generate new text based on an initial seed
	chain.Add(strings.Split("Mother should I build the wall?", " "))
	chain.Add(strings.Split("Mother should I run for President?", " "))
	chain.Add(strings.Split("Mother should I trust the government?", " "))
	next, _ := chain.Generate([]string{"I"})
	fmt.Println(next)
	next, _ = chain.Generate([]string{"wall?"})
	fmt.Println(next)

	//The chain is JSON serializable
	jsonObj, _ := json.Marshal(chain)
	err := ioutil.WriteFile("model.json", jsonObj, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
