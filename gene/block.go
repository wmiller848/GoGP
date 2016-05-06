package gene

import (
	"github.com/wmiller848/GoGP/util"
)

func randomOperator() byte {
	switch util.RandomNumber(0, 3) {
	case 0:
		return byte('+')
	case 1:
		return byte('-')
	case 2:
		return byte('*')
	case 3:
		return byte('/')
	default:
		return randomOperator()
	}
}

func randomVariable() byte {
	return byte(util.RandomNumber(0, 9))
}

func randomNumber() byte {
	return byte(util.RandomNumber(0, 9))
}

const (
	StageSpawn  int = 1
	StageAlive  int = 2
	StageDieing int = 3
	StageDead   int = 3
)

type BuildingBlock interface {
	Stages() []int
	// Bytes() [][]byte
	// Distribution() (map[string]int, int)
	Random(int, int) Gene
}
