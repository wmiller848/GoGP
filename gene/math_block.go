package gene

import (
	_ "strconv"

	"github.com/wmiller848/GoGP/util"
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

func (m *MathBuildingBlock) Random(varCount, knobCount int) Gene {
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
				var num byte = 0
				if cursor != CursorNumber {
					for num <= 48 {
						num = randomNumber() + 48
					}
				} else {
					num = randomNumber() + 48
				}
				g = append(g, num)
				cursor = CursorNumber
			} else if pick < 90 && cursor != CursorOperator {
				g = append(g, randomOperator())
				cursor = CursorOperator
			} else if pick < 95 && knobCount > 4 {
				g = append(g, byte('{'))
				g = append(g, m.Random(0, knobCount/2).Clone()...)
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
			g = append(g, []byte(GetVariableBlock(j)+",")...)
			cursor = CursorVariable
			j++
		}
	}
	g = MathGene(g.Heal())
	return g
}
