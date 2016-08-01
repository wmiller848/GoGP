package gene

import (
	"math"
	"strconv"
	// "github.com/wmiller848/govaluate"
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
	if len(n.Children) > 1 {
		bytes = append(bytes, byte('('))
	}
	for i, child := range n.Children {
		if i != 0 {
			bytes = append(bytes, []byte(n.Value)...)
		}
		//if len(child.Children) > 1 {
		//bytes = append(bytes, byte('('))
		//}
		cbytes, _ := child.MarshalExpression()
		bytes = append(bytes, cbytes...)
		//if len(child.Children) > 1 {
		//bytes = append(bytes, byte(')'))
		//}
	}
	if len(n.Children) > 1 {
		bytes = append(bytes, byte(')'))
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
			return val
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
				val = float64(int64(val) | int64(_val))
			}
		case "&":
			if math.IsNaN(val) {
				val = _val
			} else {
				val = float64(int64(val) & int64(_val))
			}
		case "^":
			if math.IsNaN(val) {
				val = _val
			} else {
				val = float64(int64(val) ^ int64(_val))
			}
		}
	}
	if len(n.Children) == 0 {
		var err error
		val, err = strconv.ParseFloat(n.Value, 64)
		if err != nil && IsVariable(n.Value) {
			val = inputs[VariableLookup(n.Value)]
		}
	}
	return val
}

// func (n *GeneNode) Variables(inputs ...float64) map[string]interface{} {
// 	m := make(map[string]interface{})
// 	if len(n.Children) == 0 {
// 		var err error
// 		val, err := strconv.ParseFloat(n.Value, 64)
// 		if err != nil && IsVariable(n.Value) {
// 			val = inputs[VariableLookup(n.Value)]
// 			m[n.Value] = val
// 		}
// 	} else {
// 		for i, _ := range n.Children {
// 			n := n.Children[i].Variables(inputs...)
// 			for k, v := range n {
// 				m[k] = v
// 			}
// 		}
// 	}
// 	return m
// }
//
// func (n *GeneNode) Eval(inputs ...float64) float64 {
// 	expr, _ := n.MarshalExpression()
// 	expression, err := govaluate.NewEvaluableExpression(string(expr))
// 	// expression, err := govaluate.NewEvaluableExpression("$ab + 1")
// 	if err != nil {
// 		// fmt.Println(string(expr), err.Error())
// 		return math.NaN()
// 	}
//
// 	parameters := n.Variables(inputs...)
// 	// parameters := make(map[string]interface{})
// 	// parameters["\"$ab\""] = 1024
//
// 	result, err := expression.Evaluate(parameters)
// 	if err != nil {
// 		return math.NaN()
// 	}
// 	return result.(float64)
// }
