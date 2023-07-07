package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAceCardValue(t *testing.T) {
	c := Card{Ace, Spades}
	
	got := c.Value()
	want := []int{1, 11}

	assert.ElementsMatch(t, got, want)
}

func TestFaceCardValue(t *testing.T) {
	c := Card{Jack, Hearts}
	
	got := c.Value()
	want := []int{10}

	assert.ElementsMatch(t, got, want)
}

func TestNumberCardValue(t *testing.T) {
	c := Card{Two, Clubs}
	
	got := c.Value()
	want := []int{2}

	assert.ElementsMatch(t, got, want)
}