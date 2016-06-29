package context

import (
	"time"

	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

type PaneType uint

const (
	Horizontal PaneType = 0
	Vertical   PaneType = 1
	Cross      PaneType = 2
)

type Position struct {
	x float64
	y float64
}

type Size struct {
	width  float64
	height float64
}

type Node struct {
	*Position
	*Size
	value    []byte
	children []*Node
	root     bool
}

func (n *Node) draw() {
	if len(n.children) > 0 {
		for i, _ := range n.children {
			n.children[i].draw()
		}
	} else {
		tWidth, tHeight := termbox.Size()
		w := int(n.width * float64(tWidth))
		h := int(n.height * float64(tHeight))
		x := int(n.x)
		y := int(n.y)
		if x > w-1 {
			w = w - 1
		}
		if y > y-1 {
			h = h - 1
		}
		var q, j, t int
		for t < h+1 {
			if q < len(n.value) {
				if n.value[q] == '\n' {
					j = 0
					t++
				} else if n.value[q] == '\r' {
					j = 0
				} else {
					termbox.SetCell(x+j, y+t, rune(n.value[q]), coldef, coldef)
					j++
				}
				q++
			} else {
				termbox.SetCell(x+j, y+t, rune(' '), coldef, coldef)
				j++
			}
			if j > w {
				j = 0
				t++
			}
		}
	}
}

func (n *Node) splitHorizontal() []*Node {
	if len(n.children) == 0 {
	}
	return n.children
}

func (n *Node) splitVertical() []*Node {
	w, _ := termbox.Size()
	if len(n.children) == 0 {
		n.children = []*Node{
			&Node{
				Position: &Position{
					x: n.x,
					y: n.y,
				},
				Size: &Size{
					width:  n.width / 2.0,
					height: n.height,
				},
				value: []byte(""),
			},
			&Node{
				Position: &Position{
					x: n.x + (float64(w) * n.width / 2.0),
					y: n.y,
				},
				Size: &Size{
					width:  n.width / 2.0,
					height: n.height,
				},
				value: []byte(""),
			},
		}
	}
	return n.children
}

func (n *Node) split() []*Node {
	if len(n.children) == 0 {

	}
	return n.children
}

type Terminal struct {
	window *Node
}

func (t *Terminal) AddPane(n Node, typ PaneType) {
	if typ == Horizontal {
		n.splitHorizontal()
	} else if typ == Vertical {
		n.splitVertical()
	} else if typ == Cross {
		n.split()
	}
}

func (t *Terminal) Start(ctx *Context) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	t.window = &Node{
		root: true,
		Size: &Size{
			width:  1.0,
			height: 1.0,
		},
		Position: &Position{
			x: 0.0, // float64(w) / 2.0,
			y: 0.0, // float64(h) / 2.0,
		},
		value: []byte(""),
	}
	// t.window.splitVertical()
	// t.window.value = []byte("LLama")
	t.draw()
	go func() {
		for {
			t.draw()
			time.Sleep((1000 / 5) * time.Millisecond)
		}
	}()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				return
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				//edit_box.MoveCursorOneRuneBackward()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				//edit_box.MoveCursorOneRuneForward()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				//edit_box.DeleteRuneBackward()
			case termbox.KeyDelete, termbox.KeyCtrlD:
				//edit_box.DeleteRuneForward()
			case termbox.KeyTab:
				//edit_box.InsertRune('\t')
			case termbox.KeySpace:
				//edit_box.InsertRune(' ')
			case termbox.KeyCtrlK:
				//edit_box.DeleteTheRestOfTheLine()
			case termbox.KeyHome, termbox.KeyCtrlA:
				//edit_box.MoveCursorToBeginningOfTheLine()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				//edit_box.MoveCursorToEndOfTheLine()
			default:
				if ev.Ch != 0 {
					//edit_box.InsertRune(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		//redraw_all()
	}
}

//func (t *Terminal) word(pos *Position, str string) {
//w, h := termbox.Size()
//strBytes := []byte(str)
//m := len(strByts)
//termbox.SetCell(pos.X, pos.Y, 'X', coldef, coldef)
//}

func (t *Terminal) draw() {
	termbox.Clear(coldef, coldef)
	t.window.draw()
	//w, h := termbox.Size()
	//x := w / 2
	//y := h / 2
	//termbox.SetCell(x, y, 'X', coldef, coldef)
	termbox.Flush()
}
