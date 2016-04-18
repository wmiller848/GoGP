package gene

type GeneNode struct {
	Value    string
	Children []*GeneNode
}

func (n *GeneNode) Add(child *GeneNode) {
	n.Children = append(n.Children, child)
}

func (n *GeneNode) MarshalExpression() ([]byte, error) {
	bytes := []byte{}
	for i, child := range n.Children {
		if len(child.Children) != 0 {
			bytes = append(bytes, byte('('))
		}
		cbytes, _ := child.MarshalExpression()
		bytes = append(bytes, cbytes...)
		if len(child.Children) != 0 {
			bytes = append(bytes, byte(')'))
		}
		if i != len(n.Children)-1 {
			bytes = append(bytes, []byte(n.Value)...)
		}
	}
	if len(n.Children) == 0 {
		bytes = append(bytes, []byte(n.Value)...)
	}
	return bytes, nil
}
