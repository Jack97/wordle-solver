package main

import "fmt"

const wordLength = 5

type Word [wordLength]byte

func (w *Word) UnmarshalText(data []byte) error {
	if len(data) != wordLength {
		return fmt.Errorf("invalid word length '%s'", data)
	}

	for i := 0; i < wordLength; i++ {
		w[i] = data[i]
	}

	return nil
}

type Dictionary struct {
	AcceptedGuesses          []Word `json:"acceptedGuesses"`
	PossibleAnswers          []Word `json:"possibleAnswers"`
	RemainingPossibleAnswers []Word
}

func (d *Dictionary) ValidGuesses() []Word {
	return append(d.AcceptedGuesses, d.RemainingPossibleAnswers...)
}

func (d *Dictionary) ResetRemainingPossibleAnswers() {
	d.RemainingPossibleAnswers = d.PossibleAnswers
}

func (d *Dictionary) UpdateRemainingPossibleAnswers(guess Word, feedback Feedback) error {
	var remainingPossibleAnswers []Word

	for _, possibleAnswer := range d.RemainingPossibleAnswers {
		keep := true

		possibleAnswerChars := map[byte]int{}
		guessChars := map[byte]int{}
		greenYellowChars := map[byte]int{}

		for i := 0; i < wordLength; i++ {
			possibleAnswerChars[possibleAnswer[i]]++
			guessChars[guess[i]]++

			if feedback[i] == GREEN || feedback[i] == YELLOW {
				greenYellowChars[guess[i]]++
			}
		}

		for i := 0; i < wordLength; i++ {
			match := possibleAnswer[i] == guess[i]

			if feedback[i] == GREEN {
				if !match {
					keep = false
					break
				}
			} else if feedback[i] == YELLOW {
				if match || possibleAnswerChars[guess[i]] < greenYellowChars[guess[i]] {
					keep = false
					break
				}
			} else { // GREY
				if match || possibleAnswerChars[guess[i]] >= guessChars[guess[i]] {
					keep = false
					break
				}
			}
		}

		if keep {
			remainingPossibleAnswers = append(remainingPossibleAnswers, possibleAnswer)
		}
	}

	if len(remainingPossibleAnswers) == 0 {
		return fmt.Errorf("no remaining possible answers")
	}

	d.RemainingPossibleAnswers = remainingPossibleAnswers

	return nil
}
