package terminal

import (
	"lazyblockchain/constant"
	"lazyblockchain/node"
	"lazyblockchain/ui"
	"sort"
	"time"

	"github.com/rivo/tview"
)

// Instance main object for all tview and node objects
type Instance struct {
	Monitor  *ui.Monitor
	RPC      *node.RPC
	inputCH  chan string
	formCH   chan map[string]string
	resultCH chan struct {
		result map[string]interface{}
		err    error
	}
}

// Setup a new terminal instance
func Setup() *Instance {
	i := &Instance{
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
			Input: tview.NewInputField(),
			Form:  make(map[string]*tview.Form),
		},
	}

	return i
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
			i.inputCH = make(chan string)
			i.resultCH = make(chan struct {
				result map[string]interface{}
				err    error
			})

			go i.Monitor.CreateInput("block hash")
			go func() {
				blockHash, ok := <-i.inputCH
				if ok {
					result, err := i.RPC.GetBlock(blockHash)
					i.resultCH <- struct {
						result map[string]interface{}
						err    error
					}{result, err}
				}
			}()

			res := <-i.resultCH
			return res.result, res.err
		},
		"getblockchaininfo": func(args ...string) (map[string]interface{}, error) {
			return i.RPC.GetBlockchainInfo()
		},
		"getblockcount": func(args ...string) (map[string]interface{}, error) {
			return i.RPC.GetBlockCount()
		},
		"getblockfilter": func(args ...string) (map[string]interface{}, error) {
			i.formCH = make(chan map[string]string)
			i.resultCH = make(chan struct {
				result map[string]interface{}
				err    error
			})

			go i.Monitor.CreateForm(constant.FormBlockFilter)
			go func() {
				form, ok := <-i.formCH
				blockHash, _ := form["blockhash"]
				filterType, _ := form["filtertype"]
				if ok {
					result, err := i.RPC.GetBlockFilter(blockHash, filterType)
					i.resultCH <- struct {
						result map[string]interface{}
						err    error
					}{result, err}
				}
			}()

			res := <-i.resultCH
			return res.result, res.err
		},
		"getblockhash": func(args ...string) (map[string]interface{}, error) {
			i.inputCH = make(chan string)
			i.resultCH = make(chan struct {
				result map[string]interface{}
				err    error
			})

			go i.Monitor.CreateInput("height index")
			go func() {
				heightIndex, ok := <-i.inputCH
				if ok {
					result, err := i.RPC.GetBlockHash(heightIndex)
					i.resultCH <- struct {
						result map[string]interface{}
						err    error
					}{result, err}
				}
			}()

			res := <-i.resultCH
			return res.result, res.err
		},
		"getblockheader": func(args ...string) (map[string]interface{}, error) {
			i.inputCH = make(chan string)
			i.resultCH = make(chan struct {
				result map[string]interface{}
				err    error
			})

			go i.Monitor.CreateInput("block hash")
			go func() {
				blockHash, ok := <-i.inputCH
				if ok {
					result, err := i.RPC.GetBlockHeader(blockHash)
					i.resultCH <- struct {
						result map[string]interface{}
						err    error
					}{result, err}
				}
			}()

			res := <-i.resultCH
			return res.result, res.err
		},
		"getblockstats": func(args ...string) (map[string]interface{}, error) {
			i.inputCH = make(chan string)
			i.resultCH = make(chan struct {
				result map[string]interface{}
				err    error
			})

			go i.Monitor.CreateInput("block hash OR height")
			go func() {
				hashOrHeight, ok := <-i.inputCH
				if ok {
					result, err := i.RPC.GetBlockHeader(hashOrHeight)
					i.resultCH <- struct {
						result map[string]interface{}
						err    error
					}{result, err}
				}
			}()

			res := <-i.resultCH
			return res.result, res.err
		},
		"getchaintips": func(args ...string) (map[string]interface{}, error) {
			return i.RPC.GetChainTips()
		},
		"getchaintxstats": func(args ...string) (map[string]interface{}, error) {
			i.inputCH = make(chan string)
			i.resultCH = make(chan struct {
				result map[string]interface{}
				err    error
			})

			go i.Monitor.CreateInput("nblocks")
			go func() {
				hashOrHeight, ok := <-i.inputCH
				if ok {
					result, err := i.RPC.GetChainTxStats(hashOrHeight)
					i.resultCH <- struct {
						result map[string]interface{}
						err    error
					}{result, err}
				}
			}()

			res := <-i.resultCH
			return res.result, res.err
		},
		"scantxoutset": func(args ...string) (map[string]interface{}, error) {
			i.inputCH = make(chan string)
			i.resultCH = make(chan struct {
				result map[string]interface{}
				err    error
			})

			go i.Monitor.CreateInput("wallet address")
			go func() {
				address, ok := <-i.inputCH
				if ok {
					result, err := i.RPC.ScanTxOutSet(address)
					i.resultCH <- struct {
						result map[string]interface{}
						err    error
					}{result, err}
				}
			}()

			res := <-i.resultCH
			return res.result, res.err
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
