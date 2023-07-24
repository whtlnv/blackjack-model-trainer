package blackjack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameInitialization(t *testing.T) {
	t.Run("Should set game initial bet", func(t *testing.T) {
		bet := 1
		game := NewGame(bet)

		assert.Equal(t, bet, game.bet)
	})
}

func TestGameSetHand(t *testing.T) {
	t.Run("Should set game hand", func(t *testing.T) {
		game := NewGame(1)

		hand := Hand{NewCard(Two, Spades), NewCard(Three, Spades)}
		game.SetHand(hand)

		assert.Equal(t, hand, game.hand)
	})
}

func TestGameHit(t *testing.T) {
	t.Run("Should hit a hand", func(t *testing.T) {
		game := NewGame(1)
		hand := Hand{NewCard(Two, Spades), NewCard(Three, Spades)}
		game.SetHand(hand)

		game.Hit(NewCard(Four, Spades))

		want := Hand{NewCard(Two, Spades), NewCard(Three, Spades), NewCard(Four, Spades)}
		got := game.hand

		assert.Equal(t, want, got)
	})
}

func TestGameDouble(t *testing.T) {
	t.Run("Should double hand bet", func(t *testing.T) {
		game := NewGame(1)
		hand := Hand{NewCard(Two, Spades), NewCard(Three, Spades)}
		game.SetHand(hand)

		game.Double(NewCard(Four, Spades))

		want := 2
		got := game.bet

		assert.Equal(t, want, got)
	})

	t.Run("Should hit once", func(t *testing.T) {
		game := NewGame(1)
		hand := Hand{NewCard(Two, Spades), NewCard(Three, Spades)}
		game.SetHand(hand)

		game.Double(NewCard(Four, Spades))

		want := Hand{NewCard(Two, Spades), NewCard(Three, Spades), NewCard(Four, Spades)}
		got := game.hand

		assert.Equal(t, want, got)
	})

	t.Run("Should flag a hand that was doubled", func(t *testing.T) {
		game := NewGame(1)
		hand := Hand{NewCard(Two, Spades), NewCard(Three, Spades)}
		game.SetHand(hand)

		game.Double(NewCard(Four, Spades))

		want := true
		got := game.Doubled

		assert.Equal(t, want, got)
	})

}
