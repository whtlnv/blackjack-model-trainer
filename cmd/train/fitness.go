package main

import (
	"github.com/whtlnv/blackjack-model-trainer/pkg/blackjack"
	"github.com/whtlnv/blackjack-model-trainer/pkg/genetics"
)

func BlackjackFitnessFunction(player *blackjack.Player, sequencing [][]byte) *genetics.Candidate {
	playerDNA := player.GetStrategy().GetEncodedStrategy()
	chromosome := genetics.NewChromosome(playerDNA, sequencing)

	statistics := player.GetStatistics()

	return &genetics.Candidate{
		Chromosome: chromosome,
		Fitness:    statistics.Bankroll,
	}
}
