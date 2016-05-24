package context

import (
	"io"
	"sync"
)

type Buffer []byte
type BufferMeta struct {
	sync.Mutex
	Done       bool
	AmountRead int
}

var bufferMap map[*Buffer]*BufferMeta = make(map[*Buffer]*BufferMeta)

func NewBuffer() *Buffer {
	buff := make(Buffer, 0)
	bufferMap[&buff] = &BufferMeta{
		Done:       false,
		AmountRead: 0,
	}
	return &buff
}

func (b *Buffer) Read(data []byte) (int, error) {
	bufferMap[b].Lock()
	if bufferMap[b].Done == false {
		bufView := *b
		bleng := len(bufView)
		dleng := len(data)
		var leng int
		if dleng > bleng {
			leng = bleng
		} else {
			leng = dleng
		}
		copy(data, bufView[:leng])
		bufferMap[b].Done = true
		if bufferMap[b].AmountRead >= bleng {
			bufferMap[b].Done = true
		}
		bufferMap[b].Unlock()
		return leng, nil
	}
	bufferMap[b].Unlock()
	delete(bufferMap, b)
	return 0, io.EOF
}
