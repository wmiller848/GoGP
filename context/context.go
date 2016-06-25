package context

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/wmiller848/GoGP/data"
	"github.com/wmiller848/GoGP/gene"
	"github.com/wmiller848/GoGP/program"
	"github.com/wmiller848/GoGP/util"
)

//type ScoreFunction func(int) int

type Context struct {
	Population  int
	Programs    Programs
	VerboseMode bool
}

func New() *Context {
	return &Context{}
}

func (c *Context) Verbose() bool {
	c.VerboseMode = !c.VerboseMode
	return c.VerboseMode
}
func (c *Context) Fitest() *ProgramInstance {
	if len(c.Programs) > 0 {
		sort.Sort(c.Programs)
		return c.Programs[0]
	} else {
		return nil
	}
}

func (c *Context) InitPopulation(inputs, population int) {
	c.Population = population
	c.Programs = make(Programs, population)
	var i int
	for i = 0; i < population; i++ {
		pgm := &ProgramInstance{
			Program:    program.New(inputs),
			ID:         util.RandomHex(16),
			Generation: 0,
			Score:      math.MaxFloat64,
		}
		c.Programs[i] = pgm
	}
}

func (c *Context) RunWithInlineScore(pipe io.Reader, threshold, score float64, inputs, population, generations int, auto bool) (string, *ProgramInstance) {
	uuid := util.RandomHex(32)
	c.InitPopulation(inputs, population)
	var i int = 0
	time.Sleep(500 * time.Millisecond)
	fountain := Multiplex(pipe)
	max := 0
	for {
		if i >= generations && !auto {
			break
		}

		parents := c.EvalInline(fountain, i, inputs, threshold, uuid)

		children := []*ProgramInstance{}
		if len(parents) > 0 && i != generations-1 {
			for i := 0; i < c.Population-len(parents); i++ {
				pgm := &ProgramInstance{
					Program:    parents[i%len(parents)].Mutate(),
					ID:         util.RandomHex(16),
					Generation: i + 1,
					Score:      math.MaxFloat64,
				}
				children = append(children, pgm)
			}
			c.Programs = append(parents, children...)
			prgm := c.Fitest()
			if c.VerboseMode {
				gns, _ := prgm.DNA.MarshalGenes()
				mathGns := gene.MathGene(gns).Heal()
				tree, _ := mathGns.MarshalTree()
				exp, _ := tree.MarshalExpression()
				str := fmt.Sprintf("\rTotal Score: %3.2f Generation: %v Expression: %v", (1.0-prgm.Score)*100.0, i, string(exp))
				strByts := []byte(str)
				if len(strByts) > max {
					max = len(strByts)
				} else {
					pad := make([]byte, max-len(strByts))
					for j := 0; j < len(pad); j++ {
						pad[j] = byte(' ')
					}
					strByts = append(strByts, pad...)
				}
				str = string(strByts)
				fmt.Printf(str)
			}
			if prgm != nil && (1.0-prgm.Score) > score {
				t := 0
				for _, grp := range prgm.Group {
					c := float64(grp.Wrong) / float64(grp.Count)
					if 1.0-c > score {
						t++
					}
				}
				if t == len(prgm.Group) && t != 0 {
					break
				}
			}
		}
		i++
	}
	fountain.Destroy()
	if c.VerboseMode {
		fmt.Printf("\n")
	}
	return uuid, c.Fitest()
}

func (c *Context) EvalInline(fountain *Multiplexer, generation, inputs int, threshold float64, uuid string) Programs {
	validPrograms := 0
	tap := fountain.Multiplex().Tap()
	var buffer []byte
	for {
		d, open := <-tap
		if open == false {
			break
		}
		buffer = append(buffer, d...)
	}
	lines := bytes.Split(buffer, []byte("\n"))
	testData := []*data.TestData{}
	for i, _ := range lines {
		if len(lines[i]) > 0 {
			nums := bytes.Split(lines[i], []byte(" "))
			if len(nums) >= inputs {
				dat := &data.TestData{}
				for j, numByts := range nums {
					num, err := strconv.ParseFloat(string(numByts), 64)
					if err == nil {
						if j < inputs {
							dat.Input = append(dat.Input, num)
						} else {
							dat.Assert = num
						}
					}
				}
				testData = append(testData, dat)
			}
		}
	}
	for i, _ := range c.Programs {
		prgm := c.Programs[i]
		gns, _ := prgm.DNA.MarshalGenes()
		mathGns := gene.MathGene(gns).Heal()
		tree, _ := mathGns.MarshalTree()
		if tree == nil {
			continue
		}
		wrong := make(map[float64]*Group)
		for _, dat := range testData {
			if wrong[dat.Assert] == nil {
				wrong[dat.Assert] = &Group{
					Count: 0,
					Wrong: 0,
				}
			}
			out := tree.Eval(dat.Input...)
			diff := math.Abs(out - dat.Assert)
			//fmt.Println(prgm.ID, inputFloats, out, assertFloat, diff)
			wrong[dat.Assert].Count++
			if diff >= threshold || math.IsNaN(out) {
				wrong[dat.Assert].Wrong++
			}
		}
		total := 0.0
		for _, grp := range wrong {
			c := float64(grp.Wrong) / float64(grp.Count)
			total += c
		}
		total /= float64(len(wrong))
		prgm.Score = total
		prgm.Group = wrong
		validPrograms++
	}

	sort.Sort(c.Programs)
	// Top 30%
	limit := validPrograms / 3
	// Extra random newbies we throw in
	variance := limit / 3
	parents := make(Programs, limit+variance)
	for i := 0; i < limit; i++ {
		parents[i] = c.Programs[i]
	}
	for i := limit; i < limit+variance; i++ {
		pgm := &ProgramInstance{
			Program:    program.New(inputs),
			ID:         util.RandomHex(16),
			Generation: generation,
			Score:      math.MaxFloat64,
		}
		parents[i] = pgm
	}
	return parents
}
func (c *Context) InlineData() []*data.TestData {
	return []*data.TestData{}
}

func (c *Context) ScoreProgram() {

}
