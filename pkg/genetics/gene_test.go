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

func TestGeneMerging(t *testing.T) {
	geneA := NewGene([]byte("AAAAA"))
	geneB := NewGene([]byte("BBBBB"))

	randomizerMock := &randomizerMock{}
	randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()
	randomizerMock.On("EventDidHappen", 0.5).Return(false).Once()
	randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()
	randomizerMock.On("EventDidHappen", 0.5).Return(false).Once()
	randomizerMock.On("EventDidHappen", 0.5).Return(true).Once()

	t.Run("Should merge two genes", func(t *testing.T) {
		want := []byte("ABABA")
		got := geneA.Merge(geneB, randomizerMock)
		assert.Equal(t, want, got.Raw())
	})
}
