package context

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"time"

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

func (c *Context) RunWithInlineScore(pipe io.Reader, inputs, population, generations int, auto bool) (string, *ProgramInstance) {
	os.Mkdir("./out", 0777)
	uuid := util.RandomHex(32)
	os.Mkdir("./out/generations", 0777)
	os.RemoveAll("./out/generations/" + uuid)
	os.Mkdir("./out/generations/"+uuid, 0777)
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
			if prgm != nil && prgm.Score < 0.1 {
				break
			}
		}
		if c.VerboseMode {
			fmt.Printf(".")
		}
		i++
	}
	fountain.Destroy()
	if c.VerboseMode {
		fmt.Printf("\n")
	}
	return uuid, c.Fitest()
}

func (c *Context) EvalInline(fountain *Multiplexer, generation, inputs int, uuid string) Programs {
	path := "./out/generations/" + uuid + "/" + strconv.Itoa(generation)
	os.Mkdir(path, 0777)

	//		* Each testBuf row ->
	//			* compute average score
	validPrograms := 0
	for i, _ := range c.Programs {
		prgm := c.Programs[i]
		cmdStr := path + "/" + prgm.ID
		cmd := exec.Command("coffee", cmdStr)
		//
		stderrBuffer := NewBuffer()
		cmd.Stderr = stderrBuffer
		//
		stdoutBuffer := NewBuffer()
		cmd.Stdout = stdoutBuffer
		// Parse out the asserted correct value from the data stream
		stdinBuffer := NewBuffer()
		var stdinTap chan []byte
		pipe := fountain.Multiplex()
		cmd.Stdin, stdinTap = stdinBuffer.Pipe(pipe)
		var data []byte
		var assert []float64
		// TODO :: get io streaming to work
		for {
			d, open := <-stdinTap
			if open == false {
				break
			}
			data = append(data, d...)
		}
		lines := bytes.Split(data, []byte("\n"))
		for i, _ := range lines {
			if len(lines[i]) > 0 {
				nums := bytes.Split(lines[i], []byte(" "))
				if len(nums) >= inputs {
					num, err := strconv.ParseFloat(string(nums[inputs]), 64)
					if err == nil {
						assert = append(assert, num)
					}
				}
			}
		}
		//
		prgmBytes, _ := prgm.MarshalProgram()
		//fmt.Println("Command - '" + cmdStr + "'")
		err := ioutil.WriteFile(path+"/"+prgm.ID, prgmBytes, 0555)
		if err != nil {
			fmt.Println(err.Error())
		}
		err = cmd.Start()
		if err != nil {
			fmt.Println(err.Error())
		}
		err = cmd.Wait()
		if err != nil {
			//fmt.Println(err.Error())
		}

		stdoutTap := stdoutBuffer.Tap()
		stdoutBuffer.Close()
		data = []byte{}
		output := []float64{}
		for {
			d, open := <-stdoutTap
			if open == false {
				break
			}
			data = append(data, d...)
		}
		lines = bytes.Split(data, []byte("\n"))
		for i, _ := range lines {
			num, err := strconv.ParseFloat(string(lines[i]), 64)
			if err == nil {
				output = append(output, num)
			} else {
				data = append(data, lines[i]...)
			}
		}
		// Compair output to assert
		if len(assert) == len(output) && len(assert) > 0 {
			score := 0.0
			for i, _ := range assert {
				diff := math.Abs(assert[i] - output[i])
				if diff > 500 || math.IsNaN(output[i]) {
					score++
				}
			}
			prgm.Score = score / float64(len(assert))
			validPrograms++
		}
	}

	sort.Sort(c.Programs)
	// Top 30%
	limit := validPrograms / 3
	parents := make(Programs, limit)
	for i := 0; i < limit; i++ {
		parents[i] = c.Programs[i]
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
