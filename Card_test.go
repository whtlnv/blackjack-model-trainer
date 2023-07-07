package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardValue(t *testing.T) {
	c := Card{Ace, Spades}
	
	got := c.Value()
	want := []int{1, 11}

	assert.ElementsMatch(t, got, want)
}