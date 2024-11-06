package ui

import (
	"lazyblockchain/node"
	"sort"
)

// "github.com/gdamore/tcell/v2"

// RegisterCommands initializes and caches a map of RPC commands
func (t *Terminal) RegisterCommands() {
	// Initialize the command map
	t.commandsCache = map[string]node.RPCCommandFunc{
		"mockresponse": func(args ...string) (map[string]interface{}, error) {
			return t.rpc.MockResponse()
		},
		"getblockchaininfo": func(args ...string) (map[string]interface{}, error) {
			return t.rpc.GetBlockchainInfo()
		},
		"CalculateBTCFromAddress": func(args ...string) (map[string]interface{}, error) {
			inputCH := make(chan string)
			resultCH := make(chan struct {
				result map[string]interface{}
				err    error
			})

			go t.promptBTCAddress(inputCH)
			go func() {
				address, ok := <-inputCH
				if ok {
					result, err := t.rpc.CalculateBTCFromAddress(address)
					resultCH <- struct {
						result map[string]interface{}
						err    error
					}{result, err}
				}
			}()

			res := <-resultCH
			return res.result, res.err
		},
		// Add more commands here as needed
	}

	// Extract command names into a slice and sort alphabetically
	commandNames := make([]string, 0, len(t.commandsCache))
	for cmdName := range t.commandsCache {
		commandNames = append(commandNames, cmdName)
	}
	sort.Strings(commandNames)

	// Populate the tview list with sorted commands
	for _, command := range commandNames {
		cmdName := command // capture range variable to avoid issues in closures
		t.menuView.AddItem(cmdName, "", rune(cmdName[0]), func() {
			t.ExecuteRPC(cmdName)
		})
	}
}
