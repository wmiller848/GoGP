package gene

type MathGene GenericGene

func (g MathGene) Bytes() []byte {
	return []byte(g)
}

func (g MathGene) Eq(ng Gene) bool {
	lg := len(g)
	if lg != ng.Len() {
		return false
	}

	for i := 0; i < lg; i++ {
		if g[i] != ng.At(i) {
			return false
		}
	}
	return true
}

func (g MathGene) Clone() Gene {
	hg := MathGene{}
	for i, _ := range g {
		if g[i] != 0x00 {
			hg = append(hg, g[i])
		}
	}
	return hg
}

func (g MathGene) Len() int {
	return len(g)
}

func (g MathGene) At(i int) byte {
	if i < 0 || i > len(g) {
		return 0x00
	}
	return g[i]
}

func (g MathGene) Heal() Gene {
	healed := []byte{}
	valid := false
	gne := g.Clone().Bytes()
	cursor := CursorNil
	cursorScope := CursorNil
	scope := 0
	for i, _ := range gne {
		switch gne[i] {
		case byte('$'):
			if valid == true {
				switch cursor {
				case CursorVariable, CursorNumber:
					healed = append(healed, byte(','))
				}
				healed = append(healed, gne[i])
				cursor = CursorVariableStart
			}
		case byte('a'), byte('b'), byte('c'), byte('d'), byte('e'), byte('f'), byte('g'), byte('h'), byte('i'), byte('j'), byte('k'), byte('l'), byte('m'), byte('n'), byte('o'), byte('p'), byte('q'), byte('r'), byte('s'), byte('t'), byte('u'), byte('v'), byte('w'), byte('x'), byte('y'), byte('z'):
			if valid == true {
				healed = append(healed, gne[i])
				cursor = CursorVariable
			}
		case byte('{'):
			if valid == true {
				healed = append(healed, gne[i])
				cursor = CursorScopeStart
				cursorScope = CursorScopeStart
				scope++
			}
		case byte('}'):
			if valid == true {
				switch cursorScope {
				case CursorScopeStart:
				case CursorScopeStop:
					if scope <= 0 {
						continue
					}
				default:
					continue
				}
				switch cursor {
				case CursorSeparator, CursorOperator:
					healed[len(healed)-1] = 0x00
				}
				healed = append(healed, gne[i])
				cursor = CursorScopeStop
				cursorScope = CursorScopeStop
				scope--
				if scope <= 0 {
					cursorScope = CursorNil
				}
			}
		case byte(','):
			if valid == true {
				switch cursor {
				case CursorSeparator, CursorOperator, CursorScopeStart, CursorScopeStop:
					continue
				}
				healed = append(healed, gne[i])
				cursor = CursorSeparator
			}
		case byte('*'), byte('/'), byte('+'), byte('-'), byte('&'), byte('|'), byte('^'):
			valid = true
			if valid == true {
				switch cursor {
				case CursorSeparator, CursorOperator:
					continue
				}
				healed = append(healed, gne[i])
				cursor = CursorOperator
			}
		case byte('0'), byte('1'), byte('2'), byte('3'), byte('4'), byte('5'), byte('6'), byte('7'), byte('8'), byte('9'):
			if valid == true {
				switch cursor {
				case CursorVariable:
					healed = append(healed, byte(','))
					if gne[i] == byte('0') {
						continue
					}
				case CursorNumber:
				default:
					if gne[i] == byte('0') {
						continue
					}
				}
				healed = append(healed, gne[i])
				cursor = CursorNumber
			}
		}
	}
	for i := len(healed) - 1; i > 0; i-- {
		switch healed[i] {
		case byte(','), byte('*'), byte('/'), byte('+'), byte('-'), byte('&'), byte('|'), byte('^'):
			healed[i] = 0x00
			continue
		}
		break
	}
	for i := 0; i < scope; i++ {
		healed = append(healed, byte('}'))
	}
	return MathGene(healed).Clone()
}

func (g MathGene) MarshalTree() (*GeneNode, error) {
	cursor := CursorNil
	var root *GeneNode = nil
	contextRoot := []*GeneNode{}
	var current *GeneNode = nil
	var numberNode *GeneNode = nil
	var variableNode *GeneNode = nil
	var lastOp string
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
		case byte('+'), byte('-'), byte('*'), byte('/'), byte('&'), byte('|'), byte('^'):
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
			lastOp = string(chrom)
			cursor = CursorOperator
		case byte('{'):
			node := &GeneNode{
				Value:    lastOp,
				Children: []*GeneNode{},
			}
			current.Add(node)
			contextRoot = append(contextRoot, current)
			current = node
			cursor = CursorSeparator
		case byte('}'):
			current = contextRoot[len(contextRoot)-1]
			contextRoot = contextRoot[:len(contextRoot)-1]
			cursor = CursorSeparator
		default:
			cursor = CursorSeparator
		}
	}
	return root, nil
}
