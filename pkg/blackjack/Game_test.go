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
		got := game.IsDoubled

		assert.Equal(t, want, got)
	})
}

func TestGameSplit(t *testing.T) {
	t.Run("Should remove a card from the split game hand", func(t *testing.T) {
		game := NewGame(1)
		hand := Hand{NewCard(Two, Spades), NewCard(Two, Hearts)}
		game.SetHand(hand)

		game.Split()

		want := Hand{NewCard(Two, Spades)}
		got := game.hand

		assert.Equal(t, want, got)
	})

	t.Run("Should flag a hand that was split", func(t *testing.T) {
		game := NewGame(1)
		hand := Hand{NewCard(Two, Spades), NewCard(Two, Hearts)}
		game.SetHand(hand)

		game.Split()

		want := true
		got := game.IsSplit

		assert.Equal(t, want, got)
	})

	t.Run("Should return a game with the split card", func(t *testing.T) {
		game := NewGame(1)
		hand := Hand{NewCard(Two, Spades), NewCard(Two, Hearts)}
		game.SetHand(hand)

		splitGame := game.Split()

		want := Hand{NewCard(Two, Hearts)}
		got := splitGame.hand

		assert.Equal(t, want, got)
	})

	t.Run("Should return a game with the same bet", func(t *testing.T) {
		bet := 1
		game := NewGame(bet)
		hand := Hand{NewCard(Two, Spades), NewCard(Two, Hearts)}
		game.SetHand(hand)

		splitGame := game.Split()

		want := bet
		got := splitGame.bet

		assert.Equal(t, want, got)
	})

	t.Run("Should return a game flagged as split", func(t *testing.T) {
		game := NewGame(1)
		hand := Hand{NewCard(Two, Spades), NewCard(Two, Hearts)}
		game.SetHand(hand)

		splitGame := game.Split()

		want := true
		got := splitGame.IsSplit

		assert.Equal(t, want, got)
	})
}

func TestGameResolution(t *testing.T) {
	testCases := []struct {
		desc          string
		playerHand    *Hand
		dealerHand    *Hand
		inputBet      int
		expectedValue float64
	}{
		{
			desc:          "Player win -> return bet + winnings",
			playerHand:    &Hand{NewCard(King, Diamonds), NewCard(Queen, Spades)},
			dealerHand:    &Hand{NewCard(Jack, Hearts), NewCard(Nine, Clubs)},
			inputBet:      10,
			expectedValue: 20.0,
		},
		{
			desc:          "Player bust -> 0",
			playerHand:    &Hand{NewCard(King, Diamonds), NewCard(Queen, Spades), NewCard(Two, Clubs)},
			dealerHand:    &Hand{NewCard(Jack, Hearts), NewCard(Nine, Clubs), NewCard(Queen, Hearts)},
			inputBet:      10,
			expectedValue: 0.0,
		},
		{
			desc:          "Dealer bust -> return bet + winnings",
			playerHand:    &Hand{NewCard(King, Diamonds), NewCard(Queen, Spades)},
			dealerHand:    &Hand{NewCard(Jack, Hearts), NewCard(Nine, Clubs), NewCard(Three, Hearts)},
			inputBet:      10,
			expectedValue: 20.0,
		},
		{
			desc:          "Player lose -> 0",
			playerHand:    &Hand{NewCard(Ace, Diamonds), NewCard(Five, Spades)},
			dealerHand:    &Hand{NewCard(Jack, Hearts), NewCard(Two, Clubs), NewCard(Five, Hearts)},
			inputBet:      10,
			expectedValue: 0.0,
		},
		{
			desc:          "Player and dealer push",
			playerHand:    &Hand{NewCard(Ace, Diamonds), NewCard(Two, Spades), NewCard(Four, Clubs)},
			dealerHand:    &Hand{NewCard(Seven, Hearts), NewCard(Ten, Clubs)},
			inputBet:      10,
			expectedValue: 10.0,
		},
		{
			desc:          "Player blackjack -> return bet + winnings",
			playerHand:    &Hand{NewCard(Ace, Diamonds), NewCard(Ten, Spades)},
			dealerHand:    &Hand{NewCard(Seven, Hearts), NewCard(Ten, Clubs)},
			inputBet:      1,
			expectedValue: 2.5,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			game := NewGame(testCase.inputBet)
			game.SetHand(*testCase.playerHand)

			want := testCase.expectedValue
			got := game.Resolve(testCase.dealerHand)

			assert.Equal(t, want, got)
		})
	}
}
