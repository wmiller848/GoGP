package dna

import (
	_ "strconv"

	"github.com/wmiller848/GoGP/util"
)

var bases []byte = []byte{
	'a', 'b', 'c', 'd',
}

type MathBlock struct {
}

func (m *MathBlock) Random(inputs int) {
	var baseBlockMap map[Base]map[Base]map[Base]Codon = make(map[Base]map[Base]map[Base]Codon)
	//
	baseBlockMap['a'] = make(map[Base]map[Base]Codon)
	////
	baseBlockMap['a']['a'] = make(map[Base]Codon)
	baseBlockMap['a']['a']['a'] = []byte("<")
	baseBlockMap['a']['a']['b'] = []byte("*")
	baseBlockMap['a']['a']['c'] = []byte("*")
	baseBlockMap['a']['a']['d'] = []byte("*")
	////
	baseBlockMap['a']['b'] = make(map[Base]Codon)
	baseBlockMap['a']['b']['a'] = []byte("/")
	baseBlockMap['a']['b']['b'] = []byte("/")
	baseBlockMap['a']['b']['c'] = []byte("/")
	baseBlockMap['a']['b']['d'] = []byte("/")
	////
	baseBlockMap['a']['c'] = make(map[Base]Codon)
	baseBlockMap['a']['c']['a'] = []byte("#")
	baseBlockMap['a']['c']['b'] = []byte("#")
	baseBlockMap['a']['c']['c'] = []byte("#")
	baseBlockMap['a']['c']['d'] = []byte("#")
	////
	baseBlockMap['a']['d'] = make(map[Base]Codon)
	baseBlockMap['a']['d']['a'] = []byte("#")
	baseBlockMap['a']['d']['b'] = []byte("#")
	baseBlockMap['a']['d']['c'] = []byte("#")
	baseBlockMap['a']['d']['d'] = []byte("#")

	//
	//
	baseBlockMap['b'] = make(map[Base]map[Base]Codon)
	////
	baseBlockMap['b']['a'] = make(map[Base]Codon)
	baseBlockMap['b']['a']['a'] = []byte("+")
	baseBlockMap['b']['a']['b'] = []byte("+")
	baseBlockMap['b']['a']['c'] = []byte("+")
	baseBlockMap['b']['a']['d'] = []byte("+")
	////
	baseBlockMap['b']['b'] = make(map[Base]Codon)
	baseBlockMap['b']['b']['a'] = []byte("!")
	baseBlockMap['b']['b']['b'] = []byte("!")
	baseBlockMap['b']['b']['c'] = []byte("!")
	baseBlockMap['b']['b']['d'] = []byte("!")
	////
	baseBlockMap['b']['c'] = make(map[Base]Codon)
	baseBlockMap['b']['c']['a'] = []byte("!")
	baseBlockMap['b']['c']['b'] = []byte("!")
	baseBlockMap['b']['c']['c'] = []byte("!")
	baseBlockMap['b']['c']['d'] = []byte("!")
	////
	baseBlockMap['b']['d'] = make(map[Base]Codon)
	baseBlockMap['b']['d']['a'] = []byte("!")
	baseBlockMap['b']['d']['b'] = []byte("!")
	baseBlockMap['b']['d']['c'] = []byte("!")
	baseBlockMap['b']['d']['d'] = []byte("!")

	//
	//
	baseBlockMap['c'] = make(map[Base]map[Base]Codon)
	////
	baseBlockMap['c']['a'] = make(map[Base]Codon)
	baseBlockMap['c']['a']['a'] = []byte("!")
	baseBlockMap['c']['a']['b'] = []byte("!")
	baseBlockMap['c']['a']['c'] = []byte("!")
	baseBlockMap['c']['a']['d'] = []byte("!")
	////
	baseBlockMap['c']['b'] = make(map[Base]Codon)
	baseBlockMap['c']['b']['a'] = []byte("!")
	baseBlockMap['c']['b']['b'] = []byte("!")
	baseBlockMap['c']['b']['c'] = []byte("!")
	baseBlockMap['c']['b']['d'] = []byte("!")
	////
	baseBlockMap['c']['c'] = make(map[Base]Codon)
	baseBlockMap['c']['c']['a'] = []byte("!")
	baseBlockMap['c']['c']['b'] = []byte("!")
	baseBlockMap['c']['c']['c'] = []byte("!")
	baseBlockMap['c']['c']['d'] = []byte("!")
	////
	baseBlockMap['c']['d'] = make(map[Base]Codon)
	baseBlockMap['c']['d']['a'] = []byte("!")
	baseBlockMap['c']['d']['b'] = []byte("!")
	baseBlockMap['c']['d']['c'] = []byte("!")
	baseBlockMap['c']['d']['d'] = []byte("!")

	//
	//
	baseBlockMap['d'] = make(map[Base]map[Base]Codon)
	////
	baseBlockMap['d']['a'] = make(map[Base]Codon)
	baseBlockMap['d']['a']['a'] = []byte("!")
	baseBlockMap['d']['a']['b'] = []byte("!")
	baseBlockMap['d']['a']['c'] = []byte("!")
	baseBlockMap['d']['a']['d'] = []byte("!")
	////
	baseBlockMap['d']['b'] = make(map[Base]Codon)
	baseBlockMap['d']['b']['a'] = []byte("!")
	baseBlockMap['d']['b']['b'] = []byte("!")
	baseBlockMap['d']['b']['c'] = []byte("!")
	baseBlockMap['d']['b']['d'] = []byte("!")
	////
	baseBlockMap['d']['c'] = make(map[Base]Codon)
	baseBlockMap['d']['c']['a'] = []byte("!")
	baseBlockMap['d']['c']['b'] = []byte("!")
	baseBlockMap['d']['c']['c'] = []byte("!")
	baseBlockMap['d']['c']['d'] = []byte("!")
	////
	baseBlockMap['d']['d'] = make(map[Base]Codon)
	baseBlockMap['d']['d']['a'] = []byte("!")
	baseBlockMap['d']['d']['b'] = []byte("!")
	baseBlockMap['d']['d']['c'] = []byte(">")
	baseBlockMap['d']['d']['d'] = []byte(">")

	pick := util.RandomNumber(0, 3)
	switch pick {
	case 0:
	case 1:
	case 2:
	case 3:
	}
}

//func (m *MathBlock) Stages() []int {
//return []int{StageAlive}
//}

//func (m *MathBlock) Random(varCount, knobCount int) Gene {
//g := MathGene{randomOperator()}
//cursor := CursorOperator
//size := util.RandomNumber(knobCount, knobCount*2) // random value
//c := int(size)
//k := varCount
//i, j := 0, 0
//for i < c || j < k {
//pick := util.RandomNumber(0, 100)
//if i < c {
//if pick < 50 {
//var num byte = 0
//if cursor != CursorNumber {
//for num <= 48 {
//num = randomNumber() + 48
//}
//} else {
//num = randomNumber() + 48
//}
//g = append(g, num)
//cursor = CursorNumber
//} else if pick < 90 && cursor != CursorOperator {
//g = append(g, randomOperator())
//cursor = CursorOperator
//} else if pick < 95 && knobCount > 4 {
//g = append(g, byte('{'))
//g = append(g, m.Random(0, knobCount/2).Clone()...)
//g = append(g, byte('}'))
//cursor = CursorSeparator
//} else {
//g = append(g, byte(','))
//cursor = CursorSeparator
//}
//i++
//} else if cursor == CursorNumber {
//g = append(g, randomOperator())
//cursor = CursorOperator
//}
//if pick > 50 && cursor != CursorNumber && j < k {
//g = append(g, []byte(GetVariableBlock(j)+",")...)
//cursor = CursorVariable
//j++
//}
//}
//g = MathGene(g.Heal())
//return g
//}
