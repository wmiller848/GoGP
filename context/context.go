package context

import (
	"fmt"
	"io"
	"math"
	"regexp"
	"sort"
	"time"

	"github.com/wmiller848/GoGP/data"
	"github.com/wmiller848/GoGP/dna"
	"github.com/wmiller848/GoGP/gene"
	"github.com/wmiller848/GoGP/program"
	"github.com/wmiller848/GoGP/util"
)

//type ScoreFunction func(int) int

var word_regex = regexp.MustCompile(`[a-zA-Z\_\-]+`)

type Context struct {
	Population int
	Programs   Programs
	visualMode bool
	terminal   *Terminal
}

func New() *Context {
	return &Context{}
}

func (c *Context) NewTerminal() {
	c.visualMode = true
	c.terminal = &Terminal{}
	c.terminal.Start(c)
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

func (c *Context) RunWithInlineScore(pipe io.Reader, score float64, inputs, population, generations int, auto bool) (string, *ProgramInstance) {
	uuid := util.RandomHex(32)
	c.InitPopulation(inputs, population)
	var i int = 0
	time.Sleep(500 * time.Millisecond)
	fountain := Multiplex(pipe)
	for {
		if i >= generations && !auto {
			break
		}

		parents := c.EvalInline(fountain, i, inputs, uuid)

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
			if c.visualMode {
				gns, _ := prgm.DNA.MarshalGenes()
				mathGns := gene.MathGene(gns).Heal()
				tree, _ := mathGns.MarshalTree()
				exp, _ := tree.MarshalExpression()
				str := fmt.Sprintf("Total Score: %3.2f\nGeneration: %v Expression: %v\n", (1.0-prgm.Score)*100.0, i, string(exp))
				str += "DNA:\n"
				str += "Yin:\n"
				codexGigasYin := prgm.DNA.Unwind(prgm.DNA.StrandYin)
				for i, _ := range codexGigasYin {
					str += fmt.Sprintf("  %v => %v\n", i, codexGigasYin[i])
				}
				str += "Yang:\n"
				codexGigasYang := prgm.DNA.Unwind(prgm.DNA.StrandYang)
				for i, _ := range codexGigasYang {
					str += fmt.Sprintf("  %v => %v\n", i, codexGigasYang[i])
				}
				chanYin := prgm.DNA.Sequence(codexGigasYin)
				chanYang := prgm.DNA.Sequence(codexGigasYang)
				dnaSeq := prgm.DNA.SpliceSequence([2]chan *dna.Sequence{
					chanYin,
					chanYang,
				})
				str += fmt.Sprintf("Sequence => %v\n", dnaSeq)
				str += "\nSub Scores:\n"
				for k, grp := range prgm.Group {
					c := float64(grp.Wrong) / float64(grp.Count)
					str += fmt.Sprintf("  %v (%v): %3.2f (%v / %v)\n", k, data.NumberFromString(k, inputs), (1.0-c)*100.00, grp.Count-grp.Wrong, grp.Count)
				}
				c.terminal.window.value = []byte(str)
			}
			if prgm != nil && (1.0-prgm.Score) > score {
				t := 0
				for _, grp := range prgm.Group {
					c := float64(grp.Wrong) / float64(grp.Count)
					if 1.0-c >= score {
						t++
					}
				}
				if t == len(prgm.Group) && t != 0 {
					// c.EvalInline(fountain, i, inputs, uuid, true)
					break
				}
			}
		}
		i++
	}
	fountain.Destroy()
	return uuid, c.Fitest()
}

func (c *Context) EvalInline(fountain *Multiplexer, generation, inputs int, uuid string) Programs {
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
	testData, threshold, inputMap, assertMap := data.New(buffer, inputs)

	for i, _ := range c.Programs {
		prgm := c.Programs[i]
		gns, _ := prgm.DNA.MarshalGenes()
		mathGns := gene.MathGene(gns).Heal()
		tree, _ := mathGns.MarshalTree()
		if tree == nil {
			continue
		}
		wrong := make(map[string]*Group)
		for _, dat := range testData {
			matched := word_regex.Match([]byte(dat.AssertStr))
			if matched {
				if wrong[dat.AssertStr] == nil {
					wrong[dat.AssertStr] = &Group{
						Count: 0,
						Wrong: 0,
					}
				}
			} else {
				str := word_regex.String()
				if wrong[str] == nil {
					wrong[str] = &Group{
						Count: 0,
						Wrong: 0,
					}
				}
			}

			out := tree.Eval(dat.Input...)
			diff := math.Abs(out - dat.Assert)
			if matched {
				wrong[dat.AssertStr].Count++
				if diff >= threshold || math.IsNaN(out) {
					wrong[dat.AssertStr].Wrong++
				}
			} else {
				str := word_regex.String()
				wrong[str].Count++
				if diff >= 1 || math.IsNaN(out) {
					wrong[str].Wrong++
				}
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
		prgm.InputMap = inputMap
		prgm.AssertMap = assertMap
		// if log {
		// 	return c.Programs
		// }
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

// func (c *Context) ScoreProgram() {
//
// }
