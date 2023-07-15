package main

import (
	"bytes"
	"fmt"

	"github.com/whtlnv/blackjack-model-trainer/pkg/blackjack"
)

func main() {
	shoe := blackjack.NewShoe(1)

	fmt.Println("Got a shoe")
	fmt.Println(shoe)

	rawStrategy := bytes.Repeat([]byte("H"), blackjack.HandCount)
	strategy, _ := blackjack.NewStrategy(rawStrategy)
	player := blackjack.NewPlayer(strategy)

	willBet, ammount := player.Bet()

	fmt.Println("Got a player")
	fmt.Println(willBet, ammount)
}
