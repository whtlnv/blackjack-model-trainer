package blackjack

type Player struct {
	strategy Strategyish
	Bankroll int
	Games    []*Game
}

// Factory

func NewPlayer(strategy Strategyish) *Player {
	player := &Player{}
	player.strategy = strategy
	player.Bankroll = strategy.GetInitialBankroll()

	player.Games = []*Game{}

	return player
}

// Public methods

func (player *Player) Bet() (bool, int) {
	bet := player.strategy.Bet()
	if bet > player.Bankroll {
		return false, 0
	}

	player.subtractFromBankroll(bet)
	player.Games = []*Game{NewGame(bet)}

	return true, bet
}

func (player *Player) Play(hand Hand, dealerHand Hand, shoe Shoeish) int {
	getAction := func(hand Hand, dealerHand Hand) PlayerAction {
		if len(hand) < 2 {
			return Hit
		}

		return player.strategy.Play(hand, dealerHand)
	}

	shoeIndex := 0
	player.Games[0].SetHand(hand)

	for i := 0; i < len(player.Games); i++ {
		for {
			game := player.Games[i]
			action := getAction(game.hand, dealerHand)
			nextCard := shoe.Peek(shoeIndex + 1)[shoeIndex]

			if action == Split {
				splitGame := game.Split()
				player.subtractFromBankroll(splitGame.bet)
				player.Games = append(player.Games, splitGame)
			}

			if action == Double {
				player.subtractFromBankroll(game.bet)
				game.Double(nextCard)
				shoeIndex++
			}

			if action == Hit {
				game.Hit(nextCard)
				shoeIndex++
			}

			if action == Stand {
				break
			}
		}
	}

	return shoeIndex
}

// Private methods

func (player *Player) subtractFromBankroll(bet int) {
	player.Bankroll -= bet
}
