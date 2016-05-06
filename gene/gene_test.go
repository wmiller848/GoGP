package gene

import (
	// "github.com/wmiller848/GoGP/gene"
	"testing"
)

//	*10,9

//	*
//	|___
//	10	9

//	(10) * (9)

// *10,20,10
//	10 * (20 * 10)
// *{10,30}{-10,18}

func TestGene(t *testing.T) {
	// t.Error("Need to be cooler")
	g := Gene("*{+$az,10,20}{-$ay,25,33}")
	g = g.Heal()
	if !g.Eq(Gene("*{+$az,10,20}{-$ay,25,33}")) {
		t.Error("Expected *{+$az,10,20}{-$ay,25,33}")
		t.Error("Got:" + string(g))
	}
	root, _ := g.MarshalTree()
	exp, _ := root.MarshalExpression()
	if string(exp) != "($az+10+20)*($ay-25-33)" {
		t.Error("Expected ($az+10+20)*($ay-25-33)")
		t.Error("Got:" + string(exp))
	}
}

func TestGeneHealWithOperators(t *testing.T) {
	g := Gene("/*//{+$az,10*10,20}*/-7{-$ay,25}")
	g = g.Heal()
	if !g.Eq(Gene("/{+$az,10*10,20}-7{-$ay,25}")) {
		t.Error("Expected /{+$az,10*10,20}-7{-$ay,25}")
		t.Error("Got:" + string(g))
	}
	root, _ := g.MarshalTree()
	exp, _ := root.MarshalExpression()
	if string(exp) != "($az+10+(10*20))/(7-($ay-25))" {
		t.Error("Expected ($az+10+(10*20))/(7-($ay-25))")
		t.Error("Got:" + string(exp))
	}
}

func TestGeneHealWithSeperators(t *testing.T) {
	g := Gene("*{+,,$az,10,,,*,,10,20},,,{-,,,$ay,,,,25}")
	g = g.Heal()
	if !g.Eq(Gene("*{+$az,10*10,20}{-$ay,25}")) {
		t.Error("Expected *{+$az,10*10,20}{-$ay,25}")
		t.Error("Got:" + string(g))
	}
	root, _ := g.MarshalTree()
	exp, _ := root.MarshalExpression()
	if string(exp) != "($az+10+(10*20))*($ay-25)" {
		t.Error("Expected ($az+10+(10*20))*($ay-25)")
		t.Error("Got:" + string(exp))
	}
}

func TestGeneHealWithDeepNested(t *testing.T) {
	g := Gene("*$ay{+99*10,-*,-11,12,{-85,,12}{*-+/11,{*123,$az}{-*+-$az,12}},,,9780,,}{-,,*-+,71,72}")
	g = g.Heal()
	if !g.Eq(Gene("*$ay{+99*10,11,12,{-85,12}{/11,{*123,$az}{-$az,12}}9780}{-71,72}")) {
		t.Error("Expected *$ay{+99*10,11,12,{-85,12}{/11,{*123,$az}{-$az,12}}9780}{-71,72}")
		t.Error("Got:" + string(g))
	}
	root, _ := g.MarshalTree()
	exp, _ := root.MarshalExpression()
	if string(exp) != "$ay*(99+(10*11*12*(85-12)*(11/(123*$az)/($az-12))*9780))*(71-72)" {
		t.Error("Expected $ay*(99+(10*11*12*(85-12)*(11/(123*$az)/($az-12))*9780))*(71-72)")
		t.Error("Got:" + string(exp))
	}
}
