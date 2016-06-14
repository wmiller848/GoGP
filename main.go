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

func run(pipe io.Reader, inputs, population, generations int, verbose bool) {
	ctx := context.New()
	if verbose {
		ctx.Verbose()
	}
	if verbose {
		fmt.Println("Learning from population of", population, "over", generations, "generations for", inputs, "inputs")
	}
	uuid, fitest := ctx.RunWithInlineScore(pipe, inputs, population, generations)
	prgm, _ := fitest.MarshalProgram()
	if verbose {
		fmt.Println(uuid)
	}
	fmt.Printf("%+v\n", string(prgm))
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
					Value:  0,
				},
				cli.IntFlag{
					Name:   "population, p",
					Usage:  "Number of programs to keep in the pool",
					EnvVar: "GOGP_POPULATION",
					Value:  50,
				},
				cli.IntFlag{
					Name:   "generations, g",
					Usage:  "Number of generations to iterate through",
					EnvVar: "GOGP_GENERATIONS",
					Value:  100,
				},
				cli.BoolFlag{
					Name:   "verbose, v",
					Usage:  "Output more then just the evolved program",
					EnvVar: "GOGP_VERBOSE",
				},
			},
			Action: func(c *cli.Context) {
				args := c.Args()
				var pipe io.Reader
				if len(args) == 0 {
					pipe = os.Stdin
				} else if len(args) == 1 {
					// Handle file io
				} else {
					fmt.Println("Too many arguments, provide path to one file.")
					return
				}
				run(pipe, c.Int("count"), c.Int("population"), c.Int("generations"), c.Bool("verbose"))
			},
		},
	}
	app.Run(os.Args)
}
