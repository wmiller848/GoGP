package dna

import (
	"math"

	"github.com/wmiller848/GoGP/gene"
	"github.com/wmiller848/GoGP/util"
)

type DNA struct {
	StrandYin  gene.GenericGene
	StrandYang gene.GenericGene
	Block      Block
}

func (d *DNA) Mutate() *DNA {
	dna := &DNA{
		StrandYin:  d.StrandYin,
		StrandYang: d.StrandYang,
		Block:      d.Block,
	}
	// Yin
	strandYin := gene.GenericGene{}
	if len(dna.StrandYin) > 0 {
		yingRnd := int(util.RandomNumber(len(dna.StrandYang)/2, len(dna.StrandYin)*2))
		if yingRnd < 24 {
			yingRnd = 24
		}
		for i := 0; i < yingRnd; i++ {
			ii := i % len(dna.StrandYin)
			codon := dna.StrandYin[ii]
			switch util.RandomNumber(0, 19) {
			// Mutate codon
			case 0:
				codon = codon ^ byte(util.RandomNumber(0, 255))
				strandYin = append(strandYin, codon)
			// Omit codon
			case 1:
			// Add extra
			case 3:
				strandYin = append(strandYin, codon)
				strandYin = append(strandYin, byte(util.RandomNumber(0, 255)))
			// No Op
			default:
				strandYin = append(strandYin, codon)
			}
		}
	}
	dna.StrandYin = strandYin

	// Yang
	strandYang := gene.GenericGene{}
	if len(dna.StrandYang) > 0 {
		yangRnd := int(util.RandomNumber(len(dna.StrandYang)/2, len(dna.StrandYang)*2))
		if yangRnd < 24 {
			yangRnd = 24
		}
		for i := 0; i < yangRnd; i++ {
			ii := i % len(dna.StrandYang)
			codon := dna.StrandYang[ii]
			switch util.RandomNumber(0, 19) {
			// Mutate codon
			case 0:
				codon = codon ^ byte(util.RandomNumber(0, 255))
				strandYang = append(strandYang, codon)
			// Omit codon
			case 1:
			// Add extra
			case 2:
				strandYang = append(strandYang, codon)
				strandYang = append(strandYang, byte(util.RandomNumber(0, 255)))
			// No Op
			default:
				strandYang = append(strandYang, codon)
			}
		}
	}
	dna.StrandYang = strandYang

	return dna
}

func (d *DNA) Unwind(strand gene.GenericGene) CodexGigas {
	leng := len(strand)
	codexGigas := CodexGigas{}
	// TODO 3 should be block size
	for i := 0; i < 3; i++ {
		codex := Codex{}
		for j := 0; j < leng; j += 3 {
			t0 := i + 0 + j
			t1 := i + 1 + j
			t2 := i + 2 + j
			if t0 > leng-1 {
				t0 -= leng
			}
			if t1 > leng-1 {
				t1 -= leng
			}
			if t2 > leng-1 {
				t2 -= leng
			}
			strand_frag := []Base{Base(strand[t0]), Base(strand[t1]), Base(strand[t2])}
			codon, _ := d.Block.Decode(strand_frag)
			codex = append(codex, codon)
		}
		codexGigas = append(codexGigas, codex)
	}
	return codexGigas
}

func (d *DNA) Sequence(codexGigas CodexGigas) chan *Sequence {
	chanSeq := make(chan *Sequence)
	go func() {
		for codexID, codex := range codexGigas {
			i := 0
			index := 0
			elements := 0
			reading := false
			codexDecoded := Codex{}
			for _, codon := range codex {
				if string(codon) == string(CodonStart) {
					reading = true
					index = i
				} else if string(codon) != string(CodonStop) && reading == true {
					codexDecoded = append(codexDecoded, codon)
					elements++
				} else if string(codon) == string(CodonStop) && reading == true {
					if len(codexDecoded) == 0 {
						reading = false
						continue
					}
					seq := Sequence{
						Codex:    codexDecoded,
						CodexID:  codexID,
						Index:    index + codexID,
						Elements: elements,
					}
					chanSeq <- &seq
					codexDecoded = Codex{}
					reading = false
					elements = 0
				}
				i++
			}
		}
		close(chanSeq)
	}()
	return chanSeq
}

func (d *DNA) SpliceSequence(chanSeqs [2]chan *Sequence) *SequenceNode {
	var headYin *SequenceNode
	var headYang *SequenceNode
	for j, chanSeq := range chanSeqs {
		var head0 *SequenceNode
		var head1 *SequenceNode
		var head2 *SequenceNode
		for {
			seq, open := <-chanSeq
			if open == false {
				break
			}
			switch seq.CodexID {
			case 0:
				if head0 == nil {
					head0 = seq.Node()
				} else {
					head0 = head0.Merge(seq)
				}
			case 1:
				if head1 == nil {
					head1 = seq.Node()
				} else {
					head1 = head1.Merge(seq)
				}
			case 2:
				if head2 == nil {
					head2 = seq.Node()
				} else {
					head2 = head2.Merge(seq)
				}
			}
		}
		head := head0
		if head == nil || (head1 != nil && head.Index > head1.Index) {
			head = head1
		}
		if head == nil || (head2 != nil && head.Index > head2.Index) {
			head = head2
		}

		if j == 0 {
			headYin = head
		} else {
			headYang = head
		}
	}

	var dnaSeq *SequenceNode
	if headYin != nil && headYang != nil {
		dnaSeqYin := headYin.Clone()
		dnaSeqYang := headYang.Clone()
		var i int = -1
		for {
			if i < dnaSeqYin.Index && dnaSeqYin.Index < dnaSeqYang.Index {
				if i == -1 {
					i = 0
				}
				if dnaSeq == nil {
					dnaSeq = dnaSeqYin.Clone()
					dnaSeq.Child = nil
					i += dnaSeqYin.Elements
				} else {
					clone := dnaSeqYin.Clone()
					dnaSeq = dnaSeq.Merge(clone.Sequence)
					i += clone.Elements
				}
			} else if i < dnaSeqYang.Index && dnaSeqYang.Index < dnaSeqYin.Index {
				if i == -1 {
					i = 0
				}
				if dnaSeq == nil {
					dnaSeq = dnaSeqYang.Clone()
					dnaSeq.Child = nil
					i += dnaSeq.Elements
				} else {
					clone := dnaSeqYang.Clone()
					dnaSeq = dnaSeq.Merge(clone.Sequence)
					i += clone.Elements
				}
			}
			if dnaSeqYin.Child == nil && dnaSeqYang.Child == nil {
				break
			}
			if dnaSeqYin.Child != nil {
				dnaSeqYin = dnaSeqYin.Child
			} else {
				dnaSeqYin.Index = math.MaxInt64
			}
			if dnaSeqYang.Child != nil {
				dnaSeqYang = dnaSeqYang.Child
			} else {
				dnaSeqYang.Index = math.MaxInt64
			}
		}
	} else if headYin != nil {
		dnaSeqYin := headYin.Clone()
		i := 0
		for {
			if i < dnaSeqYin.Index {
				if dnaSeq == nil {
					dnaSeq = dnaSeqYin.Clone()
					dnaSeq.Child = nil
					i += dnaSeq.Index + dnaSeq.Elements
				} else {
					clone := dnaSeqYin.Clone()
					dnaSeq = dnaSeq.Merge(clone.Sequence)
					i += clone.Index + clone.Elements
				}
				if dnaSeqYin.Child != nil {
					dnaSeqYin = dnaSeqYin.Child
				}
			} else {
				if dnaSeqYin.Child != nil {
					dnaSeqYin = dnaSeqYin.Child
				} else {
					break
				}
			}
		}
	} else if headYang != nil {
		dnaSeqYang := headYang.Clone()
		i := 0
		for {
			if i < dnaSeqYang.Index {
				if dnaSeq == nil {
					dnaSeq = dnaSeqYang.Clone()
					dnaSeq.Child = nil
					i += dnaSeq.Index + dnaSeq.Elements
				} else {
					clone := dnaSeqYang.Clone()
					dnaSeq = dnaSeq.Merge(clone.Sequence)
					i += clone.Index + clone.Elements
				}
				if dnaSeqYang.Child != nil {
					dnaSeqYang = dnaSeqYang.Child
				}
			} else {
				if dnaSeqYang.Child != nil {
					dnaSeqYang = dnaSeqYang.Child
				} else {
					break
				}
			}
		}
	}
	return dnaSeq
}

func (d *DNA) MarshalGenes() ([]byte, error) {
	codexGigasYin := d.Unwind(d.StrandYin)
	codexGigasYang := d.Unwind(d.StrandYang)
	chanYin := d.Sequence(codexGigasYin)
	chanYang := d.Sequence(codexGigasYang)
	dnaSeq := d.SpliceSequence([2]chan *Sequence{
		chanYin,
		chanYang,
	})

	if dnaSeq != nil {
		return dnaSeq.Bytes(), nil
	} else {
		return []byte{}, nil
	}
}

func (d *DNA) MarshalHelix() ([]byte, error) {
	return []byte(util.Hex(d.StrandYin) + "|" + util.Hex(d.StrandYang)), nil
}
