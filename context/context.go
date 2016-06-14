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

	_ "github.com/wmiller848/GoGP/gene"
	"github.com/wmiller848/GoGP/program"
	"github.com/wmiller848/GoGP/util"
)

type ScoreFunction func(int) int

type ProgramInstance struct {
	*program.Program
	ID         string
	Generation int
}

type ProgramScore struct {
	Index int
	Score float64
}

type Scores []ProgramScore

func (s Scores) Len() int           { return len(s) }
func (s Scores) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Scores) Less(i, j int) bool { return s[i].Score < s[j].Score }

type Context struct {
	Population  int
	Programs    []*ProgramInstance
	VerboseMode bool
}

func New() *Context {
	return &Context{}
}

func (c *Context) Verbose() {
	c.VerboseMode = !c.VerboseMode
}

//func (c *Context) RunWithScoreFunc(inputs, population, generations int, scoreFunc ScoreFunction) *program.Program {
//var i int
//c.InitPopulation(inputs, population)
//for i = 0; i < generations; i++ {
//c.EvalFunc(scoreFunc)
//}
//return c.Fitest()
//}

//func (c *Context) EvalFunc(scoreFunc ScoreFunction) {
//}

func (c *Context) RunWithInlineScore(pipe io.Reader, inputs, population, generations int) (string, *ProgramInstance) {
	os.Mkdir("./out", 0777)
	uuid := util.RandomHex(32)
	os.Mkdir("./out/generations", 0777)
	os.RemoveAll("./out/generations/" + uuid)
	os.Mkdir("./out/generations/"+uuid, 0777)
	c.InitPopulation(inputs, population)
	var i int
	time.Sleep(500 * time.Millisecond)
	fountain := Multiplex(pipe)
	for i = 0; i < generations; i++ {
		c.EvalInline(fountain, i, inputs, uuid)
	}
	fountain.Destroy()
	if c.VerboseMode {
		fmt.Printf("\n")
	}
	return uuid, c.Fitest()
}

func (c *Context) EvalInline(fountain *Multiplexer, generation, inputs int, uuid string) {
	path := "./out/generations/" + uuid + "/" + strconv.Itoa(generation)
	os.Mkdir(path, 0777)

	//		* Each testBuf row ->
	//			* compute average score
	scores := Scores{}
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
			mean := 0.0
			for i, _ := range assert {
				diff := assert[i] - output[i]
				//fmt.Println("Output", assert[i], output[i], diff)
				mean += diff
			}

			e := 0.0
			for i, _ := range assert {
				q := math.Pow(assert[i]-mean, 2)
				e += q
			}

			score := ProgramScore{
				Index: i,
				Score: e / float64(len(assert)),
			}
			if !math.IsNaN(score.Score) && !math.IsInf(score.Score, 0) {
				scores = append(scores, score)
				//fmt.Println("Average Score - ", score.Score)
			} else {
				//fmt.Println("Program provided invalid score, terminating DNA")
			}
		} else {
			//fmt.Println("Program provided incorrect amount of outputs, terminating DNA")
		}
	}

	sort.Sort(scores)
	// Top 25%
	limit := len(scores) / 4
	parents := []*ProgramInstance{}
	children := []*ProgramInstance{}
	for i, _ := range scores {
		//fmt.Println(scores[i].Score)
		if i < limit {
			parents = append(parents, c.Programs[scores[i].Index])
		}
	}

	if len(parents) > 0 {
		for i := 0; i < c.Population-len(parents); i++ {
			pgm := &ProgramInstance{
				Program:    parents[i%len(parents)].Mutate(),
				ID:         util.RandomHex(16),
				Generation: generation + 1,
			}
			children = append(children, pgm)
		}
	}
	c.Programs = append(parents, children...)
	if c.VerboseMode {
		fmt.Printf(".")
	}
}

func (c *Context) Fitest() *ProgramInstance {
	return c.Programs[0]
}

func (c *Context) InitPopulation(inputs, population int) {
	c.Population = population
	c.Programs = make([]*ProgramInstance, population)
	var i int
	for i = 0; i < population; i++ {
		pgm := &ProgramInstance{
			Program:    program.New(inputs),
			ID:         util.RandomHex(16),
			Generation: 0,
		}
		c.Programs[i] = pgm
	}
}
