package gene

import (
	"testing"
)

func TestMarshalTreeGeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "*",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "10",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "30",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "-",
		Children: []*GeneNode{},
	})
	root.Children[2].Add(&GeneNode{
		Value:    "21",
		Children: []*GeneNode{},
	})
	root.Children[2].Add(&GeneNode{
		Value:    "77",
		Children: []*GeneNode{},
	})
	expr, err := root.MarshalExpression()
	if err != nil {
		t.Error(err.Error())
	}
	AssertStr(t, string(expr), "(10*30*(21-77))")
}

func TestEvalBasicGeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "+",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "10",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "30",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, 40.0)
}

func TestEvalGeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "*",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "10",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "30",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "-",
		Children: []*GeneNode{},
	})
	root.Children[2].Add(&GeneNode{
		Value:    "21",
		Children: []*GeneNode{},
	})
	root.Children[2].Add(&GeneNode{
		Value:    "77",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, -16800.0)
}

func TestEvalNonCommutative1GeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "/",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "10",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "20",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "30",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, 0.016666666666666666)
}

func TestEvalNonCommutative2GeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "/",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "30",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "20",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "10",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, 0.15)
}

func TestEvalNonCommutative3GeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "-",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "10",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "20",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "30",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, -40.0)
}

func TestEvalNonCommutative4GeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "-",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "30",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "20",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "10",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, 0.0)
}

func TestEvalCommutative1GeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "+",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "30",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "20",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "10",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, 60.0)
}

func TestEvalCommutative2GeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "+",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "10",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "20",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "30",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, 60.0)
}

func TestEvalCommutative3GeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "*",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "30",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "20",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "10",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, 6000.0)
}

func TestEvalCommutative4GeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "*",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "10",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "20",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "30",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, 6000.0)
}

func TestEvalOrderOfOperations1GeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "*",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "10",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "20",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "-",
		Children: []*GeneNode{},
	})
	root.Children[2].Add(&GeneNode{
		Value:    "7",
		Children: []*GeneNode{},
	})
	root.Children[2].Add(&GeneNode{
		Value:    "2",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, 1000.0)
}

func TestEvalOrderOfOperations2GeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "-",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "10",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "20",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "/",
		Children: []*GeneNode{},
	})
	root.Children[2].Add(&GeneNode{
		Value:    "7",
		Children: []*GeneNode{},
	})
	root.Children[2].Add(&GeneNode{
		Value:    "2",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, -13.5)
}

func TestEvalOrderOfOperations3GeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "-",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "*",
		Children: []*GeneNode{},
	})
	root.Children[0].Add(&GeneNode{
		Value:    "33",
		Children: []*GeneNode{},
	})
	root.Children[0].Add(&GeneNode{
		Value:    "4",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "/",
		Children: []*GeneNode{},
	})
	root.Children[1].Add(&GeneNode{
		Value:    "2",
		Children: []*GeneNode{},
	})
	root.Children[1].Add(&GeneNode{
		Value:    "5",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, 131.6)
}

func TestEvalOrderOfOperations5GeneNode(t *testing.T) {
	root := &GeneNode{
		Value:    "|",
		Children: []*GeneNode{},
	}
	root.Add(&GeneNode{
		Value:    "7",
		Children: []*GeneNode{},
	})
	root.Add(&GeneNode{
		Value:    "+",
		Children: []*GeneNode{},
	})
	root.Children[1].Add(&GeneNode{
		Value:    "1.7",
		Children: []*GeneNode{},
	})
	root.Children[1].Add(&GeneNode{
		Value:    "*",
		Children: []*GeneNode{},
	})
	root.Children[1].Children[1].Add(&GeneNode{
		Value:    "3",
		Children: []*GeneNode{},
	})
	root.Children[1].Children[1].Add(&GeneNode{
		Value:    "2.5",
		Children: []*GeneNode{},
	})
	val := root.Eval()
	AssertFloat(t, val, 15.0)
}
