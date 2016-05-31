package dna

import "github.com/wmiller848/GoGP/gene"

type DNA struct {
	StrandYing gene.GenericGene
	StrandYang gene.GenericGene
	Blocks     []Block
}

func (d *DNA) MarshalGenes() ([]gene.Gene, error) {
	return []gene.Gene{}, nil
}

func (d *DNA) MarshalHelix() ([]byte, error) {
	return []byte{}, nil
}
