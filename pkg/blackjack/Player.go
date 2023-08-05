package blackjack

type Playerish interface {
	Bet() (bool, int)
	Play(hand Hand, dealerHand Hand, shoe Shoeish) int
	Resolve(dealerHand Hand)
	GetStatistics() PlayerStatistics
}

type Player struct {
	strategy Strategyish
	Bankroll float64
	Games    []*Game

	GamesPlayed int
	GamesWon    int
	GamesLost   int
}

type PlayerStatistics struct {
	GamesPlayed int
	GamesWon    int
	GamesLost   int
	GamesPushed int

	InitialBankroll float64
	Bankroll        float64
	BankrollDelta   float64
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

func (player *Player) Bet() (willPlay bool, bet int) {
	bet = player.strategy.Bet()
	if float64(bet) > player.Bankroll {
		return false, 0
	}

	player.subtractFromBankroll(bet)
	player.Games = []*Game{NewGame(bet)}

	return true, bet
}

func (player *Player) Play(hand Hand, dealerHand Hand, shoe Shoeish) int {
	if len(player.Games) == 0 {
		return 0
	}

	shoeIndex := 0
	player.Games[0].SetHand(hand)

	for i := 0; i < len(player.Games); i++ {
		game := player.Games[i]
		shoeIndex = player.playGame(game, dealerHand, shoe, shoeIndex)
	}

	return shoeIndex
}

func (player *Player) Resolve(dealerHand Hand) {
	totalBet := 0
	winnings := 0.0

	for _, game := range player.Games {
		totalBet += game.bet
		winnings += game.Resolve(&dealerHand)
	}

	player.updateStatistics(float64(totalBet), winnings)

	player.Bankroll += winnings
	player.Games = []*Game{}
}

func (player *Player) GetStatistics() PlayerStatistics {
	gamesPushed := player.GamesPlayed - player.GamesWon - player.GamesLost
	bankrollDelta := player.Bankroll - player.strategy.GetInitialBankroll()

	return PlayerStatistics{
		GamesPlayed:     player.GamesPlayed,
		GamesWon:        player.GamesWon,
		GamesLost:       player.GamesLost,
		GamesPushed:     gamesPushed,
		InitialBankroll: player.strategy.GetInitialBankroll(),
		Bankroll:        player.Bankroll,
		BankrollDelta:   bankrollDelta,
	}
}

func (player *Player) GetStrategy() Strategyish {
	return player.strategy
}

// Private methods

func (player *Player) subtractFromBankroll(bet int) {
	player.Bankroll -= float64(bet)
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

	// dealer has a blackjack, we lose
	if dealerHand.IsBlackjack() && dealerHand.GetHoleCard().rank != Ace {
		return Stand
	}

	idealAction := player.strategy.Play(game.hand, dealerHand)

	// if the ideal action is to double/split, but we can't afford it, hit instead
	requiresBet := idealAction == Double || idealAction == SplitOrHit
	if requiresBet && player.Bankroll < float64(game.bet) {
		return Hit
	}

	return idealAction
}

func (player *Player) playGame(game *Game, dealerHand Hand, shoe Shoeish, shoeIndex int) int {
	for {
		action := player.getAction(game, dealerHand)
		nextCard := shoe.Peek(shoeIndex + 1)[shoeIndex]

		if action == SplitOrHit {
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

func (player *Player) updateStatistics(bet float64, won float64) {
	player.GamesPlayed++
	if won > bet {
		player.GamesWon++
	} else if won < bet {
		player.GamesLost++
	}
}
