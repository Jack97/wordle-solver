package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
)

type Words struct {
	AllowedGuesses  []string `json:"allowedGuesses"`
	PossibleAnswers []string `json:"possibleAnswers"`
}

type Result struct {
	Win     bool
	Guesses []string
}

var words = Words{}

func main() {
	wordsFile, err := ioutil.ReadFile("words.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(wordsFile, &words)
	if err != nil {
		log.Fatal(err)
	}

	result := playGame(words.PossibleAnswers[rand.Intn(len(words.PossibleAnswers))])

	log.Printf("Win: %t", result.Win)
	log.Printf("Guesses: %v", result.Guesses)
}

func playGame(answer string) Result {
	var (
		win             bool
		guesses         []string
		possibleAnswers = words.PossibleAnswers
	)

	for i := 0; i < 6; i++ {
		guess := makeGuess(possibleAnswers)

		guesses = append(guesses, guess)

		if guess == answer {
			win = true
			break
		}

		pattern := buildPattern(guess, answer)
		possibleAnswers = updatePossibleAnswers(guess, pattern, possibleAnswers)
	}

	return Result{
		Win:     win,
		Guesses: guesses,
	}
}

func makeGuess(possibleAnswers []string) string {
	if len(possibleAnswers) <= 2 {
		// Return a random guess when we're left with a 50/50 choice
		return possibleAnswers[rand.Intn(len(possibleAnswers))]
	}

	return calculateOptimalGuess(possibleAnswers)
}

const (
	GREY = iota
	YELLOW
	GREEN
)

func calculateOptimalGuess(possibleAnswers []string) string {
	var (
		maxInformationGain float64
		optimalGuess       string
	)

	for _, guess := range append(words.AllowedGuesses, possibleAnswers...) {
		patterns := map[[5]byte]float64{}

		for _, possibleAnswer := range possibleAnswers {
			pattern := buildPattern(guess, possibleAnswer)
			patterns[pattern] += 1.0 / float64(len(possibleAnswers))
		}

		// Calculate the entropy
		var sum float64
		for _, value := range patterns {
			sum += value * math.Log(value)
		}
		informationGain := -sum

		// The word with the maximal entropy is the optimal guess, as it tells us the most information
		// and allows us to efficiently reduce the list of remaining possible answers
		if informationGain > maxInformationGain {
			maxInformationGain = informationGain
			optimalGuess = guess
		}
	}

	return optimalGuess
}

func buildPattern(guess, answer string) [5]byte {
	possibleAnswerChars := map[byte]uint8{}
	matchingChars := map[byte]uint8{}

	for i := 0; i < 5; i++ {
		possibleAnswerChars[answer[i]]++

		if guess[i] == answer[i] {
			matchingChars[guess[i]]++
		}
	}

	pattern := [5]uint8{GREY, GREY, GREY, GREY, GREY}
	previousChars := map[byte]uint8{}

	for i := 0; i < 5; i++ {
		if guess[i] == answer[i] {
			pattern[i] = GREEN
		} else if possibleAnswerChars[guess[i]] > 0 {
			if previousChars[guess[i]] < possibleAnswerChars[guess[i]]-matchingChars[guess[i]] {
				pattern[i] = YELLOW
			}
		}

		previousChars[guess[i]]++
	}

	return pattern
}

func updatePossibleAnswers(guess string, pattern [5]uint8, possibleAnswers []string) []string {
	var updatedPossibleAnswers []string

	for _, possibleAnswer := range possibleAnswers {
		possibleAnswerChars := map[byte]uint8{}
		guessChars := map[byte]uint8{}
		knownGuessChars := map[byte]uint8{}

		for i := 0; i < 5; i++ {
			possibleAnswerChars[possibleAnswer[i]]++
			guessChars[guess[i]]++

			if pattern[i] == GREEN || pattern[i] == YELLOW {
				knownGuessChars[guess[i]]++
			}
		}

		isPossibleAnswer := true

		for i := 0; i < 5; i++ {
			isMatching := possibleAnswer[i] == guess[i]

			if pattern[i] == GREEN {
				if !isMatching {
					isPossibleAnswer = false
					break
				}
			}

			if pattern[i] == YELLOW {
				if isMatching || possibleAnswerChars[guess[i]] < knownGuessChars[guess[i]] {
					isPossibleAnswer = false
					break
				}
			}

			if pattern[i] == GREY {
				if isMatching || possibleAnswerChars[guess[i]] >= guessChars[guess[i]] {
					isPossibleAnswer = false
					break
				}
			}
		}

		if isPossibleAnswer {
			updatedPossibleAnswers = append(updatedPossibleAnswers, possibleAnswer)
		}
	}

	return updatedPossibleAnswers
}
