package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/whtlnv/blackjack-model-trainer/internal/randomizer"
	"github.com/whtlnv/blackjack-model-trainer/pkg/blackjack"
	"github.com/whtlnv/blackjack-model-trainer/pkg/genetics"
)

func TestBlackjackFitnessFunction(t *testing.T) {
	sequencing := [][]byte{}
	randomizer := randomizer.NewRandomizer(1)

	chromosome := genetics.NewRandomChromosome(sequencing, randomizer)
	strategy, _ := blackjack.NewStrategy(chromosome.Raw())

	t.Run("A player should rank higher if they have a higher bankroll", func(t *testing.T) {
		players := []*blackjack.Player{
			blackjack.NewPlayer(strategy),
			blackjack.NewPlayer(strategy),
		}

		players[0].Bankroll = 100
		players[1].Bankroll = 200

		candidate1 := BlackjackFitnessFunction(players[0], sequencing)
		candidate2 := BlackjackFitnessFunction(players[1], sequencing)

		assert.Greater(t, candidate2.Fitness, candidate1.Fitness)
	})
	// t.Run("A player should rank higher if they have a higher win rate", func(t *testing.T) {})
	// t.Run("If a player has bankroll > 0, they should rank higher if they have fewer games played/won", func(t *testing.T) {})
	// t.Run("If a player has bankroll < 0, they should rank higher if they have more games played", func(t *testing.T) {})
}
