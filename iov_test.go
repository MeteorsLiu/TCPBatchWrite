package tcpbatch

import (
	"math/rand"
	"testing"
)

func generateBuffers() [][]byte {
	b := make([][]byte, rand.Intn(11)+1)
	for i := range b {
		b[i] = make([]byte, rand.Intn(11)+1)
	}
	return b
}

func TestIov(t *testing.T) {
	bufs := [][]byte{
		[]byte{1, 2, 3, 4},
		[]byte{5, 6, 7, 8, 9, 10, 11},
		[]byte{12, 13},
		[]byte{14, 15, 16, 17, 18, 19},
		[]byte{20},
	}
	iov := newIovLen()

	for _, b := range bufs {
		iov.append(b)
	}

	iov.resize(3)
	if iov.iovec()[0].Base != &bufs[0][3] || iov.iovec()[0].Len != 1 {
		t.Error("misbehave: 3")
		return
	}

	iov.resize(5)

	if iov.iovec()[0].Base != &bufs[1][1] || iov.iovec()[0].Len != 6 {
		t.Error("misbehave: 5")
		return
	}

	iov.resize(10)
	if iov.iovec()[0].Base != &bufs[1][6] || iov.iovec()[0].Len != 1 {
		t.Error("misbehave: 10")
		return
	}

	iov.resize(11)
	if iov.iovec()[0].Base != &bufs[2][0] || iov.iovec()[0].Len != 2 {
		t.Error("misbehave: 11")
		return
	}

	iov.resize(14)
	if iov.iovec()[0].Base != &bufs[3][1] || iov.iovec()[0].Len != 5 {
		t.Error("misbehave: 14")
		return
	}

	iov.resize(19)
	if iov.iovec()[0].Base != &bufs[4][0] || iov.iovec()[0].Len != 1 {
		t.Error("misbehave: 14")
		return
	}
}

func TestBuf(t *testing.T) {
	buf := newIovLen()

	bufs := generateBuffers()

	for _, b := range bufs {
		buf.append(b)
	}

	t.Log(buf.max, buf.len)

	for i, b := range buf.bufs {
		t.Log(i, "len: ", len(b))
	}

	start := rand.Intn(buf.max)

	buf.resize(start)

	t.Log(start, buf.max, buf.len)

	for i, b := range buf.bufs {
		t.Log(i, "len: ", len(b))
	}

	buf.reset()

	t.Log(buf)
}
