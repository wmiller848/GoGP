package gene

import (
	_ "github.com/wmiller848/GoGP/util"
)

const (
	CursorNil       int = 0
	CursorVariable  int = 1
	CursorNumber    int = 2
	CursorOperator  int = 3
	CursorSeparator int = 4
)

type Gene interface {
	Eq(Gene) bool
	Clone() []byte
	Heal() []byte
	LastChrome(int) int
	NextChrome(int) int
	Len() int
	At(int) byte
	MarshalTree() (*GeneNode, error)
}

type GenericGene []byte
