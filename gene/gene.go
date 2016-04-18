package gene

import (
	_ "fmt"
)

const (
	CursorNil       int = 0
	CursorVariable  int = 1
	CursorNumber    int = 2
	CursorOperator  int = 3
	CursorSeparator int = 4
)

type Gene []byte

func (g Gene) Heal() {
	switch g[len(g)-1] {
	case byte('+'), byte('-'), byte('*'), byte('/'):
		// g[len(g)-1] = byte(' ')
		g[len(g)-1] = 0
	}
}

func (g Gene) MarshalTree() (*GeneNode, error) {
	cursor := CursorNil
	var root *GeneNode = nil
	var current *GeneNode = nil
	var numberNode *GeneNode = nil
	var variableNode *GeneNode = nil
	for _, chrom := range g {
		switch chrom {
		case byte('$'), byte('a'), byte('b'), byte('c'), byte('d'), byte('e'), byte('f'), byte('g'), byte('h'), byte('i'), byte('j'), byte('k'), byte('l'), byte('m'), byte('n'), byte('o'), byte('p'), byte('q'), byte('r'), byte('s'), byte('t'), byte('u'), byte('v'), byte('w'), byte('x'), byte('y'), byte('z'):
			if cursor == CursorVariable {
				variableNode.Value += string(chrom)
			} else {
				node := &GeneNode{
					Value:    string(chrom),
					Children: []*GeneNode{},
				}
				current.Add(node)
				variableNode = node
			}
			cursor = CursorVariable
		case byte('0'), byte('1'), byte('2'), byte('3'), byte('4'), byte('5'), byte('6'), byte('7'), byte('8'), byte('9'):
			if cursor == CursorNumber {
				numberNode.Value += string(chrom)
			} else {
				node := &GeneNode{
					Value:    string(chrom),
					Children: []*GeneNode{},
				}
				current.Add(node)
				numberNode = node
			}
			cursor = CursorNumber
		case byte('+'), byte('-'), byte('*'), byte('/'):
			if cursor != CursorNil {
				node := &GeneNode{
					Value:    string(chrom),
					Children: []*GeneNode{},
				}
				current.Add(node)
				current = node
			} else {
				root = &GeneNode{
					Value:    string(chrom),
					Children: []*GeneNode{},
				}
				current = root
			}
			cursor = CursorOperator
		case byte(','):
			cursor = CursorSeparator
		}
	}
	return root, nil
}
