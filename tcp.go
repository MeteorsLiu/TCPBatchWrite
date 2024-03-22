package tcpbatch

import (
	"net"
	"sync"
	"syscall"
)

type BatchConn interface {
	net.Conn
	RawConn() syscall.RawConn
	WriteBuffer([]byte)
	Submit(onEagain func()) (int, error)
	SetOnEAGAIN(fn func())
}

var (
	iovPool = sync.Pool{
		New: func() any {
			return newIovLen()
		},
	}
)

type tcpBatch struct {
	net.Conn
	iovLen   *iovLen
	raw      syscall.RawConn
	onEagain func()
}

func NewTCPBatch(c net.Conn) BatchConn {
	raw, _ := c.(*net.TCPConn).SyscallConn()
	iov := iovPool.Get().(*iovLen)
	iov.reset()
	t := &tcpBatch{
		Conn:   c,
		raw:    raw,
		iovLen: iov,
	}
	return t
}

func (t *tcpBatch) WriteBuffer(b []byte) {
	t.iovLen.append(b)
	return
}

func (t *tcpBatch) Submit(onEagain func()) (n int, err error) {
	if onEagain == nil {
		n, err = Writev(t.raw, t.iovLen, t.onEagain)
	} else {
		n, err = Writev(t.raw, t.iovLen, onEagain)
	}
	t.iovLen.reset()
	return n, err
}

func (t *tcpBatch) RawConn() syscall.RawConn {
	return t.raw
}

func (t *tcpBatch) Close() error {
	iovPool.Put(t.iovLen)
	return t.Conn.Close()
}

func (t *tcpBatch) SetOnEAGAIN(fn func()) {
	t.onEagain = fn
}
