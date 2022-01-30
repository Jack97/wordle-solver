package main

import (
	"math"
	"math/rand"
)

type Game struct {
	Dictionary       *Dictionary
	FeedbackResolver FeedbackResolver
}

type GameResult struct {
	Win     bool
	Guesses []Word
}

func (g *Game) Play() GameResult {
	var (
		win     bool
		guesses []Word
	)

	g.Dictionary.ResetRemainingPossibleAnswers()

	for i := 0; i < 6; i++ {
		guess := g.guess()
		guesses = append(guesses, guess)
		feedback := g.FeedbackResolver.Resolve(guess)

		if g.isComplete(feedback) {
			win = true
			break
		}

		g.Dictionary.UpdateRemainingPossibleAnswers(guess, feedback)
	}

	return GameResult{Win: win, Guesses: guesses}
}

func (g *Game) guess() Word {
	if len(g.Dictionary.RemainingPossibleAnswers) <= 2 {
		// Return a random answer if the probability of it being correct is >= 50%
		return g.Dictionary.RemainingPossibleAnswers[rand.Intn(len(g.Dictionary.RemainingPossibleAnswers))]
	}

	var (
		maxInformationGain float64
		optimalGuess       Word
	)

	for _, guess := range g.Dictionary.ValidGuesses() {
		feedbackProbabilities := map[Feedback]float64{}

		for _, possibleAnswer := range g.Dictionary.RemainingPossibleAnswers {
			feedback := buildFeedback(guess, possibleAnswer)
			feedbackProbabilities[feedback] += 1.0 / float64(len(g.Dictionary.RemainingPossibleAnswers))
		}

		// Calculate the entropy to determine the guess that will yield the most information
		var sum float64
		for _, probability := range feedbackProbabilities {
			sum += probability * math.Log(probability)
		}
		informationGain := -sum

		if informationGain > maxInformationGain {
			maxInformationGain = informationGain
			optimalGuess = guess
		}
	}

	return optimalGuess
}

func (g *Game) isComplete(feedback Feedback) bool {
	for i := 0; i < 5; i++ {
		if feedback[i] != GREEN {
			return false
		}
	}

	return true
}
