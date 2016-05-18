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

func (c *Context) InitPopulation(inputs, population int) {
	var i int
	c.Programs = make([]*program.Program, population)
	for i = 0; i < population; i++ {
		pgm := program.New(inputs, inputs*4, &gene.MathBuildingBlock{})
		c.Programs[i] = pgm
	}
}

func (c *Context) RunWithInlineScore(traingBuf []byte, inputs, population, generations int) *program.Program {
	c.InitPopulation(inputs, population)
	var i int
	for i = 0; i < generations; i++ {
		c.EvalInline(traingBuf)
	}

	return c.Fitest()
}

func (c *Context) EvalInline(traingBuf []byte) {
	//	Each program in population ->
	//		* Apply 'life' ->
	//			* Drain Energy

	//		* Each testBuf row ->
	//			* compute average score

	//	Each in top 30% ->
	//		* Add Energy
	//		* Cross with other top 30%

	//	Each program in population ->
	//		* If energy <= 0
	//			* dead
	//	Get population - dead
}

func (c *Context) Fitest() *program.Program {
	return nil
}
