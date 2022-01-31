package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type FeedbackColour uint8

const (
	GREY FeedbackColour = iota
	YELLOW
	GREEN
)

type Feedback [5]FeedbackColour

type FeedbackResolver interface {
	Resolve(guess Word) Feedback
}

type InteractiveFeedbackResolver struct {
	ReadWriter *bufio.ReadWriter
}

func (r *InteractiveFeedbackResolver) Resolve(guess Word) Feedback {
	r.ReadWriter.WriteString("Enter feedback [Grey=0, Yellow=1, Green=2]:\n")

	for {
		r.ReadWriter.Flush()

		input, err := r.ReadWriter.ReadBytes('\n')
		if err != nil {
			r.ReadWriter.WriteString("An unexpected error occurred. Please try again:\n")
			continue
		}
		input = bytes.TrimSpace(input)

		if len(input) != 5 {
			r.ReadWriter.WriteString(fmt.Sprintf("Expected 5 digits, received %d. Please try again:\n", len(input)))
			continue
		}

		feedback := Feedback{GREY, GREY, GREY, GREY, GREY}
		validFeedback := true

		for i := 0; i < 5; i++ {
			if input[i] == '2' {
				feedback[i] = GREEN
			} else if input[i] == '1' {
				feedback[i] = YELLOW
			} else if input[i] != '0' {
				validFeedback = false
				break
			}
		}

		if !validFeedback {
			r.ReadWriter.WriteString("Invalid feedback format. Please try again using the following as an example [01210]:\n")
			continue
		}

		return feedback
	}
}

type TestFeedbackResolver struct {
	Answer Word
}

func (r *TestFeedbackResolver) Resolve(guess Word) Feedback {
	return buildFeedback(guess, r.Answer)
}

func buildFeedback(guess, answer Word) Feedback {
	answerChars := map[byte]uint8{}
	matchingChars := map[byte]uint8{}

	for i := 0; i < 5; i++ {
		answerChars[answer[i]]++

		if guess[i] == answer[i] {
			matchingChars[guess[i]]++
		}
	}

	feedback := Feedback{GREY, GREY, GREY, GREY, GREY}
	previousChars := map[byte]uint8{}

	for i := 0; i < 5; i++ {
		if guess[i] == answer[i] {
			feedback[i] = GREEN
		} else if answerChars[guess[i]] > 0 {
			if previousChars[guess[i]] < answerChars[guess[i]]-matchingChars[guess[i]] {
				feedback[i] = YELLOW
			}
		}

		previousChars[guess[i]]++
	}

	return feedback
}
