package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShoeInitialization(t *testing.T) {
	deckCount := 3
	shoe := NewShoe(deckCount)

	t.Run("Shoe has deckCount * 52 cards", func(t *testing.T) {
		got := shoe.Size()
		want := deckCount * 52
		assert.Equal(t, got, want)
	})

	t.Run("Shoe is shuffled", func(t *testing.T) {
		shoe2 := NewShoe(deckCount)
		assert.NotEqual(t, shoe.cards, shoe2.cards)
	})
}

func TestShoePeek(t *testing.T) {
	deckCount := 1
	shoe := NewShoe(deckCount)
	tenCards := shoe.Peek(0, 10)

	t.Run("Peek at the first 5 cards", func(t *testing.T) {
		drawCount := 5
		cards := shoe.Peek(0, 5)

		assert.Equal(t, len(cards), drawCount)
		assert.Equal(t, cards, tenCards[:drawCount])
	})

	t.Run("Peek at the next 4 cards", func(t *testing.T) {
		drawCount := 4
		cursor := 5
		cards := shoe.Peek(cursor, drawCount)

		assert.Equal(t, len(cards), drawCount)
		assert.Equal(t, cards, tenCards[cursor:cursor+drawCount])
	})

	t.Run("Should not peek past the end of the shoe", func(t *testing.T) {
		drawCount := 10
		cursor := 50
		cards := shoe.Peek(cursor, drawCount)

		assert.Equal(t, len(cards), 2)
	})
}
