package blackjack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helpers

type strategyMock struct {
	initialBankroll int
	alwaysHit       bool
	doubleThenHit   bool
	splitThenHit    bool
}

func (strategy *strategyMock) Play(playerHand Hand, dealerHand Hand) PlayerAction {
	_, playerIsBusted := playerHand.Score()
	if playerIsBusted {
		return Stand
	}

	if strategy.splitThenHit {
		strategy.splitThenHit = false
		strategy.alwaysHit = true
		return SplitOrHit
	}

	if strategy.doubleThenHit {
		strategy.doubleThenHit = false
		strategy.alwaysHit = true
		return Double
	}

	if strategy.alwaysHit {
		return Hit
	}

	return Stand
}

func (strategy *strategyMock) Bet() int {
	return 1
}

func (strategy *strategyMock) GetInitialBankroll() float64 {
	return float64(strategy.initialBankroll)
}

// Not used in these tests
func (strategy *strategyMock) GetEncodedStrategy() []byte { return []byte{} }

type shoeMock struct{}

func (shoe *shoeMock) Size() int {
	return 52
}

func (shoe *shoeMock) Peek(count int) []Card {
	peek := []Card{
		NewCard(Three, Clubs),
		NewCard(Three, Hearts),
		NewCard(Three, Diamonds),
		NewCard(Three, Spades),
		NewCard(Four, Clubs),
		NewCard(Four, Hearts),
		NewCard(Ten, Diamonds),
		NewCard(King, Spades),
		NewCard(Ace, Clubs),
	}
	return peek
}

// Not used in these tests
func (shoe *shoeMock) Shuffle()                              {}
func (shoe *shoeMock) AdvanceCursor(offset int) (int, error) { return 0, nil }
func (shoe *shoeMock) SetPenetration(deckPercentage float64) {}
func (shoe *shoeMock) NeedsReshuffle() bool                  { return false }

// Tests

func TestPlayerInitialization(t *testing.T) {
	t.Run("Should set player initial bankroll", func(t *testing.T) {
		strategy := &strategyMock{}
		strategy.initialBankroll = 100
		player := NewPlayer(strategy)

		assert.Equal(t, float64(strategy.initialBankroll), player.Bankroll)
	})

	t.Run("Should initialize player games", func(t *testing.T) {
		strategy := &strategyMock{}
		strategy.initialBankroll = 100
		player := NewPlayer(strategy)

		assert.Equal(t, 0, len(player.Games))
	})

	t.Run("Should initialize player games played counter", func(t *testing.T) {
		strategy := &strategyMock{}
		strategy.initialBankroll = 100
		player := NewPlayer(strategy)

		assert.Equal(t, 0, player.GamesPlayed)
	})

	t.Run("Should initialize player games won counter", func(t *testing.T) {
		strategy := &strategyMock{}
		strategy.initialBankroll = 100
		player := NewPlayer(strategy)

		assert.Equal(t, 0, player.GamesWon)
	})

	t.Run("Should initialize player games lost counter", func(t *testing.T) {
		strategy := &strategyMock{}
		strategy.initialBankroll = 100
		player := NewPlayer(strategy)

		assert.Equal(t, 0, player.GamesLost)
	})
}

func TestPlayerBet(t *testing.T) {
	strategy := &strategyMock{}
	strategy.initialBankroll = 100
	strategy.alwaysHit = true

	t.Run("Should decide to play a hand if has funds", func(t *testing.T) {
		player := NewPlayer(strategy)

		willBet, ammount := player.Bet( /* shoe? */ )
		assert.True(t, willBet)
		assert.Equal(t, 1, ammount)
	})

	t.Run("Should decide not to play a hand if has no funds", func(t *testing.T) {
		player := NewPlayer(strategy)
		player.Bankroll = 0

		willBet, _ := player.Bet( /* shoe? */ )
		assert.False(t, willBet)
	})

	t.Run("Bet should be deducted from bankroll", func(t *testing.T) {
		player := NewPlayer(strategy)
		player.Bankroll = 10

		willBet, _ := player.Bet( /* shoe? */ )
		assert.True(t, willBet)
		assert.Equal(t, 9.0, player.Bankroll)
	})

	t.Run("Should create a new game", func(t *testing.T) {
		player := NewPlayer(strategy)
		player.Bankroll = 10

		player.Bet( /* shoe? */ )
		assert.Equal(t, 1, len(player.Games))
	})

	t.Run("Should create a new game with the bet ammount", func(t *testing.T) {
		player := NewPlayer(strategy)
		player.Bankroll = 10

		player.Bet( /* shoe? */ )
		assert.Equal(t, 1, player.Games[0].bet)
	})
}

func TestPlayerPlay(t *testing.T) {
	initTest := func(bankroll int, alwaysHit bool, doubleThenHit bool, splitThenHit bool) (*Player, *shoeMock) {
		strategy := &strategyMock{}
		strategy.initialBankroll = bankroll
		strategy.alwaysHit = alwaysHit
		strategy.doubleThenHit = doubleThenHit
		strategy.splitThenHit = splitThenHit

		player := NewPlayer(strategy)
		shoe := &shoeMock{}

		return player, shoe
	}

	t.Run("Should return the number of cards dealt: regular", func(t *testing.T) {
		player, shoe := initTest(100, true, false, false)

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Ten, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 5, cardsTaken)
	})

	t.Run("Should return the number of cards dealt: split", func(t *testing.T) {
		player, shoe := initTest(100, false, false, true)

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Ten, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 8, cardsTaken)
	})

	t.Run("Should return the number of cards dealt: double", func(t *testing.T) {
		player, shoe := initTest(100, false, true, false)

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Ten, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 1, cardsTaken)
	})

	t.Run("Should deduct bet from bankroll if split", func(t *testing.T) {
		player, shoe := initTest(100, false, false, true)

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Ten, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}

		player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 98.0, player.Bankroll)
	})

	t.Run("Should not split if no funds are available", func(t *testing.T) {
		player, shoe := initTest(1, false, false, true)

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Ten, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}

		player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 0.0, player.Bankroll)
		assert.Equal(t, 1, len(player.Games))
	})

	t.Run("Should deduct bet from bankroll if double", func(t *testing.T) {
		player, shoe := initTest(100, false, true, false)

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Ten, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}

		player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 98.0, player.Bankroll)
	})

	t.Run("Should not double if no funds are available", func(t *testing.T) {
		player, shoe := initTest(1, false, true, false)

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Ten, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}

		player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 0.0, player.Bankroll)
		assert.Equal(t, 1, player.Games[0].bet)
	})

	t.Run("No cards should be dealt if dealer has BJ (Ace up)", func(t *testing.T) {
		player, shoe := initTest(100, true, false, false)

		player.Bet()
		dealerUpcard := NewCard(Ace, Clubs)
		dealerHoleCard := NewCard(Ten, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 0, cardsTaken)
	})

	t.Run("Normal game if dealer has BJ (Ace in the hole)", func(t *testing.T) {
		player, shoe := initTest(100, true, false, false)

		player.Bet()
		dealerUpcard := NewCard(King, Clubs)
		dealerHoleCard := NewCard(Ace, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 5, cardsTaken)
	})

	t.Run("Should not play anything (nor crash) if no bet is made", func(t *testing.T) {
		player, shoe := initTest(100, false, false, true)

		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Nine, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Eight, Clubs), NewCard(Eight, Hearts)}

		shoeIndex := player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 0, shoeIndex)
	})
}

func TestPlayerBankrollAfterPlay(t *testing.T) {
	initTest := func() (*Player, *shoeMock) {
		// TODO: refactor getTestStrategy into its own file
		// because it's imported from another test file
		raw, err := getTestStrategy()
		if err != nil {
			t.Fatal(err)
		}
		strategy, _ := NewStrategy(raw)

		player := NewPlayer(strategy)
		shoe := &shoeMock{}

		return player, shoe
	}

	t.Run("Should credit winnings to player bankroll", func(t *testing.T) {
		player, shoe := initTest()

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Nine, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(King, Clubs), NewCard(Queen, Hearts)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 0, cardsTaken)
		assert.Equal(t, 1001.0, player.Bankroll)
	})

	t.Run("Should reflect player loss after game", func(t *testing.T) {
		player, shoe := initTest()

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(King, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Five, Clubs), NewCard(Ace, Hearts)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 1, cardsTaken)
		assert.Equal(t, 999.0, player.Bankroll)
	})

	t.Run("Should credit winnings after spliting and winning one game", func(t *testing.T) {
		player, shoe := initTest()

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Seven, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Ace, Clubs), NewCard(Ace, Hearts)}

		player.Play(playerHand, dealerHand, shoe)
		assert.Equal(t, 2, len(player.Games))

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 1000.0, player.Bankroll)
	})

	t.Run("Should credit winnings after spliting", func(t *testing.T) {
		player, shoe := initTest()

		player.Bet()
		dealerUpcard := NewCard(Seven, Clubs)
		dealerHoleCard := NewCard(Ten, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Ace, Clubs), NewCard(Ace, Hearts)}

		player.Play(playerHand, dealerHand, shoe)
		assert.Equal(t, 2, len(player.Games))

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 1002.0, player.Bankroll)
	})

	t.Run("Should credit winnings after spliting and losing", func(t *testing.T) {
		player, shoe := initTest()

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Ace, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Ace, Clubs), NewCard(Ace, Hearts)}

		player.Play(playerHand, dealerHand, shoe)
		assert.Equal(t, 2, len(player.Games))

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 998.0, player.Bankroll)
	})

	t.Run("Should credit winnings to player after doubling", func(t *testing.T) {
		player, shoe := initTest()

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Five, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Six, Clubs), NewCard(Five, Diamonds)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		dealerHand = Hand{dealerUpcard, dealerHoleCard, NewCard(King, Clubs)}
		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 1, cardsTaken)
		assert.Equal(t, 1002.0, player.Bankroll)
	})

	t.Run("Should reflect loss after doubling", func(t *testing.T) {
		player, shoe := initTest()

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Queen, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Six, Clubs), NewCard(Five, Diamonds)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 1, cardsTaken)
		assert.Equal(t, 998.0, player.Bankroll)
	})

	t.Run("Should credit initial bet when pushing", func(t *testing.T) {
		player, shoe := initTest()

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Nine, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(King, Clubs), NewCard(Nine, Hearts)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 0, cardsTaken)
		assert.Equal(t, 1000.0, player.Bankroll)
	})

	t.Run("Should clear games after resolving", func(t *testing.T) {
		player, shoe := initTest()

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Nine, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Eight, Clubs), NewCard(Eight, Hearts)}

		player.Play(playerHand, dealerHand, shoe)

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 0, len(player.Games))
	})

	t.Run("Should increase the games played counter", func(t *testing.T) {
		player, shoe := initTest()

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Nine, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(King, Clubs), NewCard(Queen, Hearts)}

		player.Play(playerHand, dealerHand, shoe)

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 1, player.GamesPlayed)
	})

	t.Run("Should increase the gamesWon counter", func(t *testing.T) {
		player, shoe := initTest()

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Nine, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(King, Clubs), NewCard(Queen, Hearts)}

		player.Play(playerHand, dealerHand, shoe)

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 1, player.GamesWon)
	})

	t.Run("Should increase the gamesLost counter", func(t *testing.T) {
		player, shoe := initTest()

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Nine, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(King, Clubs), NewCard(Seven, Hearts)}

		player.Play(playerHand, dealerHand, shoe)

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 1, player.GamesLost)
	})

	t.Run("Should increase counters by 1 when splitting", func(t *testing.T) {
		player, shoe := initTest()

		player.Bet()
		dealerUpcard := NewCard(Seven, Clubs)
		dealerHoleCard := NewCard(Ten, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Ace, Clubs), NewCard(Ace, Hearts)}

		player.Play(playerHand, dealerHand, shoe)
		assert.Equal(t, 2, len(player.Games))

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 1, player.GamesPlayed)
		assert.Equal(t, 1, player.GamesWon)
	})
}

func TestPlayerGetStatistics(t *testing.T) {
	strategy := &strategyMock{}
	strategy.initialBankroll = 100
	strategy.alwaysHit = true

	player := NewPlayer(strategy)
	player.GamesPlayed = 10
	player.GamesWon = 1
	player.GamesLost = 7

	t.Run("Should return the games played", func(t *testing.T) {
		got := player.GetStatistics().GamesPlayed
		want := 10

		assert.Equal(t, want, got)
	})

	t.Run("Should return the games won", func(t *testing.T) {
		got := player.GetStatistics().GamesWon
		want := 1

		assert.Equal(t, want, got)
	})

	t.Run("Should return the games lost", func(t *testing.T) {
		got := player.GetStatistics().GamesLost
		want := 7

		assert.Equal(t, want, got)
	})

	t.Run("Should return the games pushed", func(t *testing.T) {
		got := player.GetStatistics().GamesPushed
		want := 2

		assert.Equal(t, want, got)
	})

	t.Run("Should return the current bankroll", func(t *testing.T) {
		got := player.GetStatistics().Bankroll
		want := 100.0

		assert.Equal(t, want, got)
	})

	t.Run("Should return bankroll delta", func(t *testing.T) {
		player.Bankroll = 200
		got := player.GetStatistics().BankrollDelta
		want := 100.0

		assert.Equal(t, want, got)
	})

	t.Run("Should return the initial bankroll", func(t *testing.T) {
		got := player.GetStatistics().InitialBankroll
		want := 100.0

		assert.Equal(t, want, got)
	})
}
