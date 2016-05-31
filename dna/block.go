package dna

import (
	"errors"
	"fmt"
	"math"
)

//func randomOperator() byte {
//switch util.RandomNumber(0, 3) {
//case 0:
//return byte('+')
//case 1:
//return byte('-')
//case 2:
//return byte('*')
//case 3:
//return byte('/')
//default:
//return randomOperator()
//}
//}

//func randomVariable() byte {
//return byte(util.RandomNumber(0, 9))
//}

//func randomNumber() byte {
//return byte(util.RandomNumber(0, 9))
//}

//const (
//StageSpawn  int = 1
//StageAlive  int = 2
//StageDieing int = 3
//StageDead   int = 3
//)

//type Base interface {
//Value() string
//}

type Base byte

type BaseNode struct {
	Children [2]*BaseNode
	Depth    int
}

type Codon []byte

var CodonStart Codon = Codon("<")
var CodonStop Codon = Codon(">")

type Block interface {
	Bases() [4]Base
	Encoding() map[Base]map[Base]map[Base]Codon
	Random() *DNA
}

type Block4x3 struct {
	bases    [4]Base
	encoding map[Base]map[Base]map[Base]Codon
}

func NewBlock4x3(bases [4]Base, codexs []Codon) (*Block4x3, error) {
	baseSize := int(math.Pow(4, 3))
	if len(codexs) > baseSize-1 {
		return nil, errors.New("Codexs can have a max of 63 items")
	}
	blk := &Block4x3{
		bases:    bases,
		encoding: make(map[Base]map[Base]map[Base]Codon),
	}

	dist := baseSize / len(codexs)
	fmt.Println("Distribution is", dist)
	fmt.Println("=== START ===")
	i := 0
	u := 0
	cursor := codexs[u]
	for _, b1 := range bases {
		for _, b2 := range bases {
			for _, b3 := range bases {
				if blk.encoding[b1] == nil {
					blk.encoding[b1] = make(map[Base]map[Base]Codon)
				}
				if blk.encoding[b1][b2] == nil {
					blk.encoding[b1][b2] = make(map[Base]Codon)
				}
				blk.encoding[b1][b2][b3] = cursor
				i++
				if i%dist == 0 {
					u++
					if u > len(codexs)-1 {
						u = 0
					}
					cursor = codexs[u]
				}
			}
		}
	}
	// First is start and assigned value
	// Last Encoding is always a stop
	blk.encoding[bases[3]][bases[3]][bases[3]] = Codon(CodonStop)
	fmt.Println(blk.encoding)
	fmt.Println("=== STOP ===")

	return blk, nil
}

func (b *Block4x3) Bases() [4]Base {
	return b.bases
}

func (b *Block4x3) Encoding() map[Base]map[Base]map[Base]Codon {
	return b.encoding
}

func (b *Block4x3) Random() *DNA {
	return &DNA{}
}
