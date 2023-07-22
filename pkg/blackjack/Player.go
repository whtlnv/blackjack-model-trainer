package blackjack

type Player struct {
	strategy Strategyish
	Bankroll int
}

// Factory

func NewPlayer(strategy Strategyish) *Player {
	player := &Player{}
	player.strategy = strategy
	player.Bankroll = strategy.GetInitialBankroll()

	return player
}

// Public methods

func (player *Player) Bet() (bool, int) {
	bet := player.strategy.Bet()
	if bet > player.Bankroll {
		return false, 0
	}

	player.Bankroll -= bet

	return true, bet
}

func (player *Player) Play(hand Hand, dealerHand Hand, shoe Shoeish) int {
	// store active game
	// for each hand
	//   get action
	//	 if split
	//     split hand
	//   if double
	//     double hand
	//   if hit
	//     hit hand
	//   if stay
	//     stay hand
	//   if bust
	//     bust hand
	//   if blackjack
	//     blackjack hand

	shoeIndex := 0
	hands := []Hand{hand}
	i := 0

	for {
		action := player.strategy.Play(hands[i], dealerHand)
		switch action {
		case Hit:
			nextCard := shoe.Peek(shoeIndex + 1)[shoeIndex]
			hands[i].Deal(nextCard)
			shoeIndex++
		case Stand:
			return shoeIndex
		default:
			panic("Action not implemented")
		}
	}
}
