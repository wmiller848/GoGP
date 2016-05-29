package dna

type Strand []byte

type DNA struct {
	StrandYing Strand
	StrandYang Strand
	Blocks     []Block
}
