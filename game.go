package main

import (
	"math"
	"math/rand"
)

type GameResult struct {
	Win     bool
	Guesses []Word
}

type Game struct {
	Dictionary       *Dictionary
	FeedbackResolver FeedbackResolver
}

const maxAttempts = 6

// The calculation to determine the optimal first guess is extremely expensive
// and will always produce the same result. Calculating the guess ahead of time
// drastically improves the solver's performance.
var optimalFirstGuess = Word{'s', 'o', 'a', 'r', 'e'}

func (g *Game) Play() (*GameResult, error) {
	var (
		win     bool
		guesses []Word
	)

	g.Dictionary.ResetRemainingPossibleAnswers()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		guess := g.guess(attempt)
		guesses = append(guesses, guess)

		feedback := g.FeedbackResolver.Resolve(guess, attempt)

		if feedback.Success() {
			win = true
			break
		}

		err := g.Dictionary.UpdateRemainingPossibleAnswers(guess, feedback)
		if err != nil {
			return nil, err
		}
	}

	return &GameResult{Win: win, Guesses: guesses}, nil
}

func (g *Game) guess(attempt int) Word {
	if attempt == 1 {
		return optimalFirstGuess
	}

	// Return a random answer if the probability of it being correct is >= 50%
	if len(g.Dictionary.RemainingPossibleAnswers) <= 2 {
		n := rand.Intn(len(g.Dictionary.RemainingPossibleAnswers))
		return g.Dictionary.RemainingPossibleAnswers[n]
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
