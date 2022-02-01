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

func (d *Dictionary) UpdateRemainingPossibleAnswers(guess Word, feedback Feedback) {
	var remainingPossibleAnswers []Word

	for _, possibleAnswer := range d.RemainingPossibleAnswers {
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

		keep := true

		for i := 0; i < wordLength; i++ {
			matching := possibleAnswer[i] == guess[i]

			if feedback[i] == GREEN {
				if !matching {
					keep = false
					break
				}
			} else if feedback[i] == YELLOW {
				if matching || possibleAnswerChars[guess[i]] < greenYellowChars[guess[i]] {
					keep = false
					break
				}
			} else { // GREY
				if matching || possibleAnswerChars[guess[i]] >= guessChars[guess[i]] {
					keep = false
					break
				}
			}
		}

		if keep {
			remainingPossibleAnswers = append(remainingPossibleAnswers, possibleAnswer)
		}
	}

	d.RemainingPossibleAnswers = remainingPossibleAnswers
}
