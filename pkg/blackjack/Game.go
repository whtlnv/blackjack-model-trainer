package blackjack

type Game struct {
	hand      Hand
	bet       int
	IsDoubled bool
	IsSplit   bool
}

// Factory

func NewGame(bet int) *Game {
	game := &Game{}
	game.bet = bet

	game.IsDoubled = false
	game.IsSplit = false

	return game
}

// Public methods

func (game *Game) SetHand(hand Hand) {
	game.hand = hand
}

func (game *Game) Hit(card Card) {
	game.hand.Deal(card)
}

func (game *Game) Double(card Card) {
	game.bet *= 2
	game.Hit(card)

	game.IsDoubled = true
}

func (game *Game) Split() *Game {
	keepCard := game.hand[0]
	splitCard := game.hand[1]

	game.hand = Hand{keepCard}
	game.IsSplit = true

	splitGame := NewGame(game.bet)
	splitGame.SetHand(Hand{splitCard})
	splitGame.IsSplit = true

	return splitGame
}

func (game *Game) Resolve(dealerHand *Hand) float64 {
	dealerScore, _ := dealerHand.Score()
	playerScore, _ := game.hand.Score()
	castedBet := float64(game.bet)

	// player busted, you lose
	if playerScore.Low > 21 {
		return 0
	}

	// blackjack, you win
	// if playerScore.High == 21 && len(game.hand) == 2 {
	// 	return game.bet * 2.5
	// }

	// dealer busted, you win
	if dealerScore.Low > 21 {
		return castedBet * 2
	}

	// better hand, you win
	if playerScore.High > dealerScore.High {
		return castedBet * 2
	}

	// same value, you push
	if playerScore.High == dealerScore.High {
		return castedBet
	}

	return 0
}
