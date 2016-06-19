package dna

import (
	"errors"
	"math"

	"github.com/wmiller848/GoGP/gene"
	"github.com/wmiller848/GoGP/util"
)

var Block4x3Bases [4]Base = [4]Base{0x00, 0x40, 0x80, 0xc0}

type Block4x3 struct {
	bases    [4]Base
	encoding map[Base]map[Base]map[Base]Codon
}

func NewBlock4x3(bases [4]Base, codex Codex) (*Block4x3, error) {
	baseSize := int(math.Pow(4, 3))
	if len(codex) > baseSize-2 {
		return nil, errors.New("Codexs can have a max of 62 items")
	}
	blk := &Block4x3{
		bases:    bases,
		encoding: make(map[Base]map[Base]map[Base]Codon),
	}

	i := 0
	u := 0
	// First Encoding Codon is start
	codexPool := append([]Codon{CodonStart}, codex...)
	// Last Encoding Codon is stop
	codexPool = append(codexPool, CodonStop)
	dist := baseSize / len(codexPool)
	cursor := codexPool[u]
	for _, b1 := range bases {
		for _, b2 := range bases {
			for _, b3 := range bases {
				if blk.encoding[b1] == nil {
					blk.encoding[b1] = make(map[Base]map[Base]Codon)
				}
				if blk.encoding[b1][b2] == nil {
					blk.encoding[b1][b2] = make(map[Base]Codon)
				}
				i++
				if i%dist == 0 {
					u++
					if u > len(codexPool)-1 {
						u = 0
					}
					cursor = codexPool[u]
				}
				blk.encoding[b1][b2][b3] = cursor
			}
		}
	}
	return blk, nil
}

func (b *Block4x3) Bases() [4]Base {
	return b.bases
}

func (b *Block4x3) Encoding() map[Base]map[Base]map[Base]Codon {
	return b.encoding
}

func (b *Block4x3) Random() *DNA {
	dna := &DNA{
		StrandYing: gene.GenericGene{},
		StrandYang: gene.GenericGene{},
		Block:      b,
	}

	seedYing := int(util.RandomNumber(SeedBase, SeedMax))
	for i := 0; i < seedYing; i++ {
		pick := byte(util.RandomNumber(0, 255))
		dna.StrandYing = append(dna.StrandYing, pick)
	}
	seedYang := int(util.RandomNumber(SeedBase, SeedMax))
	for i := 0; i < seedYang; i++ {
		pick := byte(util.RandomNumber(0, 255))
		dna.StrandYang = append(dna.StrandYang, pick)
	}
	return dna
}

func (b *Block4x3) Match(frag Base) Base {
	if frag >= b.bases[0] && frag < b.bases[1] {
		return b.bases[0]
	} else if frag >= b.bases[1] && frag < b.bases[2] {
		return b.bases[1]
	} else if frag >= b.bases[2] && frag < b.bases[3] {
		return b.bases[2]
	} else if frag >= b.bases[3] {
		return b.bases[3]
	}
	return 0x00
}

func (b *Block4x3) Decode(strand []Base) (Codon, error) {
	if len(strand) != 3 {
		return nil, errors.New("Invalid strand size, must be 3 bytes")
	}
	c0 := b.Match(strand[0])
	c1 := b.Match(strand[1])
	c2 := b.Match(strand[2])
	return b.encoding[c0][c1][c2], nil
}
