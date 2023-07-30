package genetics

// TODO: move somewhere else

type Randomizerish interface {
	EventDidHappen(probability float64) bool
}

// /TODO

type Gene struct {
	raw []byte
}

func NewGene(raw []byte) *Gene {
	return &Gene{raw}
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

	return NewGene(merged)
}
