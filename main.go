package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	errorLogger := newConsoleLogger(os.Stderr)

	dictionaryFile, err := ioutil.ReadFile("dictionary.json")
	if err != nil {
		errorLogger.Fatalf("Error reading dictionary.json: %s", err)
	}

	dictionary := &Dictionary{}
	err = json.Unmarshal(dictionaryFile, &dictionary)
	if err != nil {
		errorLogger.Fatalf("Error unmarshalling dictionary.json: %s", err)
	}

	logger := newConsoleLogger(os.Stdout)

	game := &Game{
		Dictionary: dictionary,
		FeedbackResolver: &InteractiveFeedbackResolver{
			Logger:  logger,
			Scanner: bufio.NewScanner(os.Stdin),
		},
	}

	var result *GameResult
	result, err = game.Play()
	if err != nil {
		errorLogger.Fatalf("Unexpected error: %s", err)
	}

	if result.Win {
		logger.Printf("Completed the wordle in %d/6 guesses", len(result.Guesses))
		return
	}

	errorLogger.Printf("Failed to complete the wordle, %d possible answers remaining", len(game.Dictionary.RemainingPossibleAnswers))
}

func newConsoleLogger(file *os.File) *log.Logger {
	return log.New(file, "", 0)
}
