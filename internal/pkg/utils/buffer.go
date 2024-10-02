package util

import "sync"

type Buffer[T any] struct {
	Data []T
}

func (b *Buffer[T]) Reset(length int) {
	if cap(b.Data) < length {
		b.Data = make([]T, length)
	}
	b.Data = b.Data[:length]
}

type BufferPool[T any] struct {
	pool sync.Pool
}

func NewBufferPool[T any](initcap int) *BufferPool[T] {
	return &BufferPool[T]{
		sync.Pool{New: func() interface{} {
			return &Buffer[T]{
				Data: make([]T, 0, initcap),
			}
		}},
	}
}

func (p *BufferPool[T]) Get() *Buffer[T] {
	buf := p.pool.Get().(*Buffer[T])
	buf.Data = buf.Data[:0]
	return buf
}

func (p *BufferPool[T]) Put(s *Buffer[T]) {
	p.pool.Put(s)
}
