package dna

const SeedBase int = 100
const SeedMax int = 200

type Base byte

type Codon []byte

var CodonStart Codon = Codon("<")
var CodonStop Codon = Codon(">")

type Codex []Codon

func (c Codex) String() string {
	str := ""
	for _, codon := range c {
		str += string(codon)
	}
	return str
}

type CodexGigas []Codex

type Block interface {
	Bases() [4]Base
	Encoding() map[Base]map[Base]map[Base]Codon
	Random() *DNA
	Match(Base) Base
	Decode([]Base) (Codon, error)
}
