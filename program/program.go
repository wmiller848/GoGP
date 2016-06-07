package program

import (
	"errors"
	_ "fmt"
	"io/ioutil"
	"strings"

	"github.com/wmiller848/GoGP/dna"
	"github.com/wmiller848/GoGP/gene"
	_ "github.com/wmiller848/GoGP/util"
)

type Program struct {
	InputCount int
	Block      dna.Block
	DNA        *dna.DNA
	Template   string
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
		dna.Codon("+"), dna.Codon("-"),
		dna.Codon("*"), dna.Codon("/"), dna.Codon("0"),
		dna.Codon("1"), dna.Codon("2"), dna.Codon("3"),
		dna.Codon("4"), dna.Codon("5"), dna.Codon("6"),
		dna.Codon("7"), dna.Codon("8"), dna.Codon("9"),
		dna.Codon(","),
	}
	for i := 0; i < count; i++ {
		codons = append(codons, dna.Codon(gene.Variable(i)))
	}
	blk, _ := dna.NewBlock4x3(bases, codons)
	d := blk.Random()
	tplBytes, _ := ioutil.ReadFile("./program/main.coffee")
	return &Program{
		InputCount: count,
		Block:      blk,
		DNA:        d,
		Template:   string(tplBytes),
	}
}

func (p *Program) Mutate() *Program {
	pgm := New(p.InputCount)
	pgm.DNA = p.DNA.Mutate()
	return pgm
}

func (p *Program) MarshalProgram() ([]byte, error) {
	pgm := p.Template
	gns, _ := p.DNA.MarshalGenes()
	mathGns := gene.MathGene(gns)
	mathGns = mathGns.Heal()

	if len(mathGns) == 0 {
		return nil, errors.New("DNA contains no genes")
	}
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
	// =====
	// Vars
	// =====
	pgm = strings.Replace(pgm, "{{vars}}", gene.VariableTemplate(mathGns), 1)
	//=====
	// Alive
	//=====
	root, _ := mathGns.MarshalTree()
	exp, _ := root.MarshalExpression()
	pgm = strings.Replace(pgm, "{{alive}}", string(exp), 1)
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
