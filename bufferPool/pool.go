package bufferPool

import (
	"bytes"
	"sync"
)

type Pool struct {
	p sync.Pool
}

func New(sz int) *Pool {
	return &Pool{
		p: sync.Pool{
			New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, sz)) },
		},
	}
}

func (p *Pool) Get() *bytes.Buffer {
	return p.p.Get().(*bytes.Buffer)
}

func (p *Pool) Put(buf *bytes.Buffer) {
	buf.Reset()
	p.p.Put(buf)
}
