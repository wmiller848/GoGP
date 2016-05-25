package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
	"github.com/wmiller848/GoGP/context"
)

func score(output int) int {
	return 0
}

func run(buf []byte, inline bool) {
	ctx := context.New()
	lines := bytes.Split(buf, []byte("\n"))
	inputs := len(bytes.Split(lines[0], []byte(" ")))
	population := 1
	generations := 1
	if inline == true {
		inputs -= 1
		fmt.Println("Learning from population of", population, "over", generations, "generations for", inputs, "inputs across", len(lines), "rows")
		ctx.RunWithInlineScore(buf, inputs, population, generations)
	} else {
		fmt.Println("Scoring Function Support Not Implemented")
		// ctx.RunWithScoreFunc(buf, inputs, population, generations, score)
	}
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
				cli.BoolFlag{
					Name:   "inline, i",
					Usage:  "Traing data has anwser inline.",
					EnvVar: "GOGP_INLINE",
				},
			},
			Action: func(c *cli.Context) {
				args := c.Args()
				buf := []byte{}
				if len(args) == 0 {
					// Handle stdin
					bio := bufio.NewReader(os.Stdin)
					for {
						line, _, err := bio.ReadLine()
						if err != nil {
							break
						}
						buf = append(buf, line...)
						buf = append(buf, byte('\n'))
					}
				} else if len(args) == 1 {
					// Handle file path
					inputFile := args[0]
					var err error
					buf, err = ioutil.ReadFile(inputFile)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
				} else {
					fmt.Println("Too many arguments, provide path to one file.")
					return
				}
				run(buf, c.Bool("inline"))
			},
		},
	}
	app.Run(os.Args)
}
