package node

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/decred/dcrd/rpcclient/v8"
)

// RPCCommandFunc TODO: Doc this
type RPCCommandFunc func(args ...string) (map[string]interface{}, error)

// RPC ...
type RPC struct {
	CmdCache map[string]RPCCommandFunc
	Client   *rpcclient.Client
	ctx      context.Context
}

// ConnectRPC ...
func ConnectRPC(host, port, user, password string) (*RPC, error) {
	var err error
	var config map[string]string

	if host == "" || user == "" || password == "" {
		config, err = LoadBitcoinConf()
		if err != nil {
			return nil, fmt.Errorf("either provide the arguments or bitcoin.conf: %w", err)
		}
		host = config["rpcconnect"]
		user = config["rpcuser"]
		password = config["rpcpassword"]
		if port == "" {
			port = config["rpcport"]
		}
	}

	/** Configure RPC client options
	* bitcoin-cli -rpcconnect=<host>:<port> -rpcport=8332 -rpcuser=<user> -rpcpassword=<password> getblockchaininfo
	**/
	connCfg := &rpcclient.ConnConfig{
		Host:         host + ":" + port,
		User:         user,
		Pass:         password,
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

// LoadBitcoinConf reads the bitcoin.conf file and returns the key-value pairs as a map.
func LoadBitcoinConf() (map[string]string, error) {

	// Default bitcoin.conf path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	path := filepath.Join(homeDir, ".bitcoin", "bitcoin.conf")

	config := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open bitcoin.conf: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading bitcoin.conf: %w", err)
	}

	// Extract connection details
	if config["rpcconnect"] == "" {
		return config, errors.New("rpcconnect (host) not provided")
	}

	if config["rpcuser"] == "" {
		return config, errors.New("rpcuser not provided")
	}

	if config["rpcpassword"] == "" {
		return config, errors.New("rpcpassword not provided")
	}

	if config["rpcport"] == "" {
		config["rpcport"] = "8332"
	}

	return config, nil
}
