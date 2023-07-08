package main

import "github.com/samber/lo"

type Hand []Card

func (hand Hand) Values() []int {
	var calculateSum func(cards []Card, index int, currentSum int) []int

	calculateSum = func(cards []Card, index int, currentSum int) []int {
		if index >= len(cards) {
			return []int{currentSum}
		}

		card := cards[index]
		values := card.Value()
		var sums []int
		for _, value := range values {
			moreSums := calculateSum(cards, index+1, currentSum+value)
			sums = append(sums, moreSums...)
		}

		return sums
	}

	return lo.Uniq(calculateSum(hand, 0, 0))
}

func (hand Hand) Score() (score int, isBusted bool) {
	values := hand.Values()

	notBusted := lo.Filter(values, func(value int, _ int) bool {
		return value <= 21
	})

	if len(notBusted) == 0 {
		return lo.Max(values), true
	}

	return lo.Max(notBusted), false
}
