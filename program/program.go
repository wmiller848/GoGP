package program

import (
	"github.com/wmiller848/GoGP/gene"
	_ "github.com/wmiller848/GoGP/util"
	"io/ioutil"
	"strings"
)

type Program struct {
	Gene     gene.Gene
	Template string
}

func New(varCount, knobCount int, block gene.BuildingBlock) *Program {
	g, varsTmpl := block.Random(varCount, knobCount)
	tplBytes, _ := ioutil.ReadFile("./program/main.coffee")
	return &Program{
		Gene:     g,
		Template: string(tplBytes),
	}
}

func (p *Program) MarshalProgram() ([]byte, error) {
	pgm := p.Template
	root, _ := p.Gene.MarshalTree()
	exp, _ := root.MarshalExpression()
	pgm = strings.Replace(pgm, "{{dna}}", string(p.Gene.Clone()), 1)
	//=====
	// Vars
	//=====
	pgm = strings.Replace(pgm, "{{vars}}", "", 1)
	//=====
	// Spawn
	//=====
	pgm = strings.Replace(pgm, "{{spawn}}", "", 1)
	//=====
	// Alive
	//=====
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
