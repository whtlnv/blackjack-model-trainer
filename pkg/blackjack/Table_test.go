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

	table.playAllGames(playerHand, dealerHand)

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

// func TestDealerGame(t *testing.T) {
// 	spy := &playerSpy{}
// 	players := []Playerish{spy}

// 	shoe := NewShoe(1)
// 	shoe.SetPenetration(0.5)

// 	table := NewTable(players, shoe)

// 	t.Run("Should stand on soft 17 or greater", func(t *testing.T) {})

// 	t.Run("Should stand on bust", func(t *testing.T) {})

// 	t.Run("Should advance the shoe cursor", func(t *testing.T) {})
// }

func TestTableRun(t *testing.T) {
	numberOfPlayers := 3
	cardsUsed := 5
	asPlayers := []Playerish{}
	asSpies := []*playerSpy{}
	for i := 0; i < numberOfPlayers; i++ {
		spy := &playerSpy{}
		spy.On("Bet").Return(true, 10)
		spy.On("Play", mock.AnythingOfType("Hand"), mock.AnythingOfType("Hand"), mock.Anything).Return(cardsUsed)
		spy.On("Resolve").Return()

		asPlayers = append(asPlayers, spy)
		asSpies = append(asSpies, spy)
	}

	shoe := NewShoe(1)
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

	t.Run("Should advance the shoe cursor", func(t *testing.T) {
		got := shoe.cursor
		want := 4 + cardsUsed

		assert.Equal(t, want, got)
	})
}
