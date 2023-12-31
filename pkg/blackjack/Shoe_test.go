package blackjack

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

func TestShoeShuffle(t *testing.T) {
	deckCount := 1
	shoe := NewShoe(deckCount)
	originalCards := shoe.cards

	t.Run("Shoe should shuffle", func(t *testing.T) {
		shoe.Shuffle()
		assert.NotEqual(t, originalCards, shoe.cards)
	})

	t.Run("Should reset the reshuffle flag", func(t *testing.T) {
		shoe.needsReshuffle = true
		shoe.Shuffle()
		assert.False(t, shoe.needsReshuffle)
	})

	t.Run("Should reset the cursor", func(t *testing.T) {
		shoe.cursor = 5
		shoe.Shuffle()
		assert.Equal(t, 0, shoe.cursor)
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

func TestShoePeekIndex(t *testing.T) {
	deckCount := 1
	shoe := NewShoe(deckCount)

	t.Run("Peek at 5th card", func(t *testing.T) {
		index := 5
		card, error := shoe.PeekAtIndex(index)

		assert.Equal(t, card, shoe.cards[index])
		assert.NoError(t, error)
	})

	t.Run("Peek at the 4th card after advancing cursor", func(t *testing.T) {
		offset := 5
		index := 4
		shoe.AdvanceCursor(offset)
		card, error := shoe.PeekAtIndex(index)

		assert.Equal(t, card, shoe.cards[offset+index])
		assert.NoError(t, error)
	})

	t.Run("Should not peek past the end of the shoe", func(t *testing.T) {
		index := 3

		shoe.AdvanceCursor(45) // cursor is now at 50
		card, error := shoe.PeekAtIndex(index)

		assert.Equal(t, card, Card{})
		assert.Error(t, error)
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

func TestPenetration(t *testing.T) {
	deckCount := 1

	t.Run("Should set reshuffle flag when penetration is reached", func(t *testing.T) {
		shoe := NewShoe(deckCount)
		shoe.SetPenetration(.5)
		shoe.AdvanceCursor(25)
		assert.False(t, shoe.NeedsReshuffle())

		shoe.AdvanceCursor(1)
		assert.True(t, shoe.NeedsReshuffle())
	})

	t.Run("Should reset reshuffle flag when shoe is reshuffled", func(t *testing.T) {
		shoe := NewShoe(deckCount)
		shoe.SetPenetration(.5)
		shoe.AdvanceCursor(26)
		assert.True(t, shoe.NeedsReshuffle())

		shoe.Shuffle()
		assert.False(t, shoe.NeedsReshuffle())
	})
}
