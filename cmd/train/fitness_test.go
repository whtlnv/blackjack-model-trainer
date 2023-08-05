package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helpers

type playerMock struct {
	mock.Mock
}

func (p *playerMock) GetStatistics() PlayerStatistics {
	args := p.Called()
	return args.Get(0).(PlayerStatistics)
}

// Tests

func TestBlackjackFitnessFunction(t *testing.T) {
	t.Run("A player should rank higher if they have a higher bankroll", func(t *testing.T) {
		playerA := &playerMock{}
		playerA.On("GetStatistics").Return(PlayerStatistics{Bankroll: 100, GamesPlayed: 100})

		playerB := &playerMock{}
		playerB.On("GetStatistics").Return(PlayerStatistics{Bankroll: 200, GamesPlayed: 100})

		fitnessA := BlackjackFitnessFunction(playerA)
		fitnessB := BlackjackFitnessFunction(playerB)

		assert.Greater(t, fitnessB, fitnessA)
	})

	t.Run("A player should rank higher if they have a higher win rate", func(t *testing.T) {
		playerA := &playerMock{}
		playerA.On("GetStatistics").Return(PlayerStatistics{
			Bankroll:    100,
			GamesPlayed: 100,
			GamesWon:    10,
			GamesLost:   90,
		})

		playerB := &playerMock{}
		playerB.On("GetStatistics").Return(PlayerStatistics{
			Bankroll:    100,
			GamesPlayed: 100,
			GamesWon:    20,
			GamesLost:   80,
		})

		fitnessA := BlackjackFitnessFunction(playerA)
		fitnessB := BlackjackFitnessFunction(playerB)

		assert.Greater(t, fitnessB, fitnessA)
	})

	t.Run("A player should still rank higher if they have more money", func(t *testing.T) {
		playerA := &playerMock{}
		playerA.On("GetStatistics").Return(PlayerStatistics{
			Bankroll:    100,
			GamesPlayed: 100,
			GamesWon:    99,
			GamesLost:   1,
		})

		playerB := &playerMock{}
		playerB.On("GetStatistics").Return(PlayerStatistics{
			Bankroll:    101,
			GamesPlayed: 100,
			GamesWon:    1,
			GamesLost:   99,
		})

		fitnessA := BlackjackFitnessFunction(playerA)
		fitnessB := BlackjackFitnessFunction(playerB)

		assert.Greater(t, fitnessB, fitnessA)
	})

	t.Run("A player should rank higher if they have a smaller initial bankroll", func(t *testing.T) {
		// players := []*blackjack.Player{
		// 	blackjack.NewPlayer(strategy0),
		// 	blackjack.NewPlayer(strategy0),
		// }

		// players[0].Bankroll = 100
		// players[0].GamesPlayed = 100
		// players[1].Bankroll = 200
		// players[1].GamesPlayed = 100

		// candidate0 := BlackjackFitnessFunction(players[0])
		// candidate1 := BlackjackFitnessFunction(players[1])

		// assert.Greater(t, candidate1.Fitness, candidate0.Fitness)
	})

	// t.Run("If a player has bankroll > 0, they should rank higher if they have fewer games played/won", func(t *testing.T) {})

	// t.Run("If a player has bankroll < 0, they should rank higher if they have more games played", func(t *testing.T) {})
}
