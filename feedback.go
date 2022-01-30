package main

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

type StdinFeedbackResolver struct{}

func (r *StdinFeedbackResolver) Resolve(guess Word) Feedback {
	panic("not implemented")
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
