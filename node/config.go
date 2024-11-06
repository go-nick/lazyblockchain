package node

import (
	"context"

	"github.com/decred/dcrd/rpcclient/v8"
)

// RPC ...
type RPC struct {
	Client *rpcclient.Client
	ctx    context.Context
}

// ConnectRPC ...
func ConnectRPC() (*RPC, error) {
	// Configure RPC client options
	connCfg := &rpcclient.ConnConfig{
		Host:         "192.168.1.10:8332",
		User:         "lobo",
		Pass:         "123456",
		HTTPPostMode: true, // Bitcoin Core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin Core does not use TLS by default
	}

	// Initialize client instance
	rpcClient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}
	rpc := &RPC{
		Client: rpcClient,
		ctx:    context.Background(),
	}

	return rpc, nil
}
