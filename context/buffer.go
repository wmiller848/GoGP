package context

import (
	"fmt"
	"io"
	"sync"
)

func NewBuffer() *Buffer {
	return &Buffer{
		cache:    false,
		open:     true,
		data:     []byte{},
		buffered: make(map[*byte]int),
	}
}

type Buffer struct {
	sync.Mutex
	cache    bool
	open     bool
	data     []byte
	buffered map[*byte]int
}

func (b *Buffer) Clone() *Buffer {
	return &Buffer{
		open:     true,
		data:     []byte(b.data),
		buffered: make(map[*byte]int),
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
	ref := &data[0]
	bufLeng := len(b.data)
	dataLeng := len(data)
	//fmt.Println("CACHE", b.cache, ref)
	//fmt.Println("OPEN", b.open, ref)
	//fmt.Println("PRE", bufLeng, dataLeng, ref)

	index := b.buffered[ref]
	//fmt.Println("INDEX", index, ref)
	if index+dataLeng > bufLeng && bufLeng != 0 {
		dataLeng = bufLeng - index
	}

	//fmt.Println("POST", bufLeng, dataLeng, ref)
	if ((b.open || bufLeng > 0) && dataLeng != 0) || b.cache {
		//fmt.Println("INNER - 1", ref)
		if bufLeng > 0 && index < bufLeng {
			//fmt.Println("INNER - 2", ref)
			if !b.cache {
				copy(data, b.data[:dataLeng])
				b.data = b.data[dataLeng:]
			} else {
				copy(data, b.data[index:index+dataLeng])
				b.buffered[ref] += dataLeng
			}
			b.Unlock()
			return dataLeng, nil
		} else if b.cache && index >= bufLeng && bufLeng != 0 {
			fmt.Println("Read() - closing [inner]")
			b.Unlock()
			return 0, io.EOF
		}
		b.Unlock()
		return 0, nil
	}
	fmt.Println("Read() - closing")
	b.Unlock()
	return 0, io.EOF

	//if dleng > bleng && bleng != 0 {
	//dleng = bleng
	//} else if index+dleng > bleng && index < bleng {
	//dleng = bleng - index
	//}
	//if b.open || bleng > 0 || b.cache {
	//if index+dleng <= bleng {
	//if !b.cache {
	//copy(data, b.data[:dleng])
	//b.data = b.data[dleng:]
	//b.buffered[ref] = 0
	//} else {
	//copy(data, b.data[index:index+dleng])
	//b.buffered[ref] += dleng
	//}
	//} else if index >= bleng {
	//fmt.Println(index, bleng)
	//fmt.Println("Read() - closing (inner)")
	//b.Unlock()
	//return 0, io.EOF
	//}
	//b.Unlock()
	//return dleng, nil
	//}
	//fmt.Println("Read() - closing")
	//b.Unlock()
	//return 0, io.EOF
}

func (b *Buffer) Write(data []byte) (int, error) {
	b.Lock()
	b.data = append(b.data, data...)
	leng := len(data)
	//fmt.Println("Write() dumping", leng)
	b.Unlock()
	return leng, nil
}

func (b *Buffer) Pipe(r io.Reader) (io.Reader, chan []byte) {
	tap := make(chan []byte)
	go func() {
		data := make([]byte, 1024)
		for {
			leng, err := r.Read(data)
			if leng > 0 {
				b.Write(data[:leng])
				tap <- data[:leng]
			}
			fmt.Println("Pipe slurping", leng)
			if err == io.EOF {
				fmt.Println("Pipe() - closing")
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
		data := make([]byte, 1024)
		for {
			leng, err := b.Read(data)
			//fmt.Println(string(data))
			if leng > 0 {
				tap <- data[:leng]
			}
			//fmt.Println("Tap slurping", leng)
			if err == io.EOF {
				fmt.Println("Tap() - closing")
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

type Multiplexer struct {
	well *Buffer
}

func Multiplex(r io.Reader) *Multiplexer {
	m := &Multiplexer{
		well: NewBuffer(),
	}
	m.well.cache = true
	go func() {
		data := make([]byte, 1024)
		for {
			leng, err := r.Read(data)
			if leng > 0 {
				m.well.Write(data[:leng])
			}
			fmt.Println("Multiplexer slurping from well", leng)
			if err == io.EOF {
				fmt.Println("Multiplexer() - closing")
				m.well.Close()
				break
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
			if leng > 0 {
				p.Write(data[:leng])
			}
			fmt.Println("Multiplex slurping", leng)
			if err == io.EOF {
				fmt.Println("Multiplex() - closing")
				p.Close()
				break
			}
		}
	}()
	return p
}

func (m *Multiplexer) Destroy() {
	m.well.Close()
}
