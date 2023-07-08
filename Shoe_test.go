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
		assert.Equal(t, want, got)
	})

	t.Run("Shoe is shuffled", func(t *testing.T) {
		shoe2 := NewShoe(deckCount)
		assert.NotEqual(t, shoe.cards, shoe2.cards)
	})
}

func TestShoePeek(t *testing.T) {
	deckCount := 1
	shoe := NewShoe(deckCount)
	tenCards := shoe.Peek(10)

	t.Run("Peek at the first 5 cards", func(t *testing.T) {
		drawCount := 5
		cards := shoe.Peek(5)

		assert.Equal(t, len(cards), drawCount)
		assert.Equal(t, cards, tenCards[:drawCount])
	})

	shoe.AdvanceCursor(5)

	t.Run("Peek at the next 4 cards", func(t *testing.T) {
		drawCount := 4
		offset := 5
		cards := shoe.Peek(drawCount)

		assert.Equal(t, len(cards), drawCount)
		assert.Equal(t, cards, tenCards[offset:offset+drawCount])
	})

	shoe.AdvanceCursor(45) // cursor is now at 50

	t.Run("Should not peek past the end of the shoe", func(t *testing.T) {
		drawCount := 10

		cards := shoe.Peek(drawCount)

		assert.Equal(t, len(cards), 2)
	})
}

func TestShoeAdvanceCursor(t *testing.T) {
	deckCount := 1
	shoe := NewShoe(deckCount)

	t.Run("Should advance the cursor by the offset", func(t *testing.T) {
		offset := 5
		cursor, err := shoe.AdvanceCursor(offset)

		assert.Equal(t, cursor, offset)
		assert.NoError(t, err)
	})

	t.Run("Should return an error if cursor > length", func(t *testing.T) {
		offset := 53
		cursor, err := shoe.AdvanceCursor(offset)

		assert.Equal(t, cursor, 5)
		assert.Error(t, err)
	})
}
