package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/wmiller848/GoGP/context"
	// "github.com/urfave/cli"
	"github.com/codegangsta/cli"
)

func run(pipe io.Reader, threshold, score float64, inputs, population, generations int, auto, visual bool) error {
	if inputs <= 0 {
		return errors.New("Count mut be greater then 0")
	}

	ctx := context.New()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fitest := ctx.Fitest()
			prgm, _ := fitest.MarshalProgram()
			fmt.Printf("%+v\n", string(prgm))
			os.Exit(0)
		}
	}()
	if visual {
		go func() {
			ctx.RunWithInlineScore(pipe, threshold, score, inputs, population, generations, auto)
		}()
		ctx.NewTerminal()
		//if auto {
		//fmt.Println("Learning from population of", population, "for", inputs, "inputs")
		//} else {
		//fmt.Println("Learning from population of", population, "over", generations, "generations for", inputs, "inputs")
		//}
	} else {
		_, fitest := ctx.RunWithInlineScore(pipe, threshold, score, inputs, population, generations, auto)
		prgm, _ := fitest.MarshalProgram()
		fmt.Printf("%+v\n", string(prgm))
	}
	//if visual {
	//fmt.Println("Generation ID:", uuid)
	//fmt.Println("Program ID:", fitest.ID)
	//fmt.Printf("Total Score: %3.2f\n", (1.0-fitest.Score)*100.00)
	//fmt.Println("Sub Scores:")
	//for k, grp := range fitest.Group {
	//c := float64(grp.Wrong) / float64(grp.Count)
	//fmt.Printf("%v: %3.2f (%v / %v)\n", k, (1.0-c)*100.00, grp.Count-grp.Wrong, grp.Count)
	//}
	//}
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
				cli.Float64Flag{
					Name:   "score, s",
					Usage:  "Float value for desired percentage score of program. 0.10 = 10%",
					EnvVar: "GOGP_SCORE",
					Value:  0.90,
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
				run(pipe, c.Float64("threshold"), c.Float64("score"), c.Int("count"), c.Int("population"), c.Int("generations"), !c.Bool("auto"), c.Bool("verbose"))
			},
		},
	}
	app.Run(os.Args)
}
