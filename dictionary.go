package main

import "fmt"

type Word [5]byte

func (w *Word) UnmarshalText(data []byte) error {
	if len(data) != 5 {
		return fmt.Errorf("invalid word length: %s", data)
	}

	for i := 0; i < 5; i++ {
		w[i] = data[i]
	}

	return nil
}

type Dictionary struct {
	AcceptedGuesses          []Word `json:"acceptedGuesses"`
	PossibleAnswers          []Word `json:"possibleAnswers"`
	RemainingPossibleAnswers []Word
}

func (d Dictionary) ValidGuesses() []Word {
	return append(d.AcceptedGuesses, d.RemainingPossibleAnswers...)
}

func (d *Dictionary) ResetRemainingPossibleAnswers() {
	d.RemainingPossibleAnswers = d.PossibleAnswers
}

func (d *Dictionary) UpdateRemainingPossibleAnswers(guess Word, feedback Feedback) {
	var remainingPossibleAnswers []Word

	for _, possibleAnswer := range d.RemainingPossibleAnswers {
		possibleAnswerChars := map[byte]uint8{}
		guessChars := map[byte]uint8{}
		knownGuessChars := map[byte]uint8{}

		for i := 0; i < 5; i++ {
			possibleAnswerChars[possibleAnswer[i]]++
			guessChars[guess[i]]++

			if feedback[i] == GREEN || feedback[i] == YELLOW {
				knownGuessChars[guess[i]]++
			}
		}

		isPossibleAnswer := true

		for i := 0; i < 5; i++ {
			isSameChar := possibleAnswer[i] == guess[i]

			if feedback[i] == GREEN {
				if !isSameChar {
					isPossibleAnswer = false
					break
				}
			} else if feedback[i] == YELLOW {
				if isSameChar || possibleAnswerChars[guess[i]] < knownGuessChars[guess[i]] {
					isPossibleAnswer = false
					break
				}
			} else { // GREY
				if isSameChar || possibleAnswerChars[guess[i]] >= guessChars[guess[i]] {
					isPossibleAnswer = false
					break
				}
			}
		}

		if isPossibleAnswer {
			remainingPossibleAnswers = append(remainingPossibleAnswers, possibleAnswer)
		}
	}

	d.RemainingPossibleAnswers = remainingPossibleAnswers
}
