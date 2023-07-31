package genetics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeFitnessList(t *testing.T) {
	t.Run("Should return a list of normalized fitness values", func(t *testing.T) {
		input := []*Candidate{
			{Fitness: 0.0},
			{Fitness: 1.0},
			{Fitness: 25.0},
			{Fitness: 50.0},
		}

		want := []*Candidate{{Fitness: 0.0}, {Fitness: 0.02}, {Fitness: 0.5}, {Fitness: 1.0}}
		got := NormalizeFitnessList(input)

		assert.Equal(t, want, got)
	})
}
