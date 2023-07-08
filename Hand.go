package main

import "github.com/samber/lo"

type Hand []Card

func (hand Hand) Values() []int {
	var sums []int
	var calculateSum func(cards []Card, index int, currentSum int)

	calculateSum = func(cards []Card, index int, currentSum int) {
		if index >= len(cards) {
			sums = append(sums, currentSum)
			return
		}

		card := cards[index]
		values := card.Value()
		for _, value := range values {
			calculateSum(cards, index+1, currentSum+value)
		}
	}

	calculateSum(hand, 0, 0)

	return lo.Uniq(sums)
}
