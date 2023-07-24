package blackjack

type Game struct {
	hand    Hand
	bet     int
	Doubled bool
}

// Factory

func NewGame(bet int) *Game {
	game := &Game{}
	game.bet = bet

	game.Doubled = false

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

	game.Doubled = true
}
