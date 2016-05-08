package gene

type MathGene GenericGene

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

func (g MathGene) Clone() []byte {
	hg := MathGene{}
	for i, _ := range g {
		if g[i] != 0x00 {
			hg = append(hg, g[i])
		}
	}
	return hg
}

func (g MathGene) LastChrome(i int) int {
	var lx int
	for {
		lx = i - 1
		if lx >= 0 {
			if g[lx] != 0x00 {
				break
			}
		}
		i--
		if i <= 0 {
			lx = -1
			break
		}
	}
	return lx
}

func (g MathGene) NextChrome(i int) int {
	var nx int
	for {
		nx = i + 1
		if nx < len(g) {
			if g[nx] != 0x00 {
				break
			}
		}
		i++
		if i >= len(g) {
			nx = -1
			break
		}
	}
	return nx
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

func (g MathGene) Heal() []byte {
	for {
		clean := true
		g = g.Clone()
		for i, _ := range g {
			lx := g.LastChrome(i)
			if lx >= 0 && g[i] != 0x00 {
				switch g[lx] {
				case byte('{'), byte('}'):
					switch g[i] {
					case byte(','):
						g[i] = 0x00
						clean = false
					case byte('+'), byte('-'), byte('*'), byte('/'):
						nx := g.NextChrome(i)
						if nx >= 0 && g[nx] == byte('}') {
							g[lx] = 0x00
							g[i] = 0x00
							g[nx] = 0x00
						}
					}
				case byte(','):
					switch g[i] {
					case byte(','):
						g[i] = 0x00
						clean = false
					case byte('+'), byte('-'), byte('*'), byte('/'):
						g[i] = 0x00
						clean = false
					}
				case byte('+'), byte('-'), byte('*'), byte('/'):
					switch g[i] {
					case byte(','):
						g[i] = 0x00
						clean = false
					case byte('+'), byte('-'), byte('*'), byte('/'):
						g[i] = 0x00
						clean = false
					}
				}
			}
		}
		if clean == true {
			break
		}
	}

	for {
		clean := true
		g = g.Clone()
		lg := len(g) - 1
		if lg > 0 {
			switch g[lg] {
			case byte('+'), byte('-'), byte('*'), byte('/'):
				g[lg] = 0x00
				clean = false
			case byte(','):
				g[lg] = 0x00
				clean = false
			default:
			}
			if clean == true {
				break
			}
		} else {
			break
		}
	}

	return g
}

func (g MathGene) MarshalTree() (*GeneNode, error) {
	cursor := CursorNil
	var root *GeneNode = nil
	contextRoot := []*GeneNode{}
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
		case byte('{'):
			contextRoot = append(contextRoot, current)
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
