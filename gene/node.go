package gene

import (
	"math"
	"strconv"
)

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

func (n *GeneNode) Eval(inputs ...float64) float64 {
	val := math.NaN()
	for _, child := range n.Children {
		_val := child.Eval(inputs...)
		if math.IsNaN(_val) {
			continue
		}
		switch n.Value {
		case "*":
			if math.IsNaN(val) {
				val = _val
			} else {
				val *= _val
			}
		case "/":
			if math.IsNaN(val) {
				val = _val
			} else {
				val /= _val
			}
		case "+":
			if math.IsNaN(val) {
				val = _val
			} else {
				val += _val
			}
		case "-":
			if math.IsNaN(val) {
				val = _val
			} else {
				val -= _val
			}
		case "|":
			if math.IsNaN(val) {
				val = _val
			} else {
				val = float64(uint64(val) | uint64(_val))
			}
		case "&":
			if math.IsNaN(val) {
				val = _val
			} else {
				val = float64(uint64(val) & uint64(_val))
			}
		case "^":
			if math.IsNaN(val) {
				val = _val
			} else {
				val = float64(uint64(val) ^ uint64(_val))
			}
		}
	}
	if len(n.Children) == 0 {
		var err error
		val, err = strconv.ParseFloat(n.Value, 64)
		if err != nil {
			// TODO :: Make this dynamic
			switch n.Value {
			case "$a":
				val = inputs[0]
			case "$b":
				val = inputs[1]
			case "$c":
				val = inputs[2]
			case "$d":
				val = inputs[3]
			}
		}
	}
	return val
}
