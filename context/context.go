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

	"github.com/wmiller848/GoGP/gene"
	"github.com/wmiller848/GoGP/program"
	"github.com/wmiller848/GoGP/util"
)

type ScoreFunction func(int) int

type ProgramInstance struct {
	*program.Program
	Energy int
	Stage  int
	ID     string
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
	Programs []*ProgramInstance
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
	for i = 0; i < generations; i++ {
		c.EvalInline(pipe, i, inputs, sha)
	}
	return c.Fitest()
}

func (c *Context) EvalInline(pipe io.Reader, generation, inputs int, uuid []byte) {
	path := "./out/generations/" + util.Hex(uuid) + "/" + strconv.Itoa(generation)
	os.Mkdir(path, 0777)
	//	Each program in population ->
	//		* Apply 'life' ->
	//			* Drain Energy
	for i, _ := range c.Programs {
		prgm := c.Programs[i]
		prgm.Energy -= 10
	}

	//		* Each testBuf row ->
	//			* compute average score
	scores := Scores{}
	fountain := Multiplex(pipe, len(c.Programs))
	// HACK HACK waiting for IO
	time.Sleep(100 * time.Millisecond)
	for i, _ := range c.Programs {
		prgm := c.Programs[i]
		prgm.Energy -= 1
		cmdStr := path + "/" + prgm.ID
		cmd := exec.Command("coffee", cmdStr)
		//
		cmd.Stderr = os.Stderr
		stdoutBuffer := NewBuffer()
		cmd.Stdout = stdoutBuffer
		//
		stdinBuffer := NewBuffer()
		var stdinTap chan []byte
		cmd.Stdin, stdinTap = stdinBuffer.Pipe(fountain[i])
		open := true
		var data []byte
		var assert []float64
		for open == true {
			var d []byte
			d, open = <-stdinTap
			data = append(data, d...)
			lines := bytes.Split(data, []byte("\n"))
			data = []byte{}
			for i, _ := range lines {
				nums := bytes.Split(lines[i], []byte(" "))
				if len(nums) == inputs+1 && len(nums[inputs]) != 0 {
					num, err := strconv.ParseFloat(string(nums[inputs]), 64)
					if err == nil {
						assert = append(assert, num)
					}
				} else {
					data = append(data, lines[i]...)
				}
			}
		}
		//
		//
		prgmBytes, _ := prgm.MarshalProgram()
		fmt.Println(i, "Command - '"+cmdStr+"'")
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
		open = true
		data = []byte{}
		output := []float64{}
		for open == true {
			var d []byte
			d, open = <-stdoutTap
			data = append(data, d...)
		}
		lines := bytes.Split(data, []byte("\n"))
		for i, _ := range lines {
			num, err := strconv.ParseFloat(string(lines[i]), 64)
			if err == nil {
				output = append(output, num)
			} else {
				data = append(data, lines[i]...)
			}
		}
		// Compair output to asset
		if len(assert) == len(output) {
			avgScore := 0.0
			for i, _ := range assert {
				diff := output[i] - assert[i]
				avgScore += diff
			}
			score := ProgramScore{
				Index: i,
				Score: math.Abs(avgScore / float64(len(assert))),
			}
			scores = append(scores, score)
			fmt.Println("Score - ", score.Score)
		} else {
			fmt.Println("Program provided incorect amount of outputs")
		}
	}

	//	Each in top 30% ->
	//		* Add Energy
	//		* Cross with other top 30%
	sort.Sort(scores)
	limit := len(scores) / 3
	parents := []*ProgramInstance{}
	children := []*ProgramInstance{}
	for i, _ := range scores {
		if i < limit {
			c.Programs[scores[i].Index].Energy += 20
			if c.Programs[scores[i].Index].Energy > 100 {
				c.Programs[scores[i].Index].Energy = 100
			}
			parents = append(parents, c.Programs[scores[i].Index])
		}
	}

	if len(parents) > 1 {
		for i, _ := range parents {
			mate := i
			for mate != i {
				mate = int(util.RandomNumber(0, len(parents)-1))
			}
			pgm := &ProgramInstance{
				Program: parents[i].Mate(parents[mate].Program),
				Energy:  100,
				ID:      util.RandomHex(16),
			}
			children = append(children, pgm)
		}
	}

	fmt.Println(parents, children)

	//	Each program in population ->
	//		* If energy <= 0
	//			* dead
	//	Get population - dead
}

func (c *Context) Fitest() *program.Program {
	return nil
}

func (c *Context) InitPopulation(inputs, population int) {
	c.Programs = make([]*ProgramInstance, population)
	var i int
	for i = 0; i < population; i++ {
		pgm := &ProgramInstance{
			Program: program.New(inputs, inputs*4, &gene.MathBuildingBlock{}),
			Energy:  100,
			ID:      util.RandomHex(16),
		}
		c.Programs[i] = pgm
	}
}
