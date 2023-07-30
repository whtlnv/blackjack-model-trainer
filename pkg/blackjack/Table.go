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

	dealerGame := table.playDealerHand(dealerHand)

	lo.ForEach(table.Players, func(player Playerish, index int) {
		player.Resolve(dealerGame)
	})
}

func (table *Table) RunMany(games int) {
	for i := 0; i < games; i++ {
		tableBankroll := lo.Reduce(table.Players, func(acc float64, player Playerish, _ int) float64 {
			return acc + player.GetStatistics().Bankroll
		}, 0.0)

		if tableBankroll <= 0.0 {
			break
		}

		table.Run()
	}
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

func (table *Table) playDealerHand(dealerHand Hand) Hand {
	dealerHand.Reveal()

	for {
		dealerScore, isBusted := dealerHand.Score()
		if isBusted || dealerScore.High >= 17 {
			break
		}

		dealerHand = append(dealerHand, table.Shoe.Peek(1)[0])
		table.Shoe.AdvanceCursor(1)
	}

	return dealerHand
}
