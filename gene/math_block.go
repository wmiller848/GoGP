package gene

import (
	"github.com/wmiller848/GoGP/util"
	_ "strconv"
)

var blockVars []byte = []byte{
	'a', 'b', 'c', 'd', 'e', 'f',
	'g', 'h', 'i', 'j', 'k', 'l',
	'm', 'n', 'o', 'p', 'q', 'r',
	's', 't', 'u', 'v', 'w', 'x',
	'y', 'z',
}

type MathBuildingBlock struct {
}

func (m *MathBuildingBlock) Stages() []int {
	return []int{StageAlive}
}

func (m *MathBuildingBlock) Random(varCount, knobCount int) (Gene, string) {
	g := MathGene{randomOperator()}
	cursor := CursorOperator
	size := util.RandomNumber(knobCount, knobCount*2) // random value
	c := int(size)
	k := varCount
	i, j := 0, 0
	for i < c || j < k {
		pick := util.RandomNumber(0, 100)
		if i < c {
			if pick < 50 {
				g = append(g, randomNumber()+48)
				cursor = CursorNumber
			} else if pick < 90 && cursor != CursorOperator {
				g = append(g, randomOperator())
				cursor = CursorOperator
			} else if pick < 95 && knobCount > 4 {
				g = append(g, byte('{'))
				g = append(g, m.Random(0, knobCount/4).Clone()...)
				g = append(g, byte('}'))
				cursor = CursorSeparator
			} else {
				g = append(g, byte(','))
				cursor = CursorSeparator
			}
			i++
		} else if cursor == CursorNumber {
			g = append(g, randomOperator())
			cursor = CursorOperator
		}
		if pick > 50 && cursor != CursorNumber && j < k {
			c := ""
			ji := j % len(blockVars)
			if j != 0 && ji == 0 {
				jd := j / len(blockVars)
				for t := 0; t < jd; t++ {
					c += string(blockVars[t])
				}
			}
			g = append(g, []byte("$"+c+string(blockVars[ji])+",")...)
			cursor = CursorVariable
			j++
		}
	}
	g = MathGene(g.Heal())
	return g
}
