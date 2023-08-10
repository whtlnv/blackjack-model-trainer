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

	GamesSeen   int
	GamesPlayed int
	GamesWon    int
	GamesLost   int
	GamesPushed int
}

type PlayerStatistics struct {
	Strategy []byte

	GamesSeen   int
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
	player.GamesSeen++

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
		Strategy:        player.strategy.GetEncodedStrategy(),
		GamesSeen:       player.GamesSeen,
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

	if idealAction == SplitOrStand && player.Bankroll < float64(game.bet) {
		return Stand
	}

	// // if the ideal action is to splitOrHit, but we already have, hit instead
	// if idealAction == SplitOrHit && game.IsSplit {
	// 	return Hit
	// }

	// // if the ideal action is to splitOrStand, but we already have, stand instead
	// if idealAction == SplitOrStand && game.IsSplit {
	// 	return Stand
	// }

	return idealAction
}

func (player *Player) split(game *Game) (cardsTaken int) {
	splitGame := game.Split()
	player.subtractFromBankroll(splitGame.bet)
	player.Games = append(player.Games, splitGame)
	return 0
}

func (player *Player) double(game *Game, nextCard Card) (cardsTaken int) {
	player.subtractFromBankroll(game.bet)
	game.Double(nextCard)
	return 1
}

func (player *Player) hit(game *Game, nextCard Card) (cardsTaken int) {
	game.Hit(nextCard)
	return 1
}

func (player *Player) takeTurn(game *Game, dealerHand Hand, shoe Shoeish, shoeIndex int) (cardsTaken int, stop bool) {
	getNextCard := func() Card {
		peek := shoe.Peek(shoeIndex + 1)
		if len(peek) < shoeIndex+1 {
			panic("Shoe is out of cards")
		}
		return peek[shoeIndex]
	}

	evaluateAction := func(action PlayerAction) (cardsTaken int, stop bool) {
		switch action {
		case Stand:
			return 0, true
		case SplitOrHit:
			return player.split(game), false
		case SplitOrStand:
			return player.split(game), false
		case Double:
			return player.double(game, getNextCard()), false
		case Hit:
			return player.hit(game, getNextCard()), false
		default:
			panic("Unknown action")
		}
	}

	return evaluateAction(player.getAction(game, dealerHand))
}

func (player *Player) playGame(game *Game, dealerHand Hand, shoe Shoeish, shoeIndex int) int {
	for {
		cardsUsed, stop := player.takeTurn(game, dealerHand, shoe, shoeIndex)
		shoeIndex += cardsUsed

		if stop {
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
	} else {
		player.GamesPushed++
	}
}
