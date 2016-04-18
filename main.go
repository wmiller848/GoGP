package main

import (
	"fmt"
	"github.com/wmiller848/GoGP/gene"
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
	g := gene.Gene("*$a-,30,10,$b,12*10,70,")
	g.Heal()
	fmt.Println(string(g))
	root, _ := g.MarshalTree()
	exp, _ := root.MarshalExpression()
	fmt.Println(string(exp))
}
