package context

import (
	"io"
)

func NewBuffer() *Buffer {
	return &Buffer{
		done:          false,
		amountWritten: 0,
		amountRead:    0,
		data:          []byte{},
	}
}

type Buffer struct {
	done          bool
	amountWritten int
	amountRead    int
	data          []byte
}

func (b *Buffer) Read(data []byte) (int, error) {
	if b.done == false {
		bufView := b.data
		bleng := len(bufView)
		dleng := len(data)
		var leng int
		if dleng > bleng {
			leng = bleng
		} else {
			leng = dleng
		}
		copy(data[b.amountRead:], bufView[:leng])
		b.amountRead += leng
		if b.amountRead >= bleng {
			b.done = true
		}
		return leng, nil
	}
	return 0, io.EOF
}

func (b *Buffer) Write(data []byte) (int, error) {
	b.data = append(b.data, data...)
	leng := len(data)
	b.amountWritten += leng
	return leng, nil
}

func (b *Buffer) Data() []byte {
	return b.data
}
