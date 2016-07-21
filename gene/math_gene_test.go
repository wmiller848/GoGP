package gene

import (
	"testing"
)

func TestMathGeneIsGene(t *testing.T) {
	g := MathGene("*{+$az,10,20}{-$ay,25,33}").Heal()
	AssertGene(g)
}

func TestHealthyMathGene(t *testing.T) {
	g := MathGene("*{+$az,10,20}{-$ay,25,33}").Heal().Bytes()
	assert := MathGene("*{+$az,10,20}{-$ay,25,33}").Bytes()
	AssertStr(t, string(g), string(assert))
}

func TestUnhealthyMathGene(t *testing.T) {
	g := MathGene("&*{+$az,,,10,20}{-^*{$ay,25,33}").Heal().Bytes()
	assert := MathGene("&{+$az,10,20}{-{$ay,25,33}}").Bytes()
	AssertStr(t, string(g), string(assert))
}

func TestUnhealthyMathGeneOperators(t *testing.T) {
	g := MathGene("&*{+$az,10,20}{-^*$ay,25,33}").Heal().Bytes()
	assert := MathGene("&{+$az,10,20}{-$ay,25,33}").Bytes()
	AssertStr(t, string(g), string(assert))
}

func TestUnhealthyMathGeneSeperators(t *testing.T) {
	g := MathGene("&,,{+$az,,,10,20$c}{-$ay,,,25,33},,,").Heal().Bytes()
	assert := MathGene("&{+$az,10,20,$c}{-$ay,25,33}").Bytes()
	AssertStr(t, string(g), string(assert))
}

func TestUnhealthyMathGeneScope(t *testing.T) {
	g := MathGene("*{+$az,^&+{10,20}{*$ay,{{25,33}").Heal().Bytes()
	assert := MathGene("*{+$az,{10,20}{*$ay,{{25,33}}}}").Bytes()
	AssertStr(t, string(g), string(assert))
}

func TestMarshalTree1MathGene(t *testing.T) {
	g := MathGene("*10,3-10")
	root, err := g.MarshalTree()
	if err != nil {
		t.Error(err.Error())
	}
	AssertStr(t, root.Value, "*")
	AssertInt(t, len(root.Children), 3)
}

func TestMarshalTree2MathGene(t *testing.T) {
	g := MathGene("+{-20,10}{*50,10/2,3}")
	root, err := g.MarshalTree()
	if err != nil {
		t.Error(err.Error())
	}
	AssertStr(t, root.Value, "+")
	AssertInt(t, len(root.Children), 2)
}
