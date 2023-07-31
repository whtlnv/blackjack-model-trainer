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

func TestSortByFitness(t *testing.T) {
	t.Run("Should sort a list of candidates by fitness", func(t *testing.T) {
		input := []*Candidate{
			{Fitness: 0.5},
			{Fitness: 0.0},
			{Fitness: 0.02},
			{Fitness: 1.0},
		}

		want := []*Candidate{
			{Fitness: 1.0},
			{Fitness: 0.5},
			{Fitness: 0.02},
			{Fitness: 0.0},
		}
		SortByFitness(input)

		assert.Equal(t, want, input)
	})
}

func TestRemoveWorstPerformers(t *testing.T) {
	t.Run("Should remove the worst candidates from a list", func(t *testing.T) {
		input := []*Candidate{
			{Fitness: 1.0},
			{Fitness: 0.5},
			{Fitness: 0.02},
			{Fitness: 0.0},
		}

		want := []*Candidate{
			{Fitness: 1.0},
			{Fitness: 0.5},
		}
		got := RemoveWorstPerformers(input, 0.5)

		assert.Equal(t, want, got)
	})
}
