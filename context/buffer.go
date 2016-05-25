package context

import (
	"io"
	"sync"
)

func NewBuffer() *Buffer {
	return &Buffer{
		open: true,
		data: []byte{},
	}
}

type Buffer struct {
	sync.Mutex
	open bool
	data []byte
}

func (b *Buffer) Open() {
	b.open = true
}
func (b *Buffer) Close() {
	b.open = false
}

func (b *Buffer) Read(data []byte) (int, error) {
	bleng := len(b.data)
	dleng := len(data)
	if b.open == true || bleng > 0 {
		var leng int
		if dleng > bleng {
			leng = bleng
		} else {
			leng = dleng
		}
		if leng > 0 {
			b.Lock()
			copy(data, b.data[:leng])
			b.data = b.data[leng:]
			b.Unlock()
		}
		return leng, nil
	}
	return 0, io.EOF
}

func (b *Buffer) Write(data []byte) (int, error) {
	b.Lock()
	b.data = append(b.data, data...)
	b.Unlock()
	leng := len(data)
	return leng, nil
}

func (b *Buffer) Data() []byte {
	return b.data
}
