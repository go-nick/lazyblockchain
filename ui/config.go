package ui

import (
	"lazyblockchain/node"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// "github.com/gdamore/tcell/v2"

// Terminal ...
type Terminal struct {
	// tview
	app          *tview.Application
	grid         *tview.Grid
	menuView     *tview.List
	logsView     *tview.TextArea
	statusView   *tview.TextView
	helpInfoView *tview.TextView
	inputAddress *tview.InputField

	// node package
	rpc           *node.RPC
	commandsCache map[string]node.RPCCommandFunc
}

var (
	limeGreen          tcell.Color = tcell.NewRGBColor(15, 251, 207) // Lime green
	lightBlue          tcell.Color = tcell.NewRGBColor(15, 200, 255)
	bitcoinYellow      tcell.Color = tcell.NewRGBColor(247, 148, 19)
	lightBitcoinYellow tcell.Color = tcell.NewRGBColor(255, 184, 19)
)

func newTerminal(rpc *node.RPC) *Terminal {
	return &Terminal{
		app:          tview.NewApplication(), // terminal new application
		grid:         tview.NewGrid(),        // grid 4x4
		menuView:     tview.NewList(),        // menu view with list of commands
		logsView:     tview.NewTextArea(),    // logs view with write/read/copy/paste options
		statusView:   tview.NewTextView(),    // status view with loading progress
		helpInfoView: tview.NewTextView(),    // help view with general shortcuts
		inputAddress: tview.NewInputField(),
		rpc:          rpc, // add the rpc server node connection
	}
}

// Run terminal ui app
func Run(rpc *node.RPC) error {
	t := newTerminal(rpc) // create the terminal application
	t.RegisterCommands()  // Register commands and populate the list for menu view

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
		SetRows(35, 0).
		SetColumns(50, 0).
		AddItem(t.menuView, 0, 0, 1, 1, 0, 0, true).
		AddItem(t.logsView, 0, 1, 1, 1, 0, 0, false).
		AddItem(t.statusView, 1, 0, 1, 1, 0, 0, false).
		AddItem(t.helpInfoView, 1, 1, 1, 1, 0, 0, false)

	t.resetRootFocus()
	if err := t.app.Run(); err != nil {
		panic(err)
	}

	t.app.QueueUpdateDraw(func() {
		t.logsView.SetText("Hello BTC!", false).SetBackgroundColor(tcell.ColorDefault)
		t.statusView.SetText(" ").SetBackgroundColor(tcell.ColorDefault)
		t.helpInfoView.
			SetText("Help(F1) | exit(Ctrl C) | Menu(<-) | Logs(->)").SetBackgroundColor(tcell.ColorDefault)
	})
	t.app.ForceDraw()

	return nil
}

func (t *Terminal) resetRootFocus() {
	t.grid.SetBackgroundColor(tcell.ColorDefault)
	t.menuView.SetBorder(true).SetBorderColor(lightBitcoinYellow).SetTitle(" Menu ").SetBackgroundColor(tcell.ColorDefault)
	t.logsView.SetBorder(true).SetBorderColor(bitcoinYellow).SetTitle(" Logs ").SetBackgroundColor(tcell.ColorDefault)
	t.statusView.SetText(" ").SetBorder(true).SetBorderColor(lightBlue).SetTitle(" Status ").SetBackgroundColor(tcell.ColorDefault)
	t.helpInfoView.
		SetText("Help(F1) | exit(Ctrl C) | Menu(<-) | Logs(->)").
		SetBorder(true).SetBorderColor(limeGreen).SetTitle(" Help ").SetBackgroundColor(tcell.ColorDefault)

	t.app.SetRoot(t.grid, true).SetFocus(t.menuView).EnableMouse(true)
}
