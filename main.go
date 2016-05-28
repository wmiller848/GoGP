package main

import (
	"fmt"
	"io"
	"os"

	"github.com/codegangsta/cli"
	"github.com/wmiller848/GoGP/context"
)

func score(output int) int {
	return 0
}

func run(pipe io.Reader, inputs int) {
	ctx := context.New()
	population := 10
	generations := 1
	fmt.Println("Learning from population of", population, "over", generations, "generations for", inputs, "inputs")
	ctx.RunWithInlineScore(pipe, inputs, population, generations)
}

func main() {
	app := cli.NewApp()
	app.Name = "GoGP"
	app.Usage = "Learn a program for matching data based off training data"
	app.Commands = []cli.Command{
		{
			Name:    "learn",
			Aliases: []string{"l"},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:   "count, c",
					Usage:  "Number of input fields to learn",
					EnvVar: "GOGP_COUNT",
				},
			},
			Action: func(c *cli.Context) {
				args := c.Args()
				var pipe io.Reader
				if len(args) == 0 {
					//// Handle stdin
					//bio := bufio.NewReader(os.Stdin)
					//for {
					//line, _, err := bio.ReadLine()
					//if err != nil {
					//break
					//}
					//buf = append(buf, line...)
					//buf = append(buf, byte('\n'))
					//}
					pipe = os.Stdin
				} else if len(args) == 1 {
					// Handle file path
					//inputFile := args[0]
					//buf, err := ioutil.ReadFile(inputFile)
					//if err != nil {
					//fmt.Println(err.Error())
					//return
					//}
				} else {
					fmt.Println("Too many arguments, provide path to one file.")
					return
				}
				run(pipe, c.Int("count"))
			},
		},
	}
	app.Run(os.Args)
}
