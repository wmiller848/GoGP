package gene

import "testing"

func AssertInt(t *testing.T, a, b int) bool {
	if a != b {
		t.Error("Expected", b)
		t.Error("Got:", a)
		return false
	}
	return true
}

func TestVariable(t *testing.T) {
	varA0 := Variable(0)
	AssertStr(t, varA0, "$zz")

	varB := Variable(9)
	AssertStr(t, varB, "$zq")

	varC := Variable(25)
	AssertStr(t, varC, "$yq")

	varD := Variable(26)
	AssertStr(t, varD, "$ya")

	varE := Variable(27)
	AssertStr(t, varE, "$yb")
}

func TestVariableLookup(t *testing.T) {
	varA := "$zz"
	AssertInt(t, VariableLookup(varA), 0)

	varB := "$zq"
	AssertInt(t, VariableLookup(varB), 9)

	varC := "$yq"
	AssertInt(t, VariableLookup(varC), 25)

	varD := "$ya"
	AssertInt(t, VariableLookup(varD), 26)

	varE := "$yb"
	AssertInt(t, VariableLookup(varE), 27)
}
