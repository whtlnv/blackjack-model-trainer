package main

import (
	"time"

	"github.com/whtlnv/blackjack-model-trainer/pkg/blackjack"
	"github.com/whtlnv/blackjack-model-trainer/pkg/genetics"

	"github.com/whtlnv/blackjack-model-trainer/internal/randomizer"
)

func createPlayer(chromosome *genetics.Chromosome) blackjack.Playerish {
	strategy, error := blackjack.NewStrategy(chromosome.Raw())
	if error != nil {
		panic(error)
	}

	return blackjack.NewPlayer(strategy)
}

func sessionPlayers(candidates []*genetics.Candidate) []blackjack.Playerish {
	players := []blackjack.Playerish{}

	for _, candidate := range candidates {
		players = append(players, createPlayer(candidate.Chromosome))
	}

	return players
}

func candidateFitness(players []blackjack.Playerish, sequence [][]byte) []*genetics.Candidate {
	candidates := []*genetics.Candidate{}

	for _, player := range players {
		statistics := player.GetStatistics()

		candidates = append(candidates, &genetics.Candidate{
			Chromosome: genetics.NewChromosome(statistics.Strategy, sequence),
			Fitness:    BlackjackFitnessFunction(statistics),
		})
	}

	return candidates
}

func printResults(players []blackjack.Playerish, fittedPlayers []*genetics.Candidate) {
	bankrollSum := 0.0
	fitnessSum := 0.0

	maxBankroll := 0.0
	maxGamesPlayed := 0
	maxGamesWon := 0
	maxWinRate := 0.0
	maxFitness := 0.0

	for i, player := range players {
		statistics := player.GetStatistics()
		fitness := fittedPlayers[i].Fitness

		bankrollSum += statistics.Bankroll
		fitnessSum += fitness

		if statistics.Bankroll > maxBankroll {
			maxBankroll = statistics.Bankroll
		}

		if statistics.GamesPlayed > maxGamesPlayed {
			maxGamesPlayed = statistics.GamesPlayed
		}

		if statistics.GamesWon > maxGamesWon {
			maxGamesWon = statistics.GamesWon
		}

		winRate := float64(statistics.GamesWon) / float64(statistics.GamesPlayed)
		if winRate > maxWinRate {
			maxWinRate = winRate
		}

		if fitness > maxFitness {
			maxFitness = fitness
		}
	}

	averageBankroll := bankrollSum / float64(len(players))
	averageFitness := fitnessSum / float64(len(players))

	println("Average bankroll:", averageBankroll)
	println("Max bankroll:", maxBankroll)
	println("Max games played:", maxGamesPlayed)
	println("Max games won:", maxGamesWon)
	println("Max win rate:", maxWinRate)

	println("Average fitness:", averageFitness)
	println("Max fitness:", maxFitness)
}

func trainingSession() {
	generations := 100
	sequence := blackjack.GetSequencing()
	options := genetics.GenerationOptions{
		PopulationSize: 100,
		MutationRate:   0.1,
		CutoffRate:     0.2,
	}
	deckSize := 6
	penetration := 0.5
	handsPerGeneration := 1000

	seed := time.Now().UnixNano()
	randomizer := randomizer.NewRandomizer(seed)

	fittedPlayers := [](*genetics.Candidate){}
	for i := 0; i < generations; i++ {
		currentGen := genetics.NewGenerationFromPrevious(
			fittedPlayers,
			sequence,
			options,
			randomizer,
		)
		players := sessionPlayers(currentGen)

		shoe := blackjack.NewShoe(deckSize)
		shoe.SetPenetration(penetration)
		table := blackjack.NewTable(players, shoe)
		table.RunMany(handsPerGeneration)

		fittedPlayers = candidateFitness(table.Players, sequence)

		printResults(table.Players, fittedPlayers)

		println("Gen", i+1, "of", generations)
	}
}
