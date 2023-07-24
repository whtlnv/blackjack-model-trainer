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
	shoeIndex := 0
	player.Games[0].SetHand(hand)

	for i := 0; i < len(player.Games); i++ {
		game := player.Games[i]
		shoeIndex = player.playGame(game, dealerHand, shoe, shoeIndex)
	}

	return shoeIndex
}

// Private methods

func (player *Player) subtractFromBankroll(bet int) {
	player.Bankroll -= bet
}

func (player *Player) getAction(game *Game, dealerHand Hand) PlayerAction {
	// a split game with 1 card should always hit
	if len(game.hand) < 2 {
		return Hit
	}

	// a doubled game with 3 cards should always stand
	if game.IsDoubled && len(game.hand) > 2 {
		return Stand
	}

	idealAction := player.strategy.Play(game.hand, dealerHand)

	// if the ideal action is to double/split, but we can't afford it, hit instead
	requiresBet := idealAction == Double || idealAction == Split
	if requiresBet && player.Bankroll < game.bet {
		return Hit
	}

	return idealAction
}

func (player *Player) playGame(game *Game, dealerHand Hand, shoe Shoeish, shoeIndex int) int {
	for {
		action := player.getAction(game, dealerHand)
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

	return shoeIndex
}
