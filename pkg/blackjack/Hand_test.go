package blackjack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDealingACard(t *testing.T) {
	t.Run("Should deal a card to a hand", func(t *testing.T) {
		hand := Hand{}
		card := NewCard(Queen, Spades)

		hand.Deal(card)

		assert.Equal(t, 1, len(hand))
		assert.Equal(t, card, hand[0])
	})

	t.Run("Should append a card to a hand", func(t *testing.T) {
		hand := Hand{NewCard(King, Spades)}
		card := NewCard(Queen, Spades)

		hand.Deal(card)

		assert.Equal(t, 2, len(hand))
		assert.Equal(t, card, hand[1])
	})
}

func TestHandValues(t *testing.T) {
	t.Run("Should return the value of a hand", func(t *testing.T) {
		hand := Hand{
			NewCard(Two, Spades),
			NewCard(Seven, Hearts),
			NewCard(Queen, Diamonds),
		}

		got := hand.Values()
		want := []int{19}

		assert.ElementsMatch(t, want, got)
	})

	t.Run("Should return the correct values for an Ace", func(t *testing.T) {
		hand := Hand{
			NewCard(Ace, Spades),
			NewCard(King, Hearts),
			NewCard(Queen, Diamonds),
		}

		got := hand.Values()
		want := []int{21, 31}

		assert.Equal(t, want, got)
	})

	t.Run("Should return the correct values for several Aces", func(t *testing.T) {
		hand := Hand{
			NewCard(Ace, Spades),
			NewCard(Ace, Spades),
			NewCard(Three, Diamonds),
		}

		got := hand.Values()
		want := []int{5, 15, 25}

		assert.Equal(t, want, got)
	})

	t.Run("Should return the correct value for a hole card", func(t *testing.T) {
		hand := Hand{
			NewCard(Nine, Spades),
			NewCard(King, Clubs),
		}
		hand[1].SetHole()

		got := hand.Values()
		want := []int{9}

		assert.Equal(t, want, got)
	})
}

func TestHandScore(t *testing.T) {
	t.Run("Should return the highest score of a hand", func(t *testing.T) {
		hand := Hand{
			NewCard(Two, Spades),
			NewCard(Seven, Hearts),
			NewCard(Queen, Diamonds),
		}

		got, isBusted := hand.Score()
		want := HandScore{19, 19}

		assert.Equal(t, want, got)
		assert.Equal(t, false, isBusted)
	})

	t.Run("Should return the correct score for Ace", func(t *testing.T) {
		hand := Hand{
			NewCard(Ace, Spades),
			NewCard(King, Hearts),
		}

		got, isBusted := hand.Score()
		want := HandScore{11, 21}

		assert.Equal(t, want, got)
		assert.Equal(t, false, isBusted)
	})

	t.Run("Should return the correct score for several Aces", func(t *testing.T) {
		hand := Hand{
			NewCard(Ace, Spades),
			NewCard(Ace, Spades),
			NewCard(Three, Diamonds),
		}

		got, isBusted := hand.Score()
		want := HandScore{5, 15}

		assert.Equal(t, want, got)
		assert.Equal(t, false, isBusted)
	})

	t.Run("Should return true if busted", func(t *testing.T) {
		hand := Hand{
			NewCard(Queen, Spades),
			NewCard(Queen, Spades),
			NewCard(Queen, Diamonds),
		}

		got, isBusted := hand.Score()
		want := HandScore{30, 30}

		assert.Equal(t, want, got)
		assert.Equal(t, true, isBusted)
	})

	t.Run("Should return the correct score for a hole card", func(t *testing.T) {
		hand := Hand{
			NewCard(Ace, Spades),
			NewCard(King, Clubs),
		}
		hand[1].SetHole()

		got, isBusted := hand.Score()
		want := HandScore{1, 11}

		assert.Equal(t, want, got)
		assert.Equal(t, false, isBusted)
	})
}

func TestIsPair(t *testing.T) {
	t.Run("Should return true if pair", func(t *testing.T) {
		hand := Hand{
			NewCard(Queen, Spades),
			NewCard(Queen, Spades),
		}

		got := hand.IsPair()
		want := true

		assert.Equal(t, want, got)
	})

	t.Run("Should return false if not pair", func(t *testing.T) {
		hand := Hand{
			NewCard(Queen, Spades),
			NewCard(King, Spades),
		}

		got := hand.IsPair()
		want := false

		assert.Equal(t, want, got)
	})
}

func TestHasSoftValue(t *testing.T) {
	t.Run("Should return true if soft value", func(t *testing.T) {
		hand := Hand{
			NewCard(Ace, Spades),
			NewCard(King, Spades),
		}

		got := hand.HasSoftValue()
		want := true

		assert.Equal(t, want, got)
	})

	t.Run("Should return false if not soft value", func(t *testing.T) {
		hand := Hand{
			NewCard(Queen, Spades),
			NewCard(King, Spades),
		}

		got := hand.HasSoftValue()
		want := false

		assert.Equal(t, want, got)
	})
}

func TestIsBlackjack(t *testing.T) {
	t.Run("Should return true if blackjack", func(t *testing.T) {
		hand := Hand{
			NewCard(Ace, Spades),
			NewCard(King, Spades),
		}
		hand[1].SetHole()

		got := hand.IsBlackjack()
		want := true

		assert.Equal(t, want, got)
	})

	t.Run("Should return false if not blackjack", func(t *testing.T) {
		hand := Hand{
			NewCard(Queen, Spades),
			NewCard(King, Spades),
		}

		got := hand.IsBlackjack()
		want := false

		assert.Equal(t, want, got)
	})

	t.Run("Input hand should no be modified", func(t *testing.T) {
		hand := Hand{
			NewCard(Ace, Spades),
			NewCard(King, Spades),
		}
		hand[1].SetHole()

		hand.IsBlackjack()

		assert.Equal(t, true, hand[1].hole)
	})
}

func TestGetHoleCard(t *testing.T) {
	t.Run("Should return the hole card", func(t *testing.T) {
		hand := Hand{
			NewCard(Ace, Spades),
			NewCard(King, Spades),
		}
		hand[1].SetHole()

		got := hand.GetHoleCard()
		want := &hand[1]

		assert.Equal(t, want, got)
	})

	t.Run("Should return nil if no hole card", func(t *testing.T) {
		hand := Hand{
			NewCard(Ace, Spades),
			NewCard(King, Spades),
		}

		got := hand.GetHoleCard()
		want := &Card{}

		assert.Equal(t, want, got)
	})
}
