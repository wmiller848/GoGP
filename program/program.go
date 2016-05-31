package program

import (
	"io/ioutil"
	"strings"

	"github.com/wmiller848/GoGP/dna"
	"github.com/wmiller848/GoGP/gene"
	_ "github.com/wmiller848/GoGP/util"
)

type Program struct {
	Block    dna.Block
	DNA      *dna.DNA
	Template string
}

//func New(varCount, knobCount int, block gene.BuildingBlock) *Program {
//g := block.Random(varCount, knobCount)
//tplBytes, _ := ioutil.ReadFile("./program/main.coffee")
//return &Program{
//Block:    block,
//Gene:     g,
//Template: string(tplBytes),
//}
//}

func New(count int) *Program {
	bases := [4]dna.Base{0x00, 0x40, 0x80, 0xc0}
	codons := []dna.Codon{
		dna.Codon(dna.CodonStop), dna.Codon("+"), dna.Codon("-"),
		dna.Codon("*"), dna.Codon("/"), dna.Codon("#"),
		dna.Codon("$a"), dna.Codon("$b"), dna.Codon("$c"),
		dna.Codon("$d"),
	}
	blk, _ := dna.NewBlock4x3(bases, codons)
	d := blk.Random()
	tplBytes, _ := ioutil.ReadFile("./program/main.coffee")
	return &Program{
		Block:    blk,
		DNA:      d,
		Template: string(tplBytes),
	}
}

func (p *Program) Mate(mate *Program) *Program {
	return New(4)
}

func (p *Program) MarshalProgram() ([]byte, error) {
	pgm := p.Template
	gns, _ := p.DNA.MarshalGenes()

	//=====
	// Coffee Path
	//=====
	pgm = strings.Replace(pgm, "{{coffee_path}}", "/usr/local/bin/coffee", 1)
	//=====
	// DNA
	//=====
	helix, _ := p.DNA.MarshalHelix()
	pgm = strings.Replace(pgm, "{{dna}}", string(helix), 1)
	//=====
	// Spawn
	//=====
	pgm = strings.Replace(pgm, "{{spawn}}", "", 1)

	for i, _ := range gns {
		//=====
		// Vars
		//=====
		pgm = strings.Replace(pgm, "{{vars}}", gene.VarsTemplate(gns[i])+"\n{{vars}}", 1)
		root, _ := gns[i].MarshalTree()
		exp, _ := root.MarshalExpression()
		//=====
		// Alive
		//=====
		pgm = strings.Replace(pgm, "{{alive}}", string(exp)+"\n{{alive}}", 1)
	}
	pgm = strings.Replace(pgm, "{{vars}}", "", -1)
	pgm = strings.Replace(pgm, "{{alive}}", "", -1)
	//=====
	// Dieing
	//=====
	pgm = strings.Replace(pgm, "{{dieing}}", "", 1)
	//=====
	// Dead
	//=====
	pgm = strings.Replace(pgm, "{{dead}}", "", 1)
	return []byte(pgm), nil
}
