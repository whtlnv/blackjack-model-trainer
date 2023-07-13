package blackjack

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerInitialization(t *testing.T) {
	t.Run("Should set player initial bankroll", func(t *testing.T) {
		rawStrategy := bytes.Repeat([]byte("H"), HandCount)
		strategy, _ := NewStrategy(rawStrategy)
		player := NewPlayer(strategy)

		assert.Equal(t, strategy.InitialBankroll, player.Bankroll)
	})
}

func TestPlayerTurn(t *testing.T) {
	t.Run("Should decide to play a hand", func(t *testing.T) {
		rawStrategy := bytes.Repeat([]byte("H"), HandCount)
		strategy, _ := NewStrategy(rawStrategy)
		player := NewPlayer(strategy)

		willBet, ammount := player.ShouldBet( /* shoe? */ )
		assert.True(t, willBet)
		assert.Equal(t, 1, ammount)
	})
}
