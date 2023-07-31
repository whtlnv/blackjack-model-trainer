package genetics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChromosomeMerging(t *testing.T) {
	bases := []byte("ABC")
	sequencing := [][]byte{bases, bases, bases, bases, bases}
	mutationRate := 0.1
	subjectA := NewChromosome([]byte("AAAAA"), sequencing, mutationRate)
	subjectB := NewChromosome([]byte("BBBBB"), sequencing, mutationRate)

	t.Run("Should merge two chromosomes", func(t *testing.T) {
		randomizerMock := &RandomizerMock{}
		// alternate between A and B
		randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()
		randomizerMock.On("EventDidHappen", 0.5).Return(false).Once()
		randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()
		randomizerMock.On("EventDidHappen", 0.5).Return(false).Once()
		randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()
		// don't mutate
		randomizerMock.On("EventDidHappen", mutationRate).Return(false)

		want := []byte("ABABA")
		got := subjectA.Merge(subjectB, randomizerMock)
		assert.Equal(t, want, got.Raw())
	})

	t.Run("Should mutate offspring genes", func(t *testing.T) {
		randomizerMock := &RandomizerMock{}
		// pick all my genes
		randomizerMock.On("EventDidHappen", 0.5).Return(true)
		// mutate only first gene, keep the rest
		randomizerMock.On("EventDidHappen", mutationRate).Return(true).Once()
		randomizerMock.On("EventDidHappen", mutationRate).Return(false)
		// pick C for the first gene
		randomizerMock.On("PickOne", bases).Return(byte('C')).Once()

		want := []byte("CAAAA")
		got := subjectA.Merge(subjectB, randomizerMock)
		assert.Equal(t, want, got.Raw())
	})
}
