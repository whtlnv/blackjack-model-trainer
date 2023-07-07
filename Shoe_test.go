package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShoeInitialization(t *testing.T) {
	s := NewShoe(1)

	t.Run("Shoe has 52 cards", func(t *testing.T) {
		got := s.Size()
		want := 52
		assert.Equal(t, got, want)
	})
}
