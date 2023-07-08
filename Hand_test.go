package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		want := 19

		assert.Equal(t, want, got)
		assert.Equal(t, false, isBusted)
	})

	t.Run("Should return the correct score for Ace", func(t *testing.T) {
		hand := Hand{
			Card{Ace, Spades},
			Card{King, Hearts},
			Card{Queen, Diamonds},
		}

		got, isBusted := hand.Score()
		want := 21

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
		want := 15

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
		want := 30

		assert.Equal(t, want, got)
		assert.Equal(t, true, isBusted)
	})
}
