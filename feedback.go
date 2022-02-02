package main

import (
	"bufio"
	"fmt"
	"log"
)

type FeedbackColour uint8

const (
	GREY FeedbackColour = iota
	YELLOW
	GREEN
)

type Feedback [wordLength]FeedbackColour

func (f *Feedback) Success() bool {
	for _, colour := range f {
		if colour != GREEN {
			return false
		}
	}

	return true
}

type FeedbackResolver interface {
	Resolve(guess Word, attempt int) Feedback
}

type InteractiveFeedbackResolver struct {
	Logger  *log.Logger
	Scanner *bufio.Scanner
}

func (r *InteractiveFeedbackResolver) Resolve(guess Word, attempt int) Feedback {
	r.Logger.Printf("Guess #%d: %s", attempt, guess)
	r.Logger.Printf("Enter feedback (Grey = 0, Yellow = 1, Green = 2):")

	for {
		r.Scanner.Scan()

		feedback, err := r.parseFeedback(r.Scanner.Bytes())
		if err != nil {
			r.Logger.Printf("Feedback validation failed: %s", err)
			r.Logger.Printf("Please try again (e.g. 21012):")
			continue
		}

		return feedback
	}
}

func (r *InteractiveFeedbackResolver) parseFeedback(data []byte) (Feedback, error) {
	feedback := Feedback{GREY, GREY, GREY, GREY, GREY}

	if len(data) != 5 {
		return feedback, fmt.Errorf("invalid length")
	}

	for i := 0; i < wordLength; i++ {
		char := data[i]

		switch char {
		case '2':
			feedback[i] = GREEN
		case '1':
			feedback[i] = YELLOW
		case '0':
			// Do nothing
		default:
			return feedback, fmt.Errorf("invalid character '%c'", char)
		}
	}

	return feedback, nil
}

type SimulationFeedbackResolver struct {
	Answer Word
}

func (r *SimulationFeedbackResolver) Resolve(guess Word, _ int) Feedback {
	return buildFeedback(guess, r.Answer)
}

func buildFeedback(guess, answer Word) Feedback {
	feedback := Feedback{GREY, GREY, GREY, GREY, GREY}

	answerChars := map[byte]int{}
	greenChars := map[byte]int{}
	yellowChars := map[byte]int{}

	for i := 0; i < wordLength; i++ {
		answerChars[answer[i]]++

		if guess[i] == answer[i] {
			feedback[i] = GREEN
			greenChars[guess[i]]++
		}
	}

	for i := 0; i < wordLength; i++ {
		if feedback[i] == GREEN {
			continue
		}

		// The character is present in the answer, and it's not been seen as many times as it occurs
		if answerChars[guess[i]] > 0 && yellowChars[guess[i]] < answerChars[guess[i]]-greenChars[guess[i]] {
			feedback[i] = YELLOW
			yellowChars[guess[i]]++
		}
	}

	return feedback
}
