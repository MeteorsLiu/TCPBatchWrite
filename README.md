# TCP Batch Write

[Why](https://zhuanlan.zhihu.com/p/673532129)

# Example
```
batch := tcpbatch.NewTCPBatch(TCPConn)

// replace
TCPConn.Write(buf1)
TCPConn.Write(buf2)
...
TCPConn.Write(bufn)

// to
batch.WriteBuffer(buf1)
batch.WriteBuffer(buf2)
batch.WriteBuffer(buf3)
...
batch.WriteBuffer(bufn)
// max size 1024
n, err := batch.Submit(nil)
```

