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
	b.Lock()
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
			copy(data, b.data[:leng])
			b.data = b.data[leng:]
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

func (b *Buffer) Tap(r io.Reader) (io.Reader, chan []byte) {
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

func (b *Buffer) Pipe() chan []byte {
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

func (b *Buffer) Data() []byte {
	return b.data
}
