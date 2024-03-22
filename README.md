# TCP Batch Write

[Why](https://zhuanlan.zhihu.com/p/673532129)

# Example
```
batch := tcpbatch.NewTCPBatch(TCPConn)

// replace
TCPConn.Write(buf)
// to
batch.WriteBuffer(buf)
n, err := batch.Submit(nil)
```

