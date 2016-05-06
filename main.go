package main

import (
	"fmt"
	"github.com/wmiller848/GoGP/gene"
	"github.com/wmiller848/GoGP/program"
)

func main() {
	//	GENE
	//	*21-30,7
	//
	//	TREE
	//	*
	//	21	-
	//			30	7
	//
	//	EXPRESSION
	//	21 * (30 - 7)

	// g := gene.Gene("+{-10,20{-25,11}{+9,7}}{+50,77}{$a/15,4}")
	// g := gene.Random(50)
	// fmt.Println(string(g))
	// g = g.Heal()
	// fmt.Println(string(g))
	pgm := program.New(2, 12, &gene.MathBuildingBlock{})
	// fmt.Println(string(pgm.Gene.Clone()))
	// root, _ := pgm.Gene.MarshalTree()
	// exp, _ := root.MarshalExpression()
	// fmt.Println(string(exp))
	pgmStr, _ := pgm.MarshalProgram()
	fmt.Println(string(pgmStr))
}
