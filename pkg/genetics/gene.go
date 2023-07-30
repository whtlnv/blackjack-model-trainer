package genetics

// TODO: move somewhere else

type Randomizerish interface {
	EventDidHappen(probability float64) bool
	PickOne(options []byte) byte
}

// /TODO

type Gene struct {
	raw          []byte
	sequencing   [][]byte
	mutationRate float64
}

func NewGene(raw []byte, sequencing [][]byte, mutationRate float64) *Gene {
	return &Gene{raw, sequencing, mutationRate}
}

func (gene *Gene) Raw() []byte {
	return gene.raw
}

func (gene *Gene) Merge(other *Gene, randomizer Randomizerish) *Gene {
	merged := make([]byte, len(gene.raw))
	myGeneWasChosen := func() bool { return randomizer.EventDidHappen(0.5) }

	for i := 0; i < len(gene.raw); i++ {
		if myGeneWasChosen() {
			merged[i] = gene.raw[i]
		} else {
			merged[i] = other.raw[i]
		}
	}

	return NewGene(merged, gene.sequencing, gene.mutationRate).Mutate(randomizer)
}

func (gene *Gene) Mutate(randomizer Randomizerish) *Gene {
	mutated := make([]byte, len(gene.raw))
	shouldMutate := func() bool { return randomizer.EventDidHappen(gene.mutationRate) }

	for i := 0; i < len(gene.raw); i++ {
		if shouldMutate() {
			mutated[i] = randomizer.PickOne(gene.sequencing[i])
		} else {
			mutated[i] = gene.raw[i]
		}
	}

	gene.raw = mutated
	return gene
}
