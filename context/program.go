package context

import "github.com/wmiller848/GoGP/program"

type ProgramInstance struct {
	*program.Program
	ID         string
	Generation int
	Score      float64
	Group      map[float64]*Group
}

type Programs []*ProgramInstance

func (p Programs) Len() int           { return len(p) }
func (p Programs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Programs) Less(i, j int) bool { return p[i].Score < p[j].Score }
