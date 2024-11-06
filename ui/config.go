package ui

import (
	"lazyblockchain/node"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// "github.com/gdamore/tcell/v2"

var (
	limeGreen          tcell.Color = tcell.NewRGBColor(15, 251, 207) // Lime green
	lightBlue          tcell.Color = tcell.NewRGBColor(15, 200, 255)
	bitcoinYellow      tcell.Color = tcell.NewRGBColor(247, 148, 19)
	lightBitcoinYellow tcell.Color = tcell.NewRGBColor(255, 184, 19)
)

// Terminal main object for all tview and node objects
type Terminal struct {
	// tview
	app          *tview.Application
	grid         *tview.Grid
	menuView     *tview.List
	logsView     *tview.TextView
	statusView   *tview.TextView
	helpInfoView *tview.TextView
	inputAddress *tview.InputField

	// node package
	rpc           *node.RPC
	commandsCache map[string]node.RPCCommandFunc
}

func newTerminal(rpc *node.RPC) *Terminal {
	t := &Terminal{
		app:          tview.NewApplication(), // terminal new application
		grid:         tview.NewGrid(),        // grid 4x4
		menuView:     tview.NewList(),        // menu view with list of commands
		logsView:     tview.NewTextView(),    // logs view with write/read/copy/paste options
		statusView:   tview.NewTextView(),    // status view with loading progress
		helpInfoView: tview.NewTextView(),    // help view with general shortcuts
		inputAddress: tview.NewInputField(),
		rpc:          rpc, // add the rpc server node connection
	}
	t.RegisterCommands() // Register commands and populate the list for menu view
	return t
}

// Run terminal ui app
func Run(rpc *node.RPC) error {
	t := newTerminal(rpc) // create the terminal application

	// General Shortcuts
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft {
			t.app.SetFocus(t.menuView)
			return nil
		}
		if event.Key() == tcell.KeyRight {
			t.app.SetFocus(t.logsView)
			return nil
		}
		return event
	})

	t.grid.
		SetRows(30, 0).
		SetColumns(50, 0).
		AddItem(t.menuView, 0, 0, 1, 1, 0, 0, true).
		AddItem(t.logsView, 0, 1, 1, 1, 0, 0, false).
		AddItem(t.statusView, 1, 0, 1, 1, 0, 0, false).
		AddItem(t.helpInfoView, 1, 1, 1, 1, 0, 0, false)

	t.resetRootFocus()
	if err := t.app.Run(); err != nil {
		panic(err)
	}

	return nil
}

func (t *Terminal) resetRootFocus() {
	t.grid.
		SetBorder(true).
		SetBorderColor(tcell.ColorBlack).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" LazyBlockchain ").Blur()

	t.menuView.
		SetBorder(true).
		SetBorderPadding(2, 2, 2, 2).
		SetBorderColor(lightBitcoinYellow).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Menu ").Blur()

	t.logsView.
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			t.app.Draw()
		}).
		SetBorder(true).
		SetBorderColor(bitcoinYellow).
		SetBorderPadding(1, 1, 1, 1).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Logs ").Blur()

	t.statusView.SetText(" ").
		SetBorder(true).
		SetBorderColor(lightBlue).
		SetBorderPadding(1, 1, 1, 1).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Status ").Blur()

	t.helpInfoView.
		SetText("Help(F1) | exit(Ctrl C) | Menu(<-) | Logs(->)").
		SetBorder(true).
		SetBorderColor(limeGreen).
		SetBorderPadding(1, 1, 1, 1).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Help ").Blur()

	t.app.SetRoot(t.grid, true).SetFocus(t.menuView).EnableMouse(true)
}
