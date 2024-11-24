# LazyBlockchain

### Quick Run

at: `node/rpc.go`

Change this for your target node:
```
	connCfg := &rpcclient.ConnConfig{
		Host:         "192.168.1.10:8332",
		User:         "lobo",
		Pass:         "123456",
		HTTPPostMode: true, // Bitcoin Core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin Core does not use TLS by default
	}
```

At the root of the project:
```
go run .
```
