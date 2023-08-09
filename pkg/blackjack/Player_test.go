package blackjack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helpers

type theRealMock struct {
	mock.Mock
}

func (strategy *theRealMock) Play(playerHand Hand, dealerHand Hand) PlayerAction {
	args := strategy.Called(playerHand, dealerHand)
	return args.Get(0).(PlayerAction)
}

func (strategy *theRealMock) Bet() int {
	args := strategy.Called()
	return args.Int(0)
}

func (strategy *theRealMock) GetInitialBankroll() float64 {
	args := strategy.Called()
	return args.Get(0).(float64)
}

func (strategy *theRealMock) GetEncodedStrategy() []byte {
	args := strategy.Called()
	return args.Get(0).([]byte)
}

type shoeeeeMock struct {
	mock.Mock
}

func (shoe *shoeeeeMock) Size() int {
	args := shoe.Called()
	return args.Int(0)
}

func (shoe *shoeeeeMock) Peek(count int) []Card {
	args := shoe.Called(count)
	return args.Get(0).([]Card)
}

func (shoe *shoeeeeMock) Shuffle() {}

func (shoe *shoeeeeMock) AdvanceCursor(offset int) (int, error) {
	args := shoe.Called(offset)
	return args.Int(0), args.Error(1)
}

func (shoe *shoeeeeMock) SetPenetration(deckPercentage float64) {}

func (shoe *shoeeeeMock) NeedsReshuffle() bool {
	args := shoe.Called()
	return args.Bool(0)
}

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

func (strategy *strategyMock) GetEncodedStrategy() []byte {
	return []byte("AAA")
}

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
		initialBankroll := 100.0
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(initialBankroll)

		player := NewPlayer(strategy)

		assert.Equal(t, initialBankroll, player.Bankroll)
	})

	t.Run("Should initialize player games", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(100.0)

		player := NewPlayer(strategy)

		assert.Equal(t, 0, len(player.Games))
	})

	t.Run("Should initialize player games played counter", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(100.0)

		player := NewPlayer(strategy)

		assert.Equal(t, 0, player.GamesPlayed)
	})

	t.Run("Should initialize player games won counter", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(100.0)

		player := NewPlayer(strategy)

		assert.Equal(t, 0, player.GamesWon)
	})

	t.Run("Should initialize player games lost counter", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(100.0)

		player := NewPlayer(strategy)

		assert.Equal(t, 0, player.GamesLost)
	})
}

func TestPlayerBet(t *testing.T) {
	strategy := &theRealMock{}
	strategy.On("GetInitialBankroll").Return(100.0)
	strategy.On("Bet").Return(1)

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
	t.Run("Should return the number of cards dealt: regular", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(100.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Hit).Times(5)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(Three, Clubs),
			NewCard(Three, Hearts),
			NewCard(Three, Diamonds),
			NewCard(Three, Spades),
			NewCard(Four, Clubs),
		}, nil)

		player.Bet()

		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}
		cardsTaken := player.Play(playerHand, Hand{}, shoe)

		assert.Equal(t, 5, cardsTaken)
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should return the number of cards dealt: split", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(100.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(SplitOrHit).Once()
		strategy.On("Play", mock.Anything, mock.Anything).Return(Hit).Times(4)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Once()
		strategy.On("Play", mock.Anything, mock.Anything).Return(Hit).Once()
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(Three, Clubs),
			NewCard(Three, Hearts),
			NewCard(Three, Diamonds),
			NewCard(Three, Spades),
			NewCard(Four, Clubs),
			NewCard(Four, Hearts),
			NewCard(Ten, Diamonds),
		}, nil)

		player.Bet()

		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}
		cardsTaken := player.Play(playerHand, Hand{}, shoe)

		assert.Equal(t, 7, cardsTaken)
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should return the number of cards dealt: double", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(100.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Double).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(Three, Clubs),
		}, nil)

		player.Bet()

		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}
		cardsTaken := player.Play(playerHand, Hand{}, shoe)

		assert.Equal(t, 1, cardsTaken)
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should deduct bet from bankroll if split", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(100.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(SplitOrHit).Once()
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Times(2)

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(King, Clubs),
			NewCard(King, Hearts),
		}, nil)

		player.Bet()

		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}
		player.Play(playerHand, Hand{}, shoe)

		assert.Equal(t, 98.0, player.Bankroll)
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should not split if no funds are available", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(SplitOrHit).Once()
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			// TODO: when split or hold is implemented, remove this
			NewCard(King, Clubs),
		}, nil)

		player.Bet()
		playerHand := Hand{NewCard(King, Clubs), NewCard(King, Hearts)}

		player.Play(playerHand, Hand{}, shoe)

		assert.Equal(t, 0.0, player.Bankroll)
		assert.Equal(t, 1, len(player.Games))
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should deduct bet from bankroll if double", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(100.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Double).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(Three, Clubs),
		}, nil)

		player.Bet()
		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}

		player.Play(playerHand, Hand{}, shoe)

		assert.Equal(t, 98.0, player.Bankroll)
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should not double if no funds are available, hit instead", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Double).Once()
		strategy.On("Play", mock.Anything, mock.Anything).Return(Hit).Once()
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(Three, Clubs),
			NewCard(King, Clubs),
		}, nil)

		player.Bet()

		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}
		cardsTaken := player.Play(playerHand, Hand{}, shoe)

		assert.Equal(t, 0.0, player.Bankroll)
		assert.Equal(t, 1, player.Games[0].bet)
		assert.Equal(t, 2, cardsTaken)

		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("No cards should be dealt if dealer has BJ (Ace up)", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1.0)
		strategy.On("Bet").Return(1)

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}

		player.Bet()
		dealerUpcard := NewCard(Ace, Clubs)
		dealerHoleCard := NewCard(Ten, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 0, cardsTaken)
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Normal game if dealer has BJ (Ace in the hole)", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Hit).Once()
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(Three, Clubs),
		}, nil)

		player.Bet()
		dealerUpcard := NewCard(King, Clubs)
		dealerHoleCard := NewCard(Ace, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Three, Clubs), NewCard(Three, Hearts)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 1, cardsTaken)
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should not play anything (nor crash) if no bet is made", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1.0)

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}

		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Nine, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Eight, Clubs), NewCard(Eight, Hearts)}

		shoeIndex := player.Play(playerHand, dealerHand, shoe)

		assert.Equal(t, 0, shoeIndex)
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})
}

func TestPlayerBankrollAfterPlay(t *testing.T) {
	t.Run("Should credit winnings to player bankroll", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}

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
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should reflect player loss after game", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(King, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Five, Clubs), NewCard(Ace, Hearts)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 0, cardsTaken)
		assert.Equal(t, 999.0, player.Bankroll)
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should credit winnings after spliting and winning one game", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(SplitOrHit).Once()
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Times(2)

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(King, Clubs),
			NewCard(Five, Hearts),
		}, nil)

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Seven, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Queen, Clubs), NewCard(Ace, Hearts)}

		player.Play(playerHand, dealerHand, shoe)
		assert.Equal(t, 2, len(player.Games))

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 1000.0, player.Bankroll)
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should credit winnings after spliting", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(SplitOrHit).Once()
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Times(2)

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(King, Clubs),
			NewCard(Jack, Hearts),
		}, nil)

		player.Bet()
		dealerUpcard := NewCard(Seven, Clubs)
		dealerHoleCard := NewCard(Ten, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Queen, Clubs), NewCard(Queen, Hearts)}

		player.Play(playerHand, dealerHand, shoe)
		assert.Equal(t, 2, len(player.Games))

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 1002.0, player.Bankroll)
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should reflect losses after spliting and losing", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(SplitOrHit).Once()
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Times(2)

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(Five, Clubs),
			NewCard(Six, Hearts),
		}, nil)

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
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should credit winnings to player after doubling", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Double).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(Jack, Clubs),
		}, nil)

		player.Bet()
		dealerUpcard := NewCard(Ten, Clubs)
		dealerHoleCard := NewCard(Seven, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Six, Clubs), NewCard(Five, Diamonds)}

		cardsTaken := player.Play(playerHand, dealerHand, shoe)

		dealerHand = Hand{dealerUpcard, dealerHoleCard, NewCard(King, Clubs)}
		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 1, cardsTaken)
		assert.Equal(t, 1002.0, player.Bankroll)
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should reflect loss after doubling", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Double).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(Three, Clubs),
		}, nil)

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
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should credit initial bet when pushing", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}

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
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should clear games after resolving", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(SplitOrHit).Once()
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Times(2)

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(Five, Clubs),
			NewCard(Six, Hearts),
		}, nil)

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
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should increase the games seen counter", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)

		player := NewPlayer(strategy)

		player.Bet()

		assert.Equal(t, 1, player.GamesSeen)
		strategy.AssertExpectations(t)
	})

	t.Run("Should increase the games played counter", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}

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
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should increase the gamesWon counter", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}

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
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should increase the gamesLost counter", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}

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
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should increase the gamesPushed counter", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Once()

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}

		player.Bet()
		dealerUpcard := NewCard(King, Clubs)
		dealerHoleCard := NewCard(King, Hearts)
		dealerHoleCard.SetHole()

		dealerHand := Hand{dealerUpcard, dealerHoleCard}
		playerHand := Hand{NewCard(Queen, Clubs), NewCard(Queen, Hearts)}

		player.Play(playerHand, dealerHand, shoe)

		dealerHand.Reveal()
		player.Resolve(dealerHand)

		assert.Equal(t, 1, player.GamesPushed)
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})

	t.Run("Should increase counters by 1 when splitting", func(t *testing.T) {
		strategy := &theRealMock{}
		strategy.On("GetInitialBankroll").Return(1000.0)
		strategy.On("Bet").Return(1)
		strategy.On("Play", mock.Anything, mock.Anything).Return(SplitOrHit).Once()
		strategy.On("Play", mock.Anything, mock.Anything).Return(Stand).Times(2)

		player := NewPlayer(strategy)

		shoe := &shoeeeeMock{}
		shoe.On("Peek", mock.Anything).Return([]Card{
			NewCard(Ten, Clubs),
			NewCard(Queen, Hearts),
		}, nil)

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
		strategy.AssertExpectations(t)
		shoe.AssertExpectations(t)
	})
}

func TestPlayerGetStatistics(t *testing.T) {
	strategy := &strategyMock{}
	strategy.initialBankroll = 100
	strategy.alwaysHit = true

	player := NewPlayer(strategy)
	player.GamesSeen = 11
	player.GamesPlayed = 10
	player.GamesWon = 1
	player.GamesLost = 7

	t.Run("Should return the games played", func(t *testing.T) {
		got := player.GetStatistics().GamesPlayed
		want := 10

		assert.Equal(t, want, got)
	})

	t.Run("Should return the games seen", func(t *testing.T) {
		got := player.GetStatistics().GamesSeen
		want := 11

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

	t.Run("Should return the strategy", func(t *testing.T) {
		got := player.GetStatistics().Strategy
		want := strategy.GetEncodedStrategy()

		assert.Equal(t, want, got)
	})
}

// func TestGameRestrictions(t *testing.T) {

// }
