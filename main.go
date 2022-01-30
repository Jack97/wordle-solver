package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func main() {
	dictionaryFile, err := ioutil.ReadFile("dictionary.json")
	if err != nil {
		log.Fatal(err)
	}

	dictionary := &Dictionary{}
	err = json.Unmarshal(dictionaryFile, &dictionary)
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		Dictionary: dictionary,
		FeedbackResolver: &TestFeedbackResolver{
			Answer: Word{'k', 'n', 'o', 'l', 'l'},
		},
	}

	result := game.Play()

	log.Printf("Win: %t", result.Win)

	for i, guess := range result.Guesses {
		log.Printf("Guess #%d: %s", i+1, guess)
	}
}
