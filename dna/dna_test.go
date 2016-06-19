package dna

import (
	"testing"

	"github.com/wmiller848/GoGP/gene"
)

func TestDNA(t *testing.T) {
	dna := &DNA{
		StrandYing: gene.GenericGene(""),
		StrandYang: gene.GenericGene(""),
	}
}
