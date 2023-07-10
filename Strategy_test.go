package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrategyParsing(t *testing.T) {
	t.Run("Should parse a strategy from string", func(t *testing.T) {
		// Hit every time
		raw := bytes.Repeat([]byte("H"), 350)
		strategy := NewStrategy(raw)

		playerHand := Hand{Card{Ace, Spades}, Card{Ace, Hearts}}
		dealerHand := Hand{Card{Two, Clubs}}

		got := strategy.Play(playerHand, dealerHand)
		want := Hit

		assert.Equal(t, want, got)
	})
}
