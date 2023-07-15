package blackjack

type Player struct {
	strategy Strategy
	Bankroll int
}

// Factory

func NewPlayer(strategy Strategy) Player {
	player := Player{}
	player.strategy = strategy
	player.Bankroll = strategy.InitialBankroll

	return player
}

// Public methods

func (player *Player) Bet() (bool, int) {
	bet := player.strategy.Bet()
	if bet > player.Bankroll {
		return false, 0
	}
	return true, bet
}
