package genetics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type randomizerMock struct {
	mock.Mock
}

func (r *randomizerMock) EventDidHappen(probability float64) bool {
	args := r.Called(probability)
	return args.Bool(0)
}

func (r *randomizerMock) PickOne(options []byte) byte {
	args := r.Called(options)
	return args.Get(0).(byte)
}

func TestGeneMerging(t *testing.T) {
	bases := []byte("ABC")
	sequencing := [][]byte{bases, bases, bases, bases, bases}
	mutationRate := 0.1
	geneA := NewGene([]byte("AAAAA"), sequencing, mutationRate)
	geneB := NewGene([]byte("BBBBB"), sequencing, mutationRate)

	t.Run("Should merge two genes", func(t *testing.T) {
		randomizerMock := &randomizerMock{}
		// alternate between A and B
		randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()
		randomizerMock.On("EventDidHappen", 0.5).Return(false).Once()
		randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()
		randomizerMock.On("EventDidHappen", 0.5).Return(false).Once()
		randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()
		// don't mutate
		randomizerMock.On("EventDidHappen", mutationRate).Return(false)

		want := []byte("ABABA")
		got := geneA.Merge(geneB, randomizerMock)
		assert.Equal(t, want, got.Raw())
	})

	t.Run("Should mutate offspring genes", func(t *testing.T) {
		randomizerMock := &randomizerMock{}
		// pick all my genes
		randomizerMock.On("EventDidHappen", 0.5).Return(true)
		// mutate only first gene, keep the rest
		randomizerMock.On("EventDidHappen", mutationRate).Return(true).Once()
		randomizerMock.On("EventDidHappen", mutationRate).Return(false)
		// pick C for the first gene
		randomizerMock.On("PickOne", bases).Return(byte('C')).Once()

		want := []byte("CAAAA")
		got := geneA.Merge(geneB, randomizerMock)
		assert.Equal(t, want, got.Raw())
	})
}
