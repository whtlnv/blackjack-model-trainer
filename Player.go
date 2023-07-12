package main

type Player struct {
	strategy Strategy
}

// Factory

func NewPlayer(strategy Strategy) Player {
	return Player{strategy}
}

// Public methods

func (player *Player) ShouldBet() (bool, int) {
	return true, 1
}
