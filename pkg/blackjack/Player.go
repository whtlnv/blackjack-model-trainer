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
	return true, bet
}
