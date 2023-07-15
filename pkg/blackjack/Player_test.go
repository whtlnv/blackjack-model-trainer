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

func TestPlayerBet(t *testing.T) {
	alwaysHitStrategy := bytes.Repeat([]byte("H"), HandCount)
	t.Run("Should decide to play a hand if has funds", func(t *testing.T) {
		rawStrategy := append(alwaysHitStrategy, []byte("00100001")...)
		strategy, _ := NewStrategy(rawStrategy)
		player := NewPlayer(strategy)

		willBet, ammount := player.Bet( /* shoe? */ )
		assert.True(t, willBet)
		assert.Equal(t, 1, ammount)
	})

	t.Run("Should decide not to play a hand if has no funds", func(t *testing.T) {
		rawStrategy := append(alwaysHitStrategy, []byte("00000001")...)
		strategy, _ := NewStrategy(rawStrategy)
		player := NewPlayer(strategy)
		player.Bankroll = 0 // Redundant, but explicit

		willBet, _ := player.Bet( /* shoe? */ )
		assert.False(t, willBet)
	})

	t.Run("Bet should be deducted from bankroll", func(t *testing.T) {
		rawStrategy := append(alwaysHitStrategy, []byte("00100001")...)
		strategy, _ := NewStrategy(rawStrategy)
		player := NewPlayer(strategy)
		player.Bankroll = 10

		willBet, _ := player.Bet( /* shoe? */ )
		assert.True(t, willBet)
		assert.Equal(t, 9, player.Bankroll)
	})
}
