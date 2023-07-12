package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAceCardValue(t *testing.T) {
	card := Card{Ace, Spades}

	got := card.Value()
	want := []int{1, 11}

	assert.ElementsMatch(t, want, got)
}

func TestFaceCardValue(t *testing.T) {
	card := Card{Jack, Hearts}

	got := card.Value()
	want := []int{10}

	assert.ElementsMatch(t, want, got)
}

func TestNumberCardValue(t *testing.T) {
	card := Card{Two, Clubs}

	got := card.Value()
	want := []int{2}

	assert.ElementsMatch(t, want, got)
}
