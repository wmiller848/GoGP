package gene

import (
	_ "github.com/wmiller848/GoGP/util"
	"strconv"
)

const (
	CursorNil       int = 0
	CursorVariable  int = 1
	CursorNumber    int = 2
	CursorOperator  int = 3
	CursorSeparator int = 4
)

type Gene interface {
	Eq(Gene) bool
	Clone() []byte
	Heal() []byte
	LastChrome(int) int
	NextChrome(int) int
	Len() int
	At(int) byte
	MarshalTree() (*GeneNode, error)
}

type GenericGene []byte

var blockVars []byte = []byte{
	'a', 'b', 'c', 'd', 'e', 'f',
	'g', 'h', 'i', 'j', 'k', 'l',
	'm', 'n', 'o', 'p', 'q', 'r',
	's', 't', 'u', 'v', 'w', 'x',
	'y', 'z',
}

func Variable(j int) string {
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

func VariableTemplate(g Gene) string {
	vars := make(map[string]string)
	tmpl := ""
	cursor := CursorNil
	j := 0
	for i, _ := range g.Clone() {
		switch g.At(i) {
		case byte('$'):
			if tmpl != "" && vars[tmpl] == "" {
				vars[tmpl] = tmpl + " = args[" + strconv.Itoa(j) + "];"
				j++
			}
			tmpl = string(g.At(i))
			cursor = CursorVariable
		case byte('a'), byte('b'), byte('c'), byte('d'), byte('e'), byte('f'), byte('g'), byte('h'), byte('i'), byte('j'), byte('k'), byte('l'), byte('m'), byte('n'), byte('o'), byte('p'), byte('q'), byte('r'), byte('s'), byte('t'), byte('u'), byte('v'), byte('w'), byte('x'), byte('y'), byte('z'):
			tmpl += string(g.At(i))
			cursor = CursorVariable
		default:
			if cursor == CursorVariable && vars[tmpl] == "" {
				vars[tmpl] = tmpl + " = args[" + strconv.Itoa(j) + "];"
				j++
			}
			cursor = CursorNil
		}
	}
	if tmpl != "" && vars[tmpl] == "" {
		vars[tmpl] = tmpl + " = args[" + strconv.Itoa(j) + "];"
	}

	tmpl = ""
	for _, val := range vars {
		tmpl += val
	}
	return tmpl
}
