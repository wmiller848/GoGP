package main

import (
	"fmt"

	"github.com/wmiller848/GoGP/gene"
	"github.com/wmiller848/GoGP/program"
)

func main() {
	pgm := program.New(2, 10, &gene.MathBuildingBlock{})
	pgmStr, _ := pgm.MarshalProgram()
	fmt.Println(string(pgmStr))
}
