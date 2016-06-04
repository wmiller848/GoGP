package dna

import (
	"fmt"
	"github.com/wmiller848/GoGP/gene"
	"github.com/wmiller848/GoGP/util"
)

type SequenceContext struct {
}

type Sequence struct {
	Codex    Codex
	CodexID  int
	Index    int
	Elements int
}

type DNA struct {
	StrandYing gene.GenericGene
	StrandYang gene.GenericGene
	Block      Block
}

func (d *DNA) Unwind(strand gene.GenericGene) CodexGigas {
	// bases := d.Block.Bases()
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
			// fmt.Println(strand_frag, string(codon))
			codex = append(codex, codon)
		}
		codexGigas = append(codexGigas, codex)
	}
	return codexGigas
}

func (d *DNA) Sequence(codexGigas CodexGigas) chan *Sequence {
	chanSeq := make(chan *Sequence)
	go func() {
		for o, codex := range codexGigas {
			i := 0
			j := 0
			n := 0
			reading := false
			codexDecoded := Codex{}
			for _, codon := range codex {
				if string(codon) == string(CodonStart) {
					reading = true
					j = i + o
				} else if string(codon) != string(CodonStop) && reading == true {
					codexDecoded = append(codexDecoded, codon)
					n++
				} else if string(codon) == string(CodonStop) && reading == true {
					if len(codexDecoded) == 0 {
						reading = false
						continue
					}
					seq := Sequence{
						Codex:    codexDecoded,
						CodexID:  o,
						Index:    j,
						Elements: n,
					}
					chanSeq <- &seq
					codexDecoded = Codex{}
					reading = false
					n = 0
				}
				i++
			}
		}
		close(chanSeq)
	}()
	return chanSeq
}

func (d *DNA) SpliceSequence(ctx *SequenceContext, chanSeqs [2]chan *Sequence) {
	go func() {
		codexs := [2][]CodexGigas{
			[]CodexGigas{},
			[]CodexGigas{},
		}
		for i, chanSeq := range chanSeqs {
			seq := <-chanSeq
			ctx.Lock()
			ctx.CodexGigas[i][seq.CodexID] = append(codexGigas, seq.Codex)
			sort.Sort(ctx.CodexGigas[i][seq.CodexID])
			codexs[i] = append(codexs[i], ctx.CodexGigas[i][seq.CodexID][0])
			ctx.Unlock()
		}
	}()
}

func (d *DNA) MarshalGenes() (CodexGigas, error) {
	codexGigasYing := d.Unwind(d.StrandYing)
	codexGigasYang := d.Unwind(d.StrandYang)
	fmt.Println(codexGigasYing)
	fmt.Println(codexGigasYang)
	sequenceCtx := SequenceCtx{
		CodexGigas: CodexGigas{},
	}
	chanYing := d.Sequence(codexGigasYing)
	chanYang := d.Sequence(codexGigasYang)
	for {
		seq, open := <-chanYing
		if seq == nil || open == false {
			break
		}
		fmt.Println("Sequence Ying -", seq)
		codexGigas = append(codexGigas, seq.Codex)
	}
	for {
		seq, open := <-chanYang
		if seq == nil || open == false {
			break
		}
		fmt.Println("Sequence Yang", seq)
		codexGigas = append(codexGigas, seq.Codex)
	}

	return codexGigas, nil
}

func (d *DNA) MarshalHelix() ([]byte, error) {
	return []byte(util.Hex(d.StrandYing) + "|" + util.Hex(d.StrandYang)), nil
}
