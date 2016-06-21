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
	varA := Variable(0)
	AssertStr(t, varA, "$a")

	varB := Variable(9)
	AssertStr(t, varB, "$j")

	varC := Variable(25)
	AssertStr(t, varC, "$z")

	varD := Variable(26)
	AssertStr(t, varD, "$aa")
}

func TestVariableLookup(t *testing.T) {
	varA := Variable(0)
	AssertInt(t, VariableLookup(varA), 0)

	varB := Variable(9)
	AssertInt(t, VariableLookup(varB), 9)

	varC := Variable(25)
	AssertInt(t, VariableLookup(varC), 25)

	varD := Variable(26)
	AssertInt(t, VariableLookup(varD), 26)
}
