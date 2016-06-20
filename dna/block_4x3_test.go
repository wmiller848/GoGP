package dna

import "testing"

func AssertBlock(b Block) {}

func AssertStr(t *testing.T, a, b string) bool {
	if a != b {
		t.Error("Expected: " + b)
		t.Error("Got: " + a)
		return false
	}
	return true
}

func AssertBase(t *testing.T, a, b Base) bool {
	if a != b {
		t.Error("Expected " + string(b))
		t.Error("Got:" + string(a))
		return false
	}
	return true
}

func InitBlock(t *testing.T) *Block4x3 {
	codex := Codex{Codon("a"), Codon("b")}
	blk, err := NewBlock4x3(Block4x3Bases, codex)
	if err != nil {
		t.Error(err.Error())
		return nil
	}
	return blk
}

func TestBlock4x3IsBlock(t *testing.T) {
	blk := InitBlock(t)
	AssertBlock(blk)
}

func TestBlock4x3Random(t *testing.T) {
	blk := InitBlock(t)
	dna := blk.Random()
	if dna == nil {
		t.Error("Random() failed to create DNA")
	}
}

func TestBlock4x3Match(t *testing.T) {
	blk := InitBlock(t)
	base0 := blk.Match(0x10)
	AssertBase(t, base0, 0x00)

	base1 := blk.Match(0x50)
	AssertBase(t, base1, 0x40)

	base2 := blk.Match(0x90)
	AssertBase(t, base2, 0x80)

	base3 := blk.Match(0xd0)
	AssertBase(t, base3, 0xc0)
}

func TestBlock4x3Decode(t *testing.T) {
	blk := InitBlock(t)
	codonStart, err := blk.Decode([]Base{0x00, 0x00, 0x00})
	if err != nil {
		t.Error(err.Error())
		return
	}
	AssertStr(t, string(codonStart), "<")

	codonA, err := blk.Decode([]Base{0x40, 0x00, 0x00})
	if err != nil {
		t.Error(err.Error())
		return
	}
	AssertStr(t, string(codonA), "a")

	codonB, err := blk.Decode([]Base{0x80, 0x00, 0x00})
	if err != nil {
		t.Error(err.Error())
		return
	}
	AssertStr(t, string(codonB), "b")

	codonEnd, err := blk.Decode([]Base{0xff, 0x00, 0x00})
	if err != nil {
		t.Error(err.Error())
		return
	}
	AssertStr(t, string(codonEnd), ">")
}
