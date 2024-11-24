package main

import (
	"lazyblockchain/node"
	"lazyblockchain/terminal"
	"lazyblockchain/ui"
	"log"
)

func main() {
	terminal := terminal.Setup()

	rpc, err := node.ConnectRPC()
	if err != nil {
		log.Fatal(err)
	}

	terminal.RegisterCommands(rpc)
	terminal.Shortcuts()

	ui.Setup(terminal.Monitor)
	terminal.Run()
}
