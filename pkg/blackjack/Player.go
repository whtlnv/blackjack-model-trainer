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

	player.Bankroll -= bet
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

			if action == Split {
				splitHand := game.Split()
				player.Games = append(player.Games, splitHand)
			}

			if action == Hit {
				nextCard := shoe.Peek(shoeIndex + 1)[shoeIndex]
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
