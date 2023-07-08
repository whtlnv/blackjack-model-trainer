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
