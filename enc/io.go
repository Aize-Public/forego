package enc

import (
	"io"

	"github.com/Aize-Public/forego/ctx"
)

type Writer interface {
	Write(c ctx.C, n Node) error
}

type WriteCloser interface {
	io.Closer
	Writer
}

type Reader interface {
	Read(c ctx.C) (Node, error)
}

type ReaderCloser interface {
	Reader
	io.Closer
}

type ReadWriter interface {
	Reader
	Writer
}

type ReadWriteCloser interface {
	Reader
	Writer
	io.Closer
}
