package gene

import (
	"testing"
)

func AssertStr(t *testing.T, a, b string) bool {
	if a != b {
		t.Error("Expected " + b)
		t.Error("Got:" + a)
		return false
	}
	return true
}

func AssertGene(g Gene) {}

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
