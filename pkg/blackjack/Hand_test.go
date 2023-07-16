package blackjack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDealingACard(t *testing.T) {
	t.Run("Should deal a card to a hand", func(t *testing.T) {
		hand := Hand{}
		card := Card{Queen, Spades}

		hand.Deal(card)

		assert.Equal(t, 1, len(hand))
		assert.Equal(t, card, hand[0])
	})

	t.Run("Should append a card to a hand", func(t *testing.T) {
		hand := Hand{{King, Spades}}
		card := Card{Queen, Spades}

		hand.Deal(card)

		assert.Equal(t, 2, len(hand))
		assert.Equal(t, card, hand[1])
	})
}

func TestHandValues(t *testing.T) {
	t.Run("Should return the value of a hand", func(t *testing.T) {
		hand := Hand{
			Card{Two, Spades},
			Card{Seven, Hearts},
			Card{Queen, Diamonds},
		}

		got := hand.Values()
		want := []int{19}

		assert.ElementsMatch(t, want, got)
	})

	t.Run("Should return the correct score for Ace", func(t *testing.T) {
		hand := Hand{
			Card{Ace, Spades},
			Card{King, Hearts},
			Card{Queen, Diamonds},
		}

		got := hand.Values()
		want := []int{21, 31}

		assert.Equal(t, want, got)
	})

	t.Run("Should return the correct score for several Aces", func(t *testing.T) {
		hand := Hand{
			Card{Ace, Spades},
			Card{Ace, Spades},
			Card{Three, Diamonds},
		}

		got := hand.Values()
		want := []int{5, 15, 25}

		assert.Equal(t, want, got)
	})
}

func TestHandScore(t *testing.T) {
	t.Run("Should return the highest score of a hand", func(t *testing.T) {
		hand := Hand{
			Card{Two, Spades},
			Card{Seven, Hearts},
			Card{Queen, Diamonds},
		}

		got, isBusted := hand.Score()
		want := HandScore{19, 19}

		assert.Equal(t, want, got)
		assert.Equal(t, false, isBusted)
	})

	t.Run("Should return the correct score for Ace", func(t *testing.T) {
		hand := Hand{
			Card{Ace, Spades},
			Card{King, Hearts},
		}

		got, isBusted := hand.Score()
		want := HandScore{11, 21}

		assert.Equal(t, want, got)
		assert.Equal(t, false, isBusted)
	})

	t.Run("Should return the correct score for several Aces", func(t *testing.T) {
		hand := Hand{
			Card{Ace, Spades},
			Card{Ace, Spades},
			Card{Three, Diamonds},
		}

		got, isBusted := hand.Score()
		want := HandScore{5, 15}

		assert.Equal(t, want, got)
		assert.Equal(t, false, isBusted)
	})

	t.Run("Should return true if busted", func(t *testing.T) {
		hand := Hand{
			Card{Queen, Spades},
			Card{Queen, Spades},
			Card{Queen, Diamonds},
		}

		got, isBusted := hand.Score()
		want := HandScore{30, 30}

		assert.Equal(t, want, got)
		assert.Equal(t, true, isBusted)
	})
}

func TestIsPair(t *testing.T) {
	t.Run("Should return true if pair", func(t *testing.T) {
		hand := Hand{
			Card{Queen, Spades},
			Card{Queen, Spades},
		}

		got := hand.IsPair()
		want := true

		assert.Equal(t, want, got)
	})

	t.Run("Should return false if not pair", func(t *testing.T) {
		hand := Hand{
			Card{Queen, Spades},
			Card{King, Spades},
		}

		got := hand.IsPair()
		want := false

		assert.Equal(t, want, got)
	})
}

func TestHasSoftValue(t *testing.T) {
	t.Run("Should return true if soft value", func(t *testing.T) {
		hand := Hand{
			Card{Ace, Spades},
			Card{King, Spades},
		}

		got := hand.HasSoftValue()
		want := true

		assert.Equal(t, want, got)
	})

	t.Run("Should return false if not soft value", func(t *testing.T) {
		hand := Hand{
			Card{Queen, Spades},
			Card{King, Spades},
		}

		got := hand.HasSoftValue()
		want := false

		assert.Equal(t, want, got)
	})
}

// This is for player
// func TestHandSplit(t *testing.T) {
// 	t.Run("Should split a hand", func(t *testing.T) {
// 		hand := Hand{
// 			Card{Queen, Spades},
// 			Card{Queen, Spades},
// 		}

// 		got := hand.Split()
// 		want := []Hand{
// 			Hand{Card{Queen, Spades}},
// 			Hand{Card{Queen, Spades}},
// 		}

// 		assert.Equal(t, want, got)
// 	})
// }
