package context

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

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

type Context struct {
	Programs []*ProgramInstance
}

func New() *Context {
	return &Context{}
}

func (c *Context) RunWithScoreFunc(inputs, population, generations int, scoreFunc ScoreFunction) *program.Program {
	var i int
	c.InitPopulation(inputs, population)
	for i = 0; i < generations; i++ {
		c.EvalFunc(scoreFunc)
	}

	return c.Fitest()
}

func (c *Context) EvalFunc(scoreFunc ScoreFunction) {

}

func (c *Context) RunWithInlineScore(pipe io.Reader, inputs, population, generations int) *program.Program {
	os.Mkdir("./out", 0777)
	sha := util.Random(32)
	os.Mkdir("./out/generations", 0777)
	os.RemoveAll("./out/generations/" + util.Hex(sha))
	os.Mkdir("./out/generations/"+util.Hex(sha), 0777)
	c.InitPopulation(inputs, population)
	var i int
	for i = 0; i < generations; i++ {
		c.EvalInline(pipe, i, sha)
	}

	return c.Fitest()
}

func (c *Context) EvalInline(pipe io.Reader, generation int, uuid []byte) {
	//fmt.Println("Traing Data", string(traingBuf), traingBuf)
	path := "./out/generations/" + util.Hex(uuid) + "/" + strconv.Itoa(generation)
	os.Mkdir(path, 0777)
	//	Each program in population ->
	//		* Apply 'life' ->
	//			* Drain Energy
	for i, _ := range c.Programs {
		prgm := c.Programs[i]
		prgm.Energy -= 1
	}

	//		* Each testBuf row ->
	//			* compute average score
	for i, _ := range c.Programs {
		prgm := c.Programs[i]
		prgm.Energy -= 1
		cmdStr := path + "/" + prgm.ID
		cmd := exec.Command("coffee", cmdStr)
		//
		cmd.Stderr = os.Stderr
		stdoutBuffer := NewBuffer()
		cmd.Stdout = stdoutBuffer
		cmd.Stdin = pipe
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
		//correctValues := make([]int, len(rows))
		//for j, _ := range rows {
		//row := rows[j]
		//rowArr := bytes.Split(row, []byte(" "))
		//correctVal, err := strconv.Atoi(string(rowArr[len(rowArr)-1]))
		//if err != nil {
		//continue
		//}
		//correctValues[j] = correctVal
		//input := append([]byte{}, bytes.Join(rowArr[:len(rowArr)-1], []byte(" "))...)
		//input = append(input, []byte("\n")...)
		//stdinBuffer.Write(input)
		//}
		err = cmd.Wait()
		if err != nil {
			fmt.Println(err.Error())
		}
		// Compair outputs to correct vals
		//fmt.Println(string(stdinBuffer.Data()), string(stdoutBuffer.Data()))
		//stdout := bytes.Split(stdoutBuffer.Data(), []byte("\n"))
		//fmt.Println(len(stdout), len(correctValues))
	}

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

func (c *Context) InitPopulation(inputs, population int) {
	c.Programs = make([]*ProgramInstance, population)
	var i int
	for i = 0; i < population; i++ {
		pgm := &ProgramInstance{
			Program: program.New(inputs, inputs*4, &gene.MathBuildingBlock{}),
			Energy:  1000,
			ID:      util.RandomHex(16),
		}
		c.Programs[i] = pgm
	}
}
