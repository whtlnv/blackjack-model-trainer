package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerTurn(t *testing.T) {
	t.Run("Should decide to play a hand", func(t *testing.T) {
		rawStrategy := bytes.Repeat([]byte("H"), 340)
		player := NewPlayer(NewStrategy(rawStrategy))
		willBet, ammount := player.ShouldBet( /* shoe? */ )
		assert.True(t, willBet)
		assert.Equal(t, 1, ammount)
	})
}
