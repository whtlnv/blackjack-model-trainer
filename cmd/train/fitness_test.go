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
	t.Run("A player should rank higher if they have a higher bankroll change", func(t *testing.T) {
		playerA := &playerMock{}
		playerA.On("GetStatistics").Return(PlayerStatistics{
			GamesPlayed: 100,
			GamesWon:    99,
			GamesLost:   1,

			InitialBankroll: 20,
			Bankroll:        100,
			BankrollDelta:   80,
		})

		playerB := &playerMock{}
		playerB.On("GetStatistics").Return(PlayerStatistics{
			GamesPlayed: 100,
			GamesWon:    99,
			GamesLost:   1,

			InitialBankroll: 10,
			Bankroll:        100,
			BankrollDelta:   90,
		})

		fitnessA := BlackjackFitnessFunction(playerA)
		fitnessB := BlackjackFitnessFunction(playerB)

		assert.Greater(t, fitnessB, fitnessA)
	})

	t.Run("A player should rank higher if they have a higher win rate", func(t *testing.T) {
		playerA := &playerMock{}
		playerA.On("GetStatistics").Return(PlayerStatistics{
			GamesPlayed:   100,
			GamesWon:      10,
			GamesLost:     90,
			Bankroll:      200,
			BankrollDelta: 100,
		})

		playerB := &playerMock{}
		playerB.On("GetStatistics").Return(PlayerStatistics{
			GamesPlayed:   100,
			GamesWon:      20,
			GamesLost:     80,
			Bankroll:      200,
			BankrollDelta: 100,
		})

		fitnessA := BlackjackFitnessFunction(playerA)
		fitnessB := BlackjackFitnessFunction(playerB)

		assert.Greater(t, fitnessB, fitnessA)
	})

	t.Run("A player should still rank higher if they have more money, even if wr is worst", func(t *testing.T) {
		playerA := &playerMock{}
		playerA.On("GetStatistics").Return(PlayerStatistics{
			GamesPlayed:   100,
			GamesWon:      99,
			GamesLost:     1,
			Bankroll:      200,
			BankrollDelta: 100,
		})

		playerB := &playerMock{}
		playerB.On("GetStatistics").Return(PlayerStatistics{
			GamesPlayed:   100,
			GamesWon:      1,
			GamesLost:     99,
			Bankroll:      201,
			BankrollDelta: 101,
		})

		fitnessA := BlackjackFitnessFunction(playerA)
		fitnessB := BlackjackFitnessFunction(playerB)

		assert.Greater(t, fitnessB, fitnessA)
	})

	t.Run("If a player has bankroll > 0, they should rank higher if they have fewer games played", func(t *testing.T) {
		playerA := &playerMock{}
		playerA.On("GetStatistics").Return(PlayerStatistics{
			GamesPlayed:   20,
			GamesWon:      10,
			GamesLost:     10,
			Bankroll:      200,
			BankrollDelta: 100,
		})

		playerB := &playerMock{}
		playerB.On("GetStatistics").Return(PlayerStatistics{
			GamesPlayed:   10,
			GamesWon:      5,
			GamesLost:     5,
			Bankroll:      200,
			BankrollDelta: 100,
		})

		fitnessA := BlackjackFitnessFunction(playerA)
		fitnessB := BlackjackFitnessFunction(playerB)

		assert.Greater(t, fitnessB, fitnessA)
	})

	t.Run("If a player has bankroll == 0, they should rank higher if they have more games played", func(t *testing.T) {
		playerA := &playerMock{}
		playerA.On("GetStatistics").Return(PlayerStatistics{
			GamesPlayed:   10,
			GamesWon:      5,
			GamesLost:     5,
			Bankroll:      0,
			BankrollDelta: -100,
		})

		playerB := &playerMock{}
		playerB.On("GetStatistics").Return(PlayerStatistics{
			GamesPlayed:   20,
			GamesWon:      10,
			GamesLost:     10,
			Bankroll:      0,
			BankrollDelta: -100,
		})

		fitnessA := BlackjackFitnessFunction(playerA)
		fitnessB := BlackjackFitnessFunction(playerB)

		assert.Greater(t, fitnessB, fitnessA)
	})
}
