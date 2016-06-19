package program

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/wmiller848/GoGP/dna"
	"github.com/wmiller848/GoGP/gene"
)

type Program struct {
	InputCount int
	Block      dna.Block
	DNA        *dna.DNA
	Template   string
}

func New(count int) *Program {
	bases := [4]dna.Base{0x00, 0x40, 0x80, 0xc0}
	codons := []dna.Codon{
		dna.Codon("&"), dna.Codon("|"), dna.Codon("^"),
		dna.Codon("+"), dna.Codon("-"), dna.Codon("*"), dna.Codon("/"),
		dna.Codon("0"), dna.Codon("1"), dna.Codon("2"), dna.Codon("3"),
		dna.Codon("4"), dna.Codon("5"), dna.Codon("6"), dna.Codon("7"),
		dna.Codon("8"), dna.Codon("9"),
		dna.Codon(","), dna.Codon("{"), dna.Codon("}"),
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
	mathGns := gene.MathGene(gns).Heal()

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
	// =====
	// Variabls
	// =====
	pgm = strings.Replace(pgm, "{{vars}}", gene.VariableTemplate(p.InputCount), 1)
	//=====
	// Output
	//=====
	root, _ := mathGns.MarshalTree()
	exp, _ := root.MarshalExpression()
	pgm = strings.Replace(pgm, "{{output}}", string(exp), 1)
	return []byte(pgm), nil
}
