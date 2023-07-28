package blackjack

import (
	"github.com/samber/lo"
)

type Table struct {
	Players []Playerish
	Shoe    Shoeish
}

// Factory

func NewTable(players []Playerish, shoe Shoeish) *Table {
	return &Table{Players: players, Shoe: shoe}
}

// Public methods

func (table *Table) Run() {
	lo.ForEach(table.Players, func(player Playerish, index int) {
		player.Bet()
	})

	playerHand, dealerHand := table.dealHands()

	table.playAllGames(playerHand, dealerHand)

	dealerHand.Reveal()

	// play dealer hand
	// advance shoe cursor

	// resolve all hands
}

// Private methods

func (table *Table) dealHands() (playerHand Hand, dealerHand Hand) {
	if table.Shoe.NeedsReshuffle() {
		table.Shoe.Shuffle()
	}

	topCards := table.Shoe.Peek(4)
	table.Shoe.AdvanceCursor(4)

	playerHand = Hand{topCards[0], topCards[2]}
	dealerHand = Hand{topCards[1], topCards[3]}
	dealerHand[1].SetHole()

	return playerHand, dealerHand
}

func (table *Table) playAllGames(playerHand Hand, dealerHand Hand) {
	index := lo.Reduce(table.Players, func(acc int, player Playerish, _ int) int {
		usedCards := player.Play(playerHand, dealerHand, table.Shoe)
		return lo.Max([]int{acc, usedCards})
	}, 0)

	table.Shoe.AdvanceCursor(index)
}
