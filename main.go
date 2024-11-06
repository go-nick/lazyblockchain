package main

import (
	"lazyblockchain/node"
	"lazyblockchain/ui"
	"log"
)

func main() {
	rpc, err := node.ConnectRPC()
	if err != nil {
		log.Fatal(err)
	}

	if err := ui.Run(rpc); err != nil {
		log.Fatal(err)
	}
}
