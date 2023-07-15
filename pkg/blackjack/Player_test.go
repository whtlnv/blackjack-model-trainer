package blackjack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helpers

type strategyMock struct {
	initialBankroll int
}

func (strategy *strategyMock) Play(playerHand Hand, dealerHand Hand) PlayerAction {
	return Hit
}

func (strategy *strategyMock) Bet() int {
	return 1
}

func (strategy *strategyMock) GetInitialBankroll() int {
	return strategy.initialBankroll
}

// Tests

func TestPlayerInitialization(t *testing.T) {
	t.Run("Should set player initial bankroll", func(t *testing.T) {
		strategy := &strategyMock{}
		strategy.initialBankroll = 100
		player := NewPlayer(strategy)

		assert.Equal(t, strategy.initialBankroll, player.Bankroll)
	})
}

func TestPlayerBet(t *testing.T) {
	strategy := &strategyMock{}
	strategy.initialBankroll = 100

	t.Run("Should decide to play a hand if has funds", func(t *testing.T) {
		player := NewPlayer(strategy)

		willBet, ammount := player.Bet( /* shoe? */ )
		assert.True(t, willBet)
		assert.Equal(t, 1, ammount)
	})

	t.Run("Should decide not to play a hand if has no funds", func(t *testing.T) {
		player := NewPlayer(strategy)
		player.Bankroll = 0

		willBet, _ := player.Bet( /* shoe? */ )
		assert.False(t, willBet)
	})

	t.Run("Bet should be deducted from bankroll", func(t *testing.T) {
		player := NewPlayer(strategy)
		player.Bankroll = 10

		willBet, _ := player.Bet( /* shoe? */ )
		assert.True(t, willBet)
		assert.Equal(t, 9, player.Bankroll)
	})
}
