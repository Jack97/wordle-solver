package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	readWriter := bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))

	game := &Game{
		Dictionary: dictionary,
		Writer:     readWriter.Writer,
		FeedbackResolver: &InteractiveFeedbackResolver{
			ReadWriter: readWriter,
		},
	}

	result := game.Play()

	if result.Win {
		readWriter.WriteString(fmt.Sprintf("Completed the wordle, %d/6 guesses used.\n", len(result.Guesses)))
	} else {
		readWriter.WriteString(fmt.Sprintf("Failed to complete the wordle, %d possible answers remaining.\n", len(game.Dictionary.RemainingPossibleAnswers)))
	}

	readWriter.Flush()
}
