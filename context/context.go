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
	ID string
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
	Population int
	Programs   []*ProgramInstance
}

func New() *Context {
	return &Context{}
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

func (c *Context) RunWithInlineScore(pipe io.Reader, inputs, population, generations int) *program.Program {
	os.Mkdir("./out", 0777)
	sha := util.Random(32)
	os.Mkdir("./out/generations", 0777)
	os.RemoveAll("./out/generations/" + util.Hex(sha))
	os.Mkdir("./out/generations/"+util.Hex(sha), 0777)
	c.InitPopulation(inputs, population)
	var i int
	time.Sleep(500 * time.Millisecond)
	fountain := Multiplex(pipe)
	for i = 0; i < generations; i++ {
		c.EvalInline(fountain, i, inputs, sha)
	}
	fountain.Destroy()
	return c.Fitest()
}

func (c *Context) EvalInline(fountain *Multiplexer, generation, inputs int, uuid []byte) {
	path := "./out/generations/" + util.Hex(uuid) + "/" + strconv.Itoa(generation)
	os.Mkdir(path, 0777)

	//		* Each testBuf row ->
	//			* compute average score
	scores := Scores{}
	genScore := 0.0
	for i, _ := range c.Programs {
		prgm := c.Programs[i]
		cmdStr := path + "/" + prgm.ID
		cmd := exec.Command("coffee", cmdStr)
		//
		cmd.Stderr = os.Stderr
		stdoutBuffer := NewBuffer()
		cmd.Stdout = stdoutBuffer
		//
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
			//for j, _ := range lines {
			//if i != j && string(lines[i]) == string(lines[j]) {
			////fmt.Println("FOUND SAME!", i, j, string(lines[i]), string(lines[j]))
			//}
			//}
			//fmt.Printf("%v", string(lines[i])+"\n")
			if len(lines[i]) > 0 {
				nums := bytes.Split(lines[i], []byte(" "))
				//fmt.Println(nums)
				if len(nums) >= inputs {
					num, err := strconv.ParseFloat(string(nums[inputs]), 64)
					if err == nil {
						assert = append(assert, num)
					}
				}
			}
		}
		//
		//
		prgmBytes, _ := prgm.MarshalProgram()
		fmt.Println(generation, i, "Command - '"+cmdStr+"'")
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
			fmt.Println(err.Error())
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
		fmt.Println(len(assert), len(output))
		if len(assert) == len(output) && len(assert) > 0 {
			avgScore := 0.0
			for i, _ := range assert {
				diff := output[i] - assert[i]
				avgScore += diff
			}
			score := ProgramScore{
				Index: i,
				Score: math.Abs(avgScore / float64(len(assert))),
			}
			genScore += score.Score
			scores = append(scores, score)
			fmt.Println("Score - ", score.Score)
		} else {
			fmt.Println("Program provided incorrect amount of outputs, terminating DNA")
		}
	}
	genScore /= float64(len(scores))
	fmt.Println("Generation Score -", genScore)

	sort.Sort(scores)
	// Top 30%
	limit := len(scores) / 3
	parents := []*ProgramInstance{}
	children := []*ProgramInstance{}
	for i, _ := range scores {
		if i < limit {
			parents = append(parents, c.Programs[scores[i].Index])
		}
	}

	if len(parents) > 0 {
		for i := 0; i < c.Population-len(parents); i++ {
			pgm := &ProgramInstance{
				Program: parents[i%len(parents)].Mutate(),
				ID:      util.RandomHex(16),
			}
			children = append(children, pgm)
		}
	}
	c.Programs = append(parents, children...)
}

func (c *Context) Fitest() *program.Program {
	return nil
}

func (c *Context) InitPopulation(inputs, population int) {
	c.Population = population
	c.Programs = make([]*ProgramInstance, population)
	var i int
	for i = 0; i < population; i++ {
		pgm := &ProgramInstance{
			Program: program.New(inputs),
			ID:      util.RandomHex(16),
		}
		c.Programs[i] = pgm
	}
}
