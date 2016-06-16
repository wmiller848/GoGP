package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	// "github.com/urfave/cli"
	"github.com/codegangsta/cli"
	"github.com/wmiller848/GoGP/context"
)

func score(output int) int {
	return 0
}

func run(pipe io.Reader, threshold float64, inputs, population, generations int, auto, verbose bool) error {
	if inputs <= 0 {
		return errors.New("Count mut be greater then 0")
	}

	ctx := context.New()
	if verbose {
		ctx.Verbose()
		if auto {
			fmt.Println("Learning from population of", population, "for", inputs, "inputs")
		} else {
			fmt.Println("Learning from population of", population, "over", generations, "generations for", inputs, "inputs")
		}
	}
	uuid, fitest := ctx.RunWithInlineScore(pipe, threshold, inputs, population, generations, auto)
	prgm, _ := fitest.MarshalProgram()
	if verbose {
		fmt.Println(uuid)
		fmt.Printf("%+v\n", fitest.Score)
	}
	fmt.Printf("%+v\n", string(prgm))
	return nil
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
				cli.Float64Flag{
					Name:   "threshold, t",
					Usage:  "Float value for how close the output needs to be to the training data",
					EnvVar: "GOGP_THRESHOLD",
					Value:  500.0,
				},
				cli.IntFlag{
					Name:   "population, p",
					Usage:  "Number of programs to keep in the pool",
					EnvVar: "GOGP_POPULATION",
					Value:  20,
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
				cli.BoolFlag{
					Name:   "auto, a",
					Usage:  "Run until a reasonable score is found",
					EnvVar: "GOGP_AUTO",
				},
			},
			// Action: func(c *cli.Context) error {
			// 	args := c.Args()
			// 	var pipe io.Reader
			// 	if len(args) == 0 {
			// 		pipe = os.Stdin
			// 	} else if len(args) == 1 {
			// 		// Handle file io
			// 	} else {
			// 		fmt.Println("Too many arguments, provide path to one file.")
			// 		return nil
			// 	}
			// 	return run(pipe, c.Float64("threshold"), c.Int("count"), c.Int("population"), c.Int("generations"), !c.Bool("auto"), c.Bool("verbose"))
			// },
			Action: func(c *cli.Context) {
				args := c.Args()
				var pipe io.Reader
				if len(args) == 0 {
					pipe = os.Stdin
				} else if len(args) == 1 {
					// Handle file io
				} else {
					fmt.Println("Too many arguments, provide path to one file.")
				}
				run(pipe, c.Float64("threshold"), c.Int("count"), c.Int("population"), c.Int("generations"), !c.Bool("auto"), c.Bool("verbose"))
			},
		},
	}
	app.Run(os.Args)
}
