package terminal

import (
	"lazyblockchain/node"
	"lazyblockchain/ui"
	"sort"
	"sync"
	"time"

	"github.com/rivo/tview"
)

// Instance for running TUI
type Instance struct {
	Monitor    *ui.Monitor
	RPC        *node.RPC
	workerPool map[string]*worker
	mut        sync.Mutex
}

// Setup a new terminal instance
func Setup() *Instance {
	return &Instance{
		Monitor: &ui.Monitor{
			App:  tview.NewApplication(), // terminal new application
			Grid: tview.NewGrid(),        // grid 4x4
			View: &ui.View{
				Menu:   tview.NewList(),     // menu view with list of commands
				Logs:   tview.NewTextView(), // logs view with write/read/copy/paste options
				Status: tview.NewTextView(), // status view with loading progress
				Info:   tview.NewTextView(), // help view with general shortcuts
			},
			Help: &ui.Help{
				View1: tview.NewTextView(),
				View2: tview.NewTextView(),
				View3: tview.NewTextView(),
			},
			Inputs: make(map[string]*ui.Input),
			Forms:  make(map[string]*ui.Form),
		},
		workerPool: make(map[string]*worker),
	}
}

// RegisterCommands initializes and caches a map of RPC commands
func (i *Instance) RegisterCommands(rpc *node.RPC) {
	i.RPC = rpc

	i.RPC.CmdCache = map[string]node.RPCCommandFunc{
		"mockresponse": func(args ...string) (map[string]interface{}, error) {
			return i.RPC.MockResponse()
		},
		"getbestblockhash": func(args ...string) (map[string]interface{}, error) {
			return i.RPC.GetBestBlockHash()
		},
		"getblock": func(args ...string) (map[string]interface{}, error) {
			res := i.workCmd("getblock", "input", i.RPC.GetBlock)
			return res.data, res.err
		},
		"getblockchaininfo": func(args ...string) (map[string]interface{}, error) {
			return i.RPC.GetBlockchainInfo()
		},
		"getblockcount": func(args ...string) (map[string]interface{}, error) {
			return i.RPC.GetBlockCount()
		},
		"getblockfilter": func(args ...string) (map[string]interface{}, error) {
			res := i.workCmd("getblockfilter", "form", i.RPC.GetBlockFilter)
			return res.data, res.err
		},
		"getblockhash": func(args ...string) (map[string]interface{}, error) {
			res := i.workCmd("getblockhash", "input", i.RPC.GetBlockHash)
			return res.data, res.err
		},
		"getblockheader": func(args ...string) (map[string]interface{}, error) {
			res := i.workCmd("getblockheader", "input", i.RPC.GetBlockHeader)
			return res.data, res.err
		},
		"getblockstats": func(args ...string) (map[string]interface{}, error) {
			res := i.workCmd("getblockstats", "input", i.RPC.GetBlockStats)
			return res.data, res.err
		},
		"getchaintips": func(args ...string) (map[string]interface{}, error) {
			return i.RPC.GetChainTips()
		},
		"getchaintxstats": func(args ...string) (map[string]interface{}, error) {
			res := i.workCmd("getchaintxstats", "input", i.RPC.GetChainTxStats)
			return res.data, res.err
		},
		"scantxoutset": func(args ...string) (map[string]interface{}, error) {
			res := i.workCmd("scantxoutset", "input", i.RPC.ScanTxOutSet)
			return res.data, res.err
		},
		// Add more commands here as needed
	}

	// Extract command names into a slice and sort alphabetically
	commandNames := make([]string, 0, len(i.RPC.CmdCache))
	for cmdName := range i.RPC.CmdCache {
		commandNames = append(commandNames, cmdName)
	}
	sort.Strings(commandNames)

	// Populate the tview list with sorted commands
	for _, command := range commandNames {
		cmdName := command // capture range variable to avoid issues in closures
		i.Monitor.View.Menu.AddItem(cmdName, "", 0, func() {
			i.ExecuteRPC(cmdName)
		})
	}
}

// ExecuteRPC runs an RPC command, handles errors, and displays the response
func (i *Instance) ExecuteRPC(command string, args ...string) {
	startTime := time.Now()                             // Capture start time
	timestamp := time.Now().Format("15:04:05-02/01/06") // Capture the timestamp

	// Channel to control the loading
	loadCH := make(chan bool)
	// Get the selected command function from the list
	commandFunc, _ := i.RPC.CmdCache[command]

	// start loading animation
	go i.Monitor.Loading(loadCH)

	// start rpc call
	go func() {
		// Ensure loadCH is closed when done
		defer close(loadCH)

		result, err := commandFunc()
		if err != nil {
			i.Monitor.LogError(err)
			return
		}

		// calculate elapsed time and print response if any
		elapsed := time.Since(startTime)
		if result != nil {
			i.Monitor.LogResponse(result, command, timestamp, elapsed.String())
		}

		defer i.Monitor.SyncFocus(i.Monitor.View.Logs, true)
	}()
}

// Run terminal cli-app
func (i *Instance) Run() {
	if err := i.Monitor.App.Run(); err != nil {
		panic(err)
	}
}
