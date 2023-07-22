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
	_, playerIsBusted := playerHand.Score()
	if playerIsBusted {
		return Stand
	}

	return Hit
}

func (strategy *strategyMock) Bet() int {
	return 1
}

func (strategy *strategyMock) GetInitialBankroll() int {
	return strategy.initialBankroll
}

type shoeMock struct{}

func (shoe *shoeMock) Size() int {
	return 52
}

func (shoe *shoeMock) Peek(count int) []Card {
	peek := []Card{
		NewCard(Three, Clubs),
		NewCard(Three, Hearts),
		NewCard(Three, Diamonds),
		NewCard(Three, Spades),
		NewCard(Four, Clubs),
		NewCard(Four, Hearts),
	}
	return peek
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

func TestPlayerPlay(t *testing.T) {
	t.Run("Should return the number of cards dealt", func(t *testing.T) {
		strategy := &strategyMock{}
		strategy.initialBankroll = 100
		player := NewPlayer(strategy)

		shoe := &shoeMock{}

		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Ten, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 5, cardsTaken)
	})
}

// Test game resolution
