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

func TestParthenogenesis(t *testing.T) {
	t.Run("Should create a new population by cloning the best candidates", func(t *testing.T) {
		randomizerMock := &RandomizerMock{}
		randomizerMock.On("EventDidHappen", 1.0).Return(true).Once()
		randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()
		randomizerMock.On("EventDidHappen", 0.02).Return(false).Once()
		randomizerMock.On("EventDidHappen", 0.0).Return(false).Once()

		input := []*Candidate{
			{Chromosome: &Chromosome{raw: []byte("AAA")}, Fitness: 1.0},
			{Chromosome: &Chromosome{raw: []byte("BBB")}, Fitness: 0.5},
			{Chromosome: &Chromosome{raw: []byte("CCC")}, Fitness: 0.02},
			{Chromosome: &Chromosome{raw: []byte("DDD")}, Fitness: 0.0},
		}

		want := []*Candidate{
			{Chromosome: &Chromosome{raw: []byte("AAA")}, Fitness: -1.0},
			{Chromosome: &Chromosome{raw: []byte("BBB")}, Fitness: -1.0},
		}
		got := Parthenogenesis(input, randomizerMock)

		assert.Equal(t, want, got)
		assert.NotEqual(t, &input[0], &got[0])
	})
}

func TestCrossover(t *testing.T) {
	t.Run("Should create a new population by merging the best candidates", func(t *testing.T) {
		randomizerMock := &RandomizerMock{}
		// mate A and B
		randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()
		// do not mate any other pair
		randomizerMock.On("EventDidHappen", 0.02).Return(false)
		randomizerMock.On("EventDidHappen", 0.0).Return(false)
		randomizerMock.On("EventDidHappen", 0.01).Return(false)
		// number of children
		randomizerMock.On("NumberBetween", 1, 10).Return(2)
		// first child
		randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()
		randomizerMock.On("EventDidHappen", 0.5).Return(false).Once()
		randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()
		// second child
		randomizerMock.On("EventDidHappen", 0.5).Return(false).Once()
		randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()
		randomizerMock.On("EventDidHappen", 0.5).Return(false).Once()

		input := []*Candidate{
			{Chromosome: &Chromosome{raw: []byte("AAA")}, Fitness: 1.0},
			{Chromosome: &Chromosome{raw: []byte("BBB")}, Fitness: 0.5},
			{Chromosome: &Chromosome{raw: []byte("CCC")}, Fitness: 0.02},
			{Chromosome: &Chromosome{raw: []byte("DDD")}, Fitness: 0.0},
		}

		want := []*Candidate{
			{Chromosome: &Chromosome{raw: []byte("ABA")}, Fitness: -1.0},
			{Chromosome: &Chromosome{raw: []byte("BAB")}, Fitness: -1.0},
		}
		got := Crossover(input, randomizerMock)

		assert.Equal(t, want, got)
		assert.NotEqual(t, &input[0], &got[0])
		randomizerMock.AssertExpectations(t)
	})
}
