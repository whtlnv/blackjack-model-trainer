package blackjack

type Player struct {
	strategy Strategyish
	Bankroll int
}

// Factory

func NewPlayer(strategy Strategyish) Player {
	player := Player{}
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
