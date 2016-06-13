package gene

import (
	"strconv"
)

const (
	CursorNil           int = 0
	CursorVariable      int = 1
	CursorVariableStart int = 2
	CursorNumber        int = 3
	CursorOperator      int = 4
	CursorSeparator     int = 5
	CursorScopeStart    int = 6
	CursorScopeStop     int = 7
)

type Gene interface {
	Eq(Gene) bool
	Clone() []byte
	Heal() []byte
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

func VariableTemplate(count int) string {
	tmpl := ""
	for i := 0; i < count; i++ {
		tmpl += Variable(i) + " = Number(args[" + strconv.Itoa(i) + "]);"
	}
	return tmpl
}
