package test

import (
	"io"
	"sync"
)

// implements io.RWC with buffers, to help write tests
type BufferedRWC struct {
	sync.Mutex
	cond sync.Cond

	readClosed bool
	ToRead     [][]byte

	Written [][]byte
	Closed  bool
}

func NewBufferedRWC() *BufferedRWC {
	this := &BufferedRWC{}
	this.cond.L = &this.Mutex
	return this
}

var _ io.ReadWriteCloser = (*BufferedRWC)(nil)

// append to the Read queue, panics if the buffer is closed
func (this *BufferedRWC) Append(p []byte) {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	if this.readClosed {
		panic(io.ErrClosedPipe)
	}
	this.ToRead = append(this.ToRead, p)
	this.cond.Signal()
}

func (this *BufferedRWC) IsReadClosed() bool {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	return this.readClosed
}

func (this *BufferedRWC) SetReadClosed() {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	this.readClosed = true
	this.cond.Broadcast()
}

// RWC interface
func (this *BufferedRWC) Read(dst []byte) (int, error) {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()

	if cap(dst) == 0 {
		if len(this.ToRead) == 0 && this.readClosed {
			return 0, io.EOF
		}
		return 0, nil
	}

	for {
		if len(this.ToRead) > 0 {
			ct := copy(dst, this.ToRead[0])
			if ct < len(this.ToRead[0]) {
				this.ToRead[0] = this.ToRead[0][ct:]
			} else {
				this.ToRead = this.ToRead[1:]
			}
			return ct, nil
		} else if this.readClosed {
			return 0, io.EOF
		} else {
			this.cond.Wait()
		}
	}
}

func (this *BufferedRWC) HasRead() bool {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	return len(this.ToRead) > 0
}

func (this *BufferedRWC) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil // NOOP
	}
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	if this.readClosed {
		return 0, io.ErrClosedPipe
	}
	this.Written = append(this.Written, b)
	return len(b), nil
}

func (this *BufferedRWC) Close() error {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	if this.readClosed {
		return io.ErrClosedPipe
	}
	this.readClosed = true
	return nil
}
