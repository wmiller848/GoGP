package gene

import (
	"strconv"
	"strings"

	"github.com/wmiller848/GoGP/util"
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
	Bytes() []byte
	Eq(Gene) bool
	Clone() Gene
	Heal() Gene
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
	'y', 'z', '_',
}

func IsVariable(v string) bool {
	for i, _ := range blockVars {
		for t, _ := range v {
			if blockVars[i] == v[t] {
				return true
			}
		}
	}
	return false
}

func Variable(j int) string {
	tmpl := "$"
	tmpl += util.Hex([]byte{byte(j)})
	tmpl = strings.Replace(tmpl, "0", "z", -1)
	tmpl = strings.Replace(tmpl, "1", "y", -1)
	tmpl = strings.Replace(tmpl, "2", "x", -1)
	tmpl = strings.Replace(tmpl, "3", "w", -1)
	tmpl = strings.Replace(tmpl, "4", "v", -1)
	tmpl = strings.Replace(tmpl, "5", "u", -1)
	tmpl = strings.Replace(tmpl, "6", "t", -1)
	tmpl = strings.Replace(tmpl, "7", "s", -1)
	tmpl = strings.Replace(tmpl, "8", "r", -1)
	tmpl = strings.Replace(tmpl, "9", "q", -1)
	return tmpl
}

func VariableLookup(v string) int {
	if string(v[:1]) == "$" {
		str := string(v[1:])
		str = strings.Replace(str, "z", "0", -1)
		str = strings.Replace(str, "y", "1", -1)
		str = strings.Replace(str, "x", "2", -1)
		str = strings.Replace(str, "w", "3", -1)
		str = strings.Replace(str, "v", "4", -1)
		str = strings.Replace(str, "u", "5", -1)
		str = strings.Replace(str, "t", "6", -1)
		str = strings.Replace(str, "s", "7", -1)
		str = strings.Replace(str, "r", "8", -1)
		str = strings.Replace(str, "q", "9", -1)
		byts := util.Unhex(str)
		return int(byts[0])
	}
	return 0
}

func VariableTemplate(count int) string {
	tmpl := ""
	for i := 0; i < count; i++ {
		tmpl += Variable(i) + " = Number(args[" + strconv.Itoa(i) + "]);"
	}
	return tmpl
}
