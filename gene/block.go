package gene

import (
	"strconv"

	"github.com/wmiller848/GoGP/util"
)

func randomOperator() byte {
	switch util.RandomNumber(0, 3) {
	case 0:
		return byte('+')
	case 1:
		return byte('-')
	case 2:
		return byte('*')
	case 3:
		return byte('/')
	default:
		return randomOperator()
	}
}

func randomVariable() byte {
	return byte(util.RandomNumber(0, 9))
}

func randomNumber() byte {
	return byte(util.RandomNumber(0, 9))
}

const (
	StageSpawn  int = 1
	StageAlive  int = 2
	StageDieing int = 3
	StageDead   int = 3
)

type BuildingBlock interface {
	Stages() []int
	// Bytes() [][]byte
	// Distribution() (map[string]int, int)
	Random(int, int) Gene
}

func GetVariableBlock(j int) string {
	tmpl := "$"
	ji := j % len(blockVars)
	if j != 0 && ji == 0 {
		jd := j / len(blockVars)
		for t := 0; t < jd; t++ {
			tmpl += string(blockVars[t])
		}
	}
	tmpl += string(blockVars[ji])
	return tmpl
}

func VarsTemplate(g Gene) string {
	tmpl := ""
	cursor := CursorNil
	j := 0
	for i, _ := range g.Clone() {
		switch g.At(i) {
		case byte('$'):
			cursor = CursorVariable
			tmpl += "$"
		case byte('a'), byte('b'), byte('c'), byte('d'), byte('e'), byte('f'), byte('g'), byte('h'), byte('i'), byte('j'), byte('k'), byte('l'), byte('m'), byte('n'), byte('o'), byte('p'), byte('q'), byte('r'), byte('s'), byte('t'), byte('u'), byte('v'), byte('w'), byte('x'), byte('y'), byte('z'):
			if cursor == CursorVariable {
				cursor = CursorVariable
				tmpl += string(g.At(i))
			}
		default:
			if cursor == CursorVariable {
				cursor = CursorNil
				tmpl += " = args[" + strconv.Itoa(j) + "];"
				j++
			}
		}
	}

	return tmpl
}
