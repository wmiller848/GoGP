package context

import (
	_ "fmt"
	"io"
	"sync"
)

func NewBuffer() *Buffer {
	return &Buffer{
		open:     true,
		data:     []byte{},
		buffered: make(map[*[]byte]int),
	}
}

type Buffer struct {
	sync.Mutex
	cache    bool
	open     bool
	data     []byte
	buffered map[*[]byte]int
}

func (b *Buffer) Clone() *Buffer {
	return &Buffer{
		open:     true,
		data:     []byte(b.data),
		buffered: make(map[*[]byte]int),
	}
}

func (b *Buffer) Open() {
	b.open = true
}
func (b *Buffer) Close() {
	b.open = false
}

func (b *Buffer) Read(data []byte) (int, error) {
	b.Lock()
	bleng := len(b.data)
	dleng := len(data)
	index := b.buffered[&data]
	if b.open == true || bleng > 0 {
		var leng int
		if dleng > bleng {
			leng = bleng
		} else {
			leng = dleng
		}
		if leng > 0 && index < leng {
			copy(data, b.data[index:index+leng])
			if !b.cache {
				b.data = b.data[leng:]
			} else {
				b.buffered[&data] += leng
			}
		} else if index > leng && b.open == false {
			return 0, io.EOF
		}
		b.Unlock()
		return leng, nil
	}
	b.Unlock()
	return 0, io.EOF
}

func (b *Buffer) Write(data []byte) (int, error) {
	b.Lock()
	b.data = append(b.data, data...)
	leng := len(data)
	b.Unlock()
	return leng, nil
}

func (b *Buffer) Pipe(r io.Reader) (io.Reader, chan []byte) {
	tap := make(chan []byte)
	go func() {
		for {
			data := make([]byte, 1024)
			leng, err := r.Read(data)
			if err == io.EOF {
				b.Close()
				close(tap)
				return
			}
			if leng > 0 {
				b.Write(data)
				tap <- data
			} else {
				b.Close()
				close(tap)
				return
			}
		}
	}()
	return b, tap
}

func (b *Buffer) Tap() chan []byte {
	tap := make(chan []byte)
	go func() {
		for {
			data := make([]byte, 1024)
			leng, err := b.Read(data)
			if err == io.EOF {
				close(tap)
				return
			}
			if leng > 0 {
				tap <- data
			} else {
				close(tap)
				return
			}
		}
	}()
	return tap
}

type Multiplexer struct {
	well *Buffer
}

func Multiplex(r io.Reader) *Multiplexer {
	m := &Multiplexer{
		well: NewBuffer(),
	}
	m.well.cache = true
	go func() {
		for {
			data := make([]byte, 1024)
			leng, err := r.Read(data)
			if err == io.EOF {
				m.well.Close()
				return
			}
			if leng > 0 {
				m.well.Write(data)
			} else {
				return
			}
		}
	}()
	return m
}

func (m *Multiplexer) Multiplex() *Buffer {
	p := NewBuffer()
	go func() {
		data := make([]byte, 1024)
		for {
			leng, err := m.well.Read(data)
			if err == io.EOF {
				p.Close()
				return
			}
			if leng > 0 {
				p.Write(data)
			} else {
				return
			}
		}
	}()
	return p
}

func (b *Buffer) Data() []byte {
	return b.data
}
