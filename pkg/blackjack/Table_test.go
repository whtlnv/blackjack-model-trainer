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

// func TestTableRun(t *testing.T) {
// 	spy := &playerSpy{}
// 	spy.On("Bet").Return(true, 10)
// 	spy.On("Play", mock.Anything, mock.Anything).Return(Stand)
// 	spy.On("Resolve").Return()

// 	players := []Playerish{spy}
// 	table := NewTable(players, 1)

// 	t.Run("Should run a game with the given players", func(t *testing.T) {
// 		table.Run()
// 		got := len(spy.Calls)
// 		want := 3
// 		assert.Equal(t, want, got)
// 	})
// }
