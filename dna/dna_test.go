package dna

import (
	"testing"

	"github.com/wmiller848/GoGP/gene"
)

var Start []byte = []byte{0x00, 0x00, 0x00}
var A []byte = []byte{0x40, 0x00, 0x00}
var B []byte = []byte{0x80, 0x00, 0x00}
var End []byte = []byte{0xc0, 0x00, 0x00}

//func InitBlock(t *testing.T) *Block4x3 {
//codex := Codex{Codon("a"), Codon("b")}
//blk, err := NewBlock4x3(Block4x3Bases, codex)
//if err != nil {
//t.Error(err.Error())
//return nil
//}
//return blk
//}

func ConstructStrand(codex Codex) []byte {
	strand := []byte{}
	for _, codon := range codex {
		strand = append(strand, codon...)
	}
	return strand
}

func TestDNASingleStrand(t *testing.T) {
	// 2nd frame reading
	ying := ConstructStrand(Codex{
		[]byte{0xc0},
		Start,
		A, A, B, B, A, B, B,
		End,
	})
	yang := []byte{}
	blk := InitBlock(t)
	dna := &DNA{
		StrandYing: gene.GenericGene(ying),
		StrandYang: gene.GenericGene(yang),
		Block:      blk,
	}

	gns, err := dna.MarshalGenes()
	if err != nil {
		t.Error(err.Error())
		return
	}
	AssertStr(t, string(gns), "aabbabb")
}

func TestDNADoubleStrand(t *testing.T) {
	// 2nd frame reading
	ying := ConstructStrand(Codex{
		[]byte{0xc0},
		Start,
		A, A, B, B, A, B, B,
		End,
	})
	// 1st frame reading
	yang := ConstructStrand(Codex{
		Start,
		B, B, B, B, A,
		End,
	})
	blk := InitBlock(t)
	dna := &DNA{
		StrandYing: gene.GenericGene(ying),
		StrandYang: gene.GenericGene(yang),
		Block:      blk,
	}

	gns, err := dna.MarshalGenes()
	if err != nil {
		t.Error(err.Error())
		return
	}
	AssertStr(t, string(gns), "bbbba")
}

func TestDNADoubleStrandComplex(t *testing.T) {
	// 2nd frame reading
	ying := ConstructStrand(Codex{
		[]byte{0xc0},
		Start,
		A, A, B, B, A, B, B,
		End,
		Start,
		A, B, B, A,
		End,
	})
	// 1st frame reading
	yang := ConstructStrand(Codex{
		Start,
		B, B, B, B, A,
		End,
	})
	blk := InitBlock(t)
	dna := &DNA{
		StrandYing: gene.GenericGene(ying),
		StrandYang: gene.GenericGene(yang),
		Block:      blk,
	}

	gns, err := dna.MarshalGenes()
	if err != nil {
		t.Error(err.Error())
		return
	}
	AssertStr(t, string(gns), "bbbbaabba")
}
