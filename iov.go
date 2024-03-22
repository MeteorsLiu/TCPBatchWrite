package tcpbatch

import (
	"sort"

	"golang.org/x/sys/unix"
)

const (
	IovMax = 1024
)

type iovLen struct {
	idx, max int
	len      []int
	iov      []unix.Iovec
	bufs     [][]byte
}

func (i *iovLen) reset() {
	i.idx, i.max = 0, 0
	i.len, i.iov, i.bufs = i.len[:0], i.iov[:0], i.bufs[:0]
}

func (iov *iovLen) resize(start int) {
	lenIdx, bufs := iov.len[iov.idx:], iov.bufs[iov.idx:]
	if lenIdx[0] > start {
		base := iov.idx
		nb := len(bufs[0])
		fromLen := lenIdx[0] - start
		iov.iov[base].Base = &bufs[0][nb-fromLen]
		iov.iov[base].SetLen(fromLen)
	} else {
		i := sort.Search(len(lenIdx), func(i int) bool {
			return lenIdx[i] > start
		})
		base := i + iov.idx
		nb := len(bufs[i])
		fromLen := lenIdx[i] - start
		iov.iov[base].Base = &bufs[i][nb-fromLen]
		iov.iov[base].SetLen(fromLen)
		iov.idx = base
	}
}

func (i *iovLen) iovec() []unix.Iovec {
	return i.iov[i.idx:]
}

func (i *iovLen) append(buf []byte) {
	i.max += len(buf)
	i.len = append(i.len, i.max)
	i.bufs = append(i.bufs, buf)
	pos := len(i.len) - 1
	i.iov = i.iov[:len(i.len)]
	i.iov[pos].Base = &buf[0]
	i.iov[pos].SetLen(len(buf))
}

func newIovLen() *iovLen {
	return &iovLen{
		len:  make([]int, 0, IovMax),
		iov:  make([]unix.Iovec, 0, IovMax),
		bufs: make([][]byte, 0, IovMax),
	}
}
