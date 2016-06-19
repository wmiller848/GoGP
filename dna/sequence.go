package dna

import "fmt"

type Sequence struct {
	Codex    Codex
	CodexID  int
	Index    int
	Elements int
}

func (s *Sequence) Node() *SequenceNode {
	return &SequenceNode{
		Sequence: s,
		Child:    nil,
	}
}

type SequenceNode struct {
	*Sequence
	Child *SequenceNode
}

func (s *SequenceNode) String() string {
	str := fmt.Sprintf("%+v ", *s.Sequence)
	if s.Child != nil {
		str += s.Child.String()
	}
	return str
}

func (s *SequenceNode) Bytes() []byte {
	bytes := []byte{}
	bytes = append(bytes, []byte(s.Codex.String())...)
	if s.Child != nil {
		bytes = append(bytes, s.Child.Bytes()...)
	}
	return bytes
}

func (s *SequenceNode) Clone() *SequenceNode {
	seq := &SequenceNode{
		Sequence: &Sequence{
			Codex:    s.Sequence.Codex,
			CodexID:  s.Sequence.CodexID,
			Index:    s.Sequence.Index,
			Elements: s.Sequence.Elements,
		},
	}
	if s.Child != nil {
		seq.Child = s.Child.Clone()
	}
	return seq
}

func (s *SequenceNode) Merge(seq *Sequence) *SequenceNode {
	node := seq.Node()
	if s.Index > seq.Index {
		node.Child = s
		return node
	} else {
		if s.Child == nil {
			s.Child = node
		} else {
			s.Child = s.Child.Merge(seq)
		}
		return s
	}
}
