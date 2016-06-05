package dna

import (
	"fmt"
	"github.com/wmiller848/GoGP/gene"
	"github.com/wmiller848/GoGP/util"
	_ "sort"
)

type Sequence struct {
	Codex    Codex
	CodexID  int
	Index    int
	Elements int
}

func (s *Sequence) Node() *SequenceNode {
	return &SequenceNode{
		Sequence: s,
		Child:    nil,
	}
}

type SequenceNode struct {
	*Sequence
	Child *SequenceNode
}

func (s *SequenceNode) String() string {
	str := fmt.Sprintf("%+v ", *s.Sequence)
	if s.Child != nil {
		str += s.Child.String()
	}
	return str
}

func (s *SequenceNode) Bytes() []byte {
	bytes := []byte{}
	bytes = append(bytes, []byte(s.Codex.String())...)
	if s.Child != nil {
		bytes = append(bytes, s.Child.Bytes()...)
	}
	return bytes
}

func (s *SequenceNode) Clone() *SequenceNode {
	seq := &SequenceNode{
		Sequence: &Sequence{
			Codex:    s.Sequence.Codex,
			CodexID:  s.Sequence.CodexID,
			Index:    s.Sequence.Index,
			Elements: s.Sequence.Elements,
		},
	}
	if s.Child != nil {
		seq.Child = s.Child.Clone()
	}
	return seq
}

func (s *SequenceNode) Merge(seq *Sequence) *SequenceNode {
	node := seq.Node()
	if s.Index > seq.Index {
		node.Child = s
		return node
	} else {
		if s.Child == nil {
			s.Child = node
		} else {
			s.Child = s.Child.Merge(seq)
		}
		return s
	}
}

type DNA struct {
	StrandYing gene.GenericGene
	StrandYang gene.GenericGene
	Block      Block
}

func (d *DNA) Unwind(strand gene.GenericGene) CodexGigas {
	leng := len(strand)
	codexGigas := CodexGigas{}
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
						Index:    index,
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
	var headYing *SequenceNode
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
		if head == nil {
			break
		}

		if j == 0 {
			headYing = head
		} else {
			headYang = head
		}
	}

	var dnaSeq *SequenceNode
	if headYing != nil && headYang != nil {
		dnaSeqYing := headYing.Clone()
		dnaSeqYang := headYang.Clone()
		i := 0
		for {
			if i < dnaSeqYing.Index && dnaSeqYing.Index < dnaSeqYang.Index {
				if dnaSeq == nil {
					dnaSeq = dnaSeqYing.Clone()
					dnaSeq.Child = nil
					i += dnaSeq.Index + dnaSeq.Elements
				} else {
					clone := dnaSeqYang.Clone()
					dnaSeq = dnaSeq.Merge(clone.Sequence)
					i += clone.Index + clone.Elements
				}
				if dnaSeqYing.Child != nil {
					dnaSeqYing = dnaSeqYing.Child
				}
			} else if i < dnaSeqYang.Index && dnaSeqYang.Index < dnaSeqYing.Index {
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
				if dnaSeqYing.Child != nil {
					dnaSeqYing = dnaSeqYing.Child
				}
				if dnaSeqYang.Child != nil {
					dnaSeqYang = dnaSeqYang.Child
				}
				if dnaSeqYing.Child == nil && dnaSeqYang.Child == nil {
					break
				}
			}
		}
	} else if headYing != nil {
		dnaSeqYing := headYing.Clone()
		i := 0
		for {
			if i < dnaSeqYing.Index {
				if dnaSeq == nil {
					dnaSeq = dnaSeqYing.Clone()
					dnaSeq.Child = nil
					i += dnaSeq.Index + dnaSeq.Elements
				} else {
					clone := dnaSeqYing.Clone()
					dnaSeq = dnaSeq.Merge(clone.Sequence)
					i += clone.Index + clone.Elements
				}
				if dnaSeqYing.Child != nil {
					dnaSeqYing = dnaSeqYing.Child
				}
			} else {
				if dnaSeqYing.Child != nil {
					dnaSeqYing = dnaSeqYing.Child
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
	codexGigasYing := d.Unwind(d.StrandYing)
	codexGigasYang := d.Unwind(d.StrandYang)
	chanYing := d.Sequence(codexGigasYing)
	chanYang := d.Sequence(codexGigasYang)
	dnaSeq := d.SpliceSequence([2]chan *Sequence{
		chanYing,
		chanYang,
	})

	if dnaSeq != nil {
		return dnaSeq.Bytes(), nil
	} else {
		return []byte{}, nil
	}
}

func (d *DNA) MarshalHelix() ([]byte, error) {
	return []byte(util.Hex(d.StrandYing) + "|" + util.Hex(d.StrandYang)), nil
}
