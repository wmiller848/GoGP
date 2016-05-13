package context

import (
	"github.com/wmiller848/GoGP/gene"
	"github.com/wmiller848/GoGP/program"
)

type ScoreFunction func(int) int

type Context struct {
	Programs []*program.Program
}

func New() *Context {
	return &Context{}
}

func (c *Context) RunWithScoreFunc(inputs, population, generations int, scoreFunc ScoreFunction) *program.Program {
	var i int
	c.Programs = make([]*program.Program, population)
	for i = 0; i < population; i++ {
		pgm := program.New(inputs, inputs*4, &gene.MathBuildingBlock{})
		c.Programs[i] = pgm
	}

	for i = 0; i < generations; i++ {
		c.EvalFunc(scoreFunc)
	}

	return c.Fitest()
}

func (c *Context) EvalFunc(scoreFunc ScoreFunction) {

}

func (c *Context) RunWithInlineScore(buf []byte, inputs, population, generations int) *program.Program {
	var i int
	c.Programs = make([]*program.Program, population)
	for i = 0; i < population; i++ {
		pgm := program.New(inputs, inputs*4, &gene.MathBuildingBlock{})
		c.Programs[i] = pgm
	}

	for i = 0; i < generations; i++ {
		c.EvalInline()
	}

	return c.Fitest()
}

func (c *Context) EvalInline() {

}

func (c *Context) Fitest() *program.Program {
	return nil
}
