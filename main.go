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
	g := gene.Gene("*{-10,20},{+50,77},{$a/15,4}")
	fmt.Println(string(g))
	g.Heal()
	fmt.Println(string(g))
	root, _ := g.MarshalTree()
	exp, _ := root.MarshalExpression()
	fmt.Println(string(exp))
}
