package blackjack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helpers

type playerSpy struct {
	mock.Mock
}

func (p *playerSpy) Bet() (bool, int) {
	args := p.Called()
	return args.Bool(0), args.Int(1)
}

func (p *playerSpy) Play(hand Hand, dealerHand Hand, shoe Shoeish) int {
	args := p.Called(hand, dealerHand, shoe)
	return args.Int(0)
}

func (p *playerSpy) Resolve(dealerHand Hand) {
	p.Called(dealerHand)
}

func (p *playerSpy) GetStatistics() PlayerStatistics {
	args := p.Called()
	return args.Get(0).(PlayerStatistics)
}

// Tests

func TestTableInitialization(t *testing.T) {
	spy := &playerSpy{}
	players := []Playerish{spy}
	shoe := NewShoe(1)
	table := NewTable(players, shoe)

	t.Run("Should initialize a table with the given players", func(t *testing.T) {
		got := table.Players
		want := players
		assert.Equal(t, want, got)
	})

	t.Run("Should initialize a table with the given shoe", func(t *testing.T) {
		got := table.Shoe
		want := shoe
		assert.Equal(t, want, got)
	})
}

func TestDealingHands(t *testing.T) {
	spy := &playerSpy{}
	players := []Playerish{spy}

	shoe := NewShoe(1)
	shoe.SetPenetration(0.5)

	table := NewTable(players, shoe)

	t.Run("Should advance the shoe cursor", func(t *testing.T) {
		table.dealHands()

		got := shoe.cursor
		want := 4

		assert.Equal(t, want, got)
	})

	t.Run("Should deal a player hand", func(t *testing.T) {
		topCards := shoe.Peek(4)
		playerHand, _ := table.dealHands()

		got := playerHand
		want := Hand{topCards[0], topCards[2]}

		assert.Equal(t, want, got)
	})

	t.Run("Should deal a dealer hand", func(t *testing.T) {
		topCards := shoe.Peek(4)
		_, dealerHand := table.dealHands()

		got := dealerHand
		want := Hand{topCards[1], topCards[3]}
		want[1].SetHole()

		assert.Equal(t, want, got)
	})

	t.Run("Should shuffle the shoe if the penetration index is reached", func(t *testing.T) {
		shoe.AdvanceCursor(27)
		table.dealHands()

		got := shoe.cursor
		want := 4

		assert.Equal(t, want, got)
	})
}

func TestPlayerGames(t *testing.T) {
	spy1 := &playerSpy{}
	spy1.On("Play", mock.AnythingOfType("Hand"), mock.AnythingOfType("Hand"), mock.Anything).Return(0)

	spy2 := &playerSpy{}
	spy2.On("Play", mock.AnythingOfType("Hand"), mock.AnythingOfType("Hand"), mock.Anything).Return(5)

	spy3 := &playerSpy{}
	spy3.On("Play", mock.AnythingOfType("Hand"), mock.AnythingOfType("Hand"), mock.Anything).Return(3)

	players := []Playerish{spy1, spy2, spy3}

	shoe := NewShoe(1)
	shoe.SetPenetration(0.5)

	table := NewTable(players, shoe)

	playerHand := Hand{}
	dealerHand := Hand{}

	table.playAllHands(playerHand, dealerHand)

	t.Run("Should play each player's hand", func(t *testing.T) {
		spy1.AssertNumberOfCalls(t, "Play", 1)
		spy2.AssertNumberOfCalls(t, "Play", 1)
		spy3.AssertNumberOfCalls(t, "Play", 1)
	})

	t.Run("Should advance the shoe cursor", func(t *testing.T) {
		want := 5
		got := shoe.cursor

		assert.Equal(t, want, got)
	})
}

func TestDealerGame(t *testing.T) {
	spy := &playerSpy{}
	players := []Playerish{spy}

	t.Run("Should stand on soft 17", func(t *testing.T) {
		shoe := NewShoe(1)
		shoe.cards = []Card{
			NewCard(Three, Hearts),
			NewCard(King, Clubs),
		}
		table := NewTable(players, shoe)

		dealerHand := Hand{NewCard(Ace, Clubs), NewCard(Three, Clubs)}
		got := table.playDealerHand(dealerHand)

		want := Hand{
			NewCard(Ace, Clubs),
			NewCard(Three, Clubs),
			NewCard(Three, Hearts),
		}

		assert.Equal(t, want, got)
	})

	t.Run("Should stand on hard 17", func(t *testing.T) {
		shoe := NewShoe(1)
		shoe.cards = []Card{
			NewCard(Three, Hearts),
			NewCard(King, Clubs),
		}
		table := NewTable(players, shoe)

		dealerHand := Hand{NewCard(Ten, Clubs), NewCard(Seven, Clubs)}
		got := table.playDealerHand(dealerHand)

		want := Hand{
			NewCard(Ten, Clubs),
			NewCard(Seven, Clubs),
		}

		assert.Equal(t, want, got)
	})

	t.Run("Should stand on bust", func(t *testing.T) {
		shoe := NewShoe(1)
		shoe.cards = []Card{
			NewCard(Three, Hearts),
			NewCard(King, Clubs),
		}
		table := NewTable(players, shoe)

		dealerHand := Hand{NewCard(King, Diamonds), NewCard(Two, Spades)}
		got := table.playDealerHand(dealerHand)

		want := Hand{
			NewCard(King, Diamonds),
			NewCard(Two, Spades),
			NewCard(Three, Hearts),
			NewCard(King, Clubs),
		}

		assert.Equal(t, want, got)
	})

	t.Run("Should advance the shoe cursor", func(t *testing.T) {
		shoe := NewShoe(1)
		shoe.cards = []Card{
			NewCard(Three, Hearts),
			NewCard(Four, Clubs),
			NewCard(Five, Diamonds),
			NewCard(Six, Spades),
		}
		table := NewTable(players, shoe)

		dealerHand := Hand{NewCard(Two, Diamonds), NewCard(Seven, Spades)}
		table.playDealerHand(dealerHand)

		want := 3
		got := shoe.cursor

		assert.Equal(t, want, got)
	})

	t.Run("Should reveal the dealer hand", func(t *testing.T) {
		shoe := NewShoe(1)
		shoe.cards = []Card{
			NewCard(Three, Hearts),
			NewCard(Four, Clubs),
		}
		table := NewTable(players, shoe)

		dealerHand := Hand{NewCard(King, Diamonds), NewCard(Seven, Spades)}
		dealerHand[1].SetHole()

		result := table.playDealerHand(dealerHand)
		got := result[1].hole

		assert.False(t, got)
	})
}

func TestTableRun(t *testing.T) {
	numberOfPlayers := 3
	dealtCards := 4
	playerCards := 5
	dealerCards := 1

	asPlayers := []Playerish{}
	asSpies := []*playerSpy{}
	for i := 0; i < numberOfPlayers; i++ {
		spy := &playerSpy{}
		spy.On("Bet").Return(true, 10)
		spy.On("Play", mock.AnythingOfType("Hand"), mock.AnythingOfType("Hand"), mock.Anything).Return(playerCards)
		spy.On("Resolve", mock.AnythingOfType("Hand")).Return()

		asPlayers = append(asPlayers, spy)
		asSpies = append(asSpies, spy)
	}

	shoe := NewShoe(1)
	shoe.cards = []Card{
		NewCard(Three, Hearts),
		NewCard(King, Diamonds),
		NewCard(Four, Clubs),
		NewCard(Five, Diamonds),
		NewCard(Two, Spades),
		NewCard(Three, Hearts),
		NewCard(King, Diamonds),
		NewCard(Four, Clubs),
		NewCard(Five, Diamonds),
		NewCard(Two, Spades),
	}
	shoe.SetPenetration(0.5)

	table := NewTable(asPlayers, shoe)
	table.Run()

	t.Run("Should call Bet() on each player once", func(t *testing.T) {
		for _, spy := range asSpies {
			spy.AssertNumberOfCalls(t, "Bet", 1)
		}
	})

	t.Run("Should call Play() on each player once", func(t *testing.T) {
		for _, spy := range asSpies {
			spy.AssertNumberOfCalls(t, "Play", 1)
		}
	})

	t.Run("Should advance the shoe cursor by the cards used", func(t *testing.T) {
		got := shoe.cursor
		want := dealtCards + playerCards + dealerCards

		assert.Equal(t, want, got)
	})

	t.Run("Should call Resolve() on each player once", func(t *testing.T) {
		for _, spy := range asSpies {
			spy.AssertNumberOfCalls(t, "Resolve", 1)
		}
	})
}

func TestTableRunMany(t *testing.T) {
	shoe := NewShoe(1)
	shoe.SetPenetration(0.5)

	// TODO: refactor getTestStrategy into its own file
	// because it's imported from another test file
	raw, err := getTestStrategy()
	if err != nil {
		t.Fatal(err)
	}
	strategy, _ := NewStrategy(raw)

	makePlayers := func(numberOfPlayers int, bankroll float64) []Playerish {
		players := []Playerish{}
		for i := 0; i < numberOfPlayers; i++ {
			aPlayer := NewPlayer(strategy)
			aPlayer.Bankroll = bankroll
			players = append(players, aPlayer)
		}

		return players
	}

	t.Run("Should run until the given number of hands is reached", func(t *testing.T) {
		runs := 10
		players := makePlayers(1, float64(runs*2))
		table := NewTable(players, shoe)

		table.RunMany(runs)

		got := players[0].GetStatistics().GamesPlayed
		want := runs

		assert.Equal(t, want, got)
	})

	t.Run("Should stop running if all players run out of money", func(t *testing.T) {
		runs := 100
		players := makePlayers(1, 0)
		table := NewTable(players, shoe)

		table.RunMany(runs)

		got := players[0].GetStatistics().GamesPlayed
		want := 0

		assert.Equal(t, want, got)
	})
}
