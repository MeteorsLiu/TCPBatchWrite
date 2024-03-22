package tcpbatch

import (
	"syscall"

	_ "unsafe"

	"golang.org/x/sys/unix"
	_ "golang.org/x/sys/unix"
)

//go:linkname writev golang.org/x/sys/unix.writev
func writev(fd int, iovs []unix.Iovec) (n int, err error)

func Writev(dst syscall.RawConn, iov *iovLen, onEagain func()) (n int, err error) {
	var max, ns int
	// always escape, so we can only use the pool.
	max = iov.max
	err2 := dst.Write(func(fd uintptr) bool {
		f := int(fd)
		for max > 0 {
			ns, err = writev(f, iov.iovec())
			if ns > 0 {
				max -= ns
				n += ns
				// resize before because when TCP_NOT_SENT_LOWAT enables
				// will cause wrong slices.
				if max > 0 {
					iov.resize(n)
				}
			}

			switch err {
			case syscall.EINTR:
				continue
			case syscall.EAGAIN:
				if onEagain != nil {
					onEagain()
				}
				return false
			}

			if err != nil {
				break
			}
		}
		return true
	})

	if err == nil && err2 != nil {
		err = err2
	}
	return
}
