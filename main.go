package main

import (
	"lazyblockchain/node"
	"lazyblockchain/terminal"
	"lazyblockchain/ui"
	"log"

	"github.com/spf13/cobra"
)

// CLI arguments
var (
	host     string
	port     string
	user     string
	password string
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "lazyblockchain",
		Short: "LazyBlockchain is a Bitcoin node monitoring app",
		Run:   runApp,
	}

	// Add flags to the root command
	rootCmd.Flags().StringVar(&host, "host", "", "RPC server host")
	rootCmd.Flags().StringVar(&port, "port", "", "RPC server port")
	rootCmd.Flags().StringVar(&user, "user", "", "RPC username")
	rootCmd.Flags().StringVar(&password, "password", "", "RPC password")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// runApp initializes and runs the application
func runApp(cmd *cobra.Command, args []string) {
	// Check if any arguments are passed, otherwise fallback to bitcoin.conf
	rpc, err := node.ConnectRPC(host, port, user, password)
	if err != nil {
		log.Fatal(err)
	}

	// Set up the terminal infra
	terminal := terminal.Setup()
	terminal.RegisterCommands(rpc)
	terminal.Shortcuts()

	// Set up the UI
	ui.Setup(terminal.Monitor)
	// Run
	terminal.Run()
}
