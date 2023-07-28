package blackjack

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
	// playersInGame := lo.Filter(table.Players, func(player Playerish, index int) bool {
	// 	willBet, _ := player.Bet()
	// 	return willBet
	// })

	// if table.Shoe.NeedsReshuffle {
	// 	table.Shoe.shuffle()
	// }

	// initialHands := table.Shoe.Peek(4)
	// table.Shoe.AdvanceCursor(4)

	// playerHand := Hand{initialHands[0], initialHands[2]}
	// dealerHand := Hand{initialHands[1], initialHands[3]}
	// dealerHand[1].SetHole()

	// play all hands
	// play dealer hand
	// resolve all hands
}

// Private methods

func (table *Table) dealHands() (playerHand Hand, dealerHand Hand) {
	// if table.Shoe.NeedsReshuffle {
	// 	table.Shoe.shuffle()
	// }

	topCards := table.Shoe.Peek(4)
	table.Shoe.AdvanceCursor(4)

	playerHand = Hand{topCards[0], topCards[2]}
	dealerHand = Hand{topCards[1], topCards[3]}
	dealerHand[1].SetHole()

	return playerHand, dealerHand
}
