package context

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/wmiller848/GoGP/gene"
	"github.com/wmiller848/GoGP/program"
	"github.com/wmiller848/GoGP/util"
)

type ScoreFunction func(int) int

type ProgramInstance struct {
	*program.Program
	ID         string
	Generation int
	Score      float64
}

type Programs []*ProgramInstance

func (p Programs) Len() int           { return len(p) }
func (p Programs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Programs) Less(i, j int) bool { return p[i].Score < p[j].Score }

type Context struct {
	Population  int
	Programs    Programs
	VerboseMode bool
}

func New() *Context {
	return &Context{}
}

func (c *Context) Verbose() {
	c.VerboseMode = !c.VerboseMode
}

func (c *Context) RunWithInlineScore(pipe io.Reader, threshold float64, inputs, population, generations int, auto bool) (string, *ProgramInstance) {
	//os.Mkdir("./out", 0777)
	uuid := util.RandomHex(32)
	//os.Mkdir("./out/generations", 0777)
	//os.RemoveAll("./out/generations/" + uuid)
	//os.Mkdir("./out/generations/"+uuid, 0777)
	c.InitPopulation(inputs, population)
	var i int = 0
	time.Sleep(500 * time.Millisecond)
	fountain := Multiplex(pipe)
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
				//fmt.Printf(".")
				fmt.Printf("\rScore - %3.2f Generation %v", (1.0-prgm.Score)*100.0, i)
			}
			if prgm != nil && (1.0-prgm.Score) > 0.95 {
				break
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
	path := "./out/generations/" + uuid + "/" + strconv.Itoa(generation)
	os.Mkdir(path, 0777)

	//		* Each testBuf row ->
	//			* compute average score
	validPrograms := 0
	tap := fountain.Multiplex().Tap()
	var data []byte
	for {
		d, open := <-tap
		if open == false {
			break
		}
		data = append(data, d...)
	}
	for i, _ := range c.Programs {
		prgm := c.Programs[i]
		gns, _ := prgm.DNA.MarshalGenes()
		mathGns := gene.MathGene(gns).Heal()
		tree, _ := mathGns.MarshalTree()
		if tree == nil {
			continue
		}
		wrong := 0
		lines := bytes.Split(data, []byte("\n"))
		for i, _ := range lines {
			if len(lines[i]) > 0 {
				nums := bytes.Split(lines[i], []byte(" "))
				if len(nums) >= inputs {
					inputFloats := []float64{}
					assertFloat := math.NaN()
					for j, numByts := range nums {
						num, err := strconv.ParseFloat(string(numByts), 64)
						if err == nil {
							if j < inputs {
								inputFloats = append(inputFloats, num)
							} else {
								assertFloat = num
							}
						}
					}
					out := tree.Eval(inputFloats...)
					diff := math.Abs(out - assertFloat)
					//fmt.Println(prgm.ID, inputFloats, out, assertFloat, diff)
					if diff >= threshold || math.IsNaN(out) {
						wrong++
					}
				}
			}
		}
		prgm.Score = float64(wrong) / float64(len(lines))
		validPrograms++
	}

	sort.Sort(c.Programs)
	// Top 30%
	limit := validPrograms / 3
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

func (c *Context) Fitest() *ProgramInstance {
	//if c.VerboseMode {
	//for i, _ := range c.Programs {
	//gn, err := c.Programs[i].DNA.MarshalGenes()
	//if err != nil {
	//fmt.Println(err.Error())
	//continue
	//}
	//d := gene.MathGene(gn)
	//fmt.Println("Program", c.Programs[i].ID, c.Programs[i].Score, string(d.Heal()))
	//}
	//}
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
