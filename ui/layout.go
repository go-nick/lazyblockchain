package ui

import (
	"lazyblockchain/constant"
	"lazyblockchain/logs"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Monitor ...
type Monitor struct {
	App    *tview.Application
	Grid   *tview.Grid
	View   *View
	Help   *Help
	Record *logs.Record
	Inputs map[string]*Input
	Forms  map[string]*Form
}

// View holds the main 4 views of the cli app
type View struct {
	Menu   *tview.List
	Logs   *tview.TextView
	Status *tview.TextView
	Info   *tview.TextView
}

// Input wraps tview.InputField
type Input struct {
	Method    string
	Focus     bool
	Primitive *tview.InputField
}

// Form wraps tview.Form
type Form struct {
	Method     string
	Focus      bool
	ItensCount int
	Primitive  *tview.Form
}

// Help holds the help support pages with shortcuts
type Help struct {
	Frame *tview.Frame
	View1 *tview.TextView
	View2 *tview.TextView
	View3 *tview.TextView
}

// Setup terminal ui app
func Setup(m *Monitor) {

	// TODO for removal later
	// m.Help.View1.SetDynamicColors(true).SetText(constant.HelpSupport1)
	// m.Help.View2.SetDynamicColors(true).SetText(constant.HelpSupport2)
	// m.Help.View3.SetDynamicColors(true).SetText(constant.HelpSupport3)
	// m.Help.Frame = tview.NewFrame(m.Help.View1).
	// 	SetBorders(1, 1, 0, 0, 2, 2).
	// 	SetBorder(true)

	m.Grid.
		SetRows(0, 5).
		SetColumns(33, 0).
		AddItem(m.View.Menu, 0, 0, 1, 1, 0, 0, true).
		AddItem(m.View.Logs, 0, 1, 1, 1, 0, 0, false).
		AddItem(m.View.Status, 1, 0, 1, 1, 0, 0, false).
		AddItem(m.View.Info, 1, 1, 1, 1, 0, 0, false)

	m.DefaultLayout()
}

// DefaultLayout of the tview cli app
func (m *Monitor) DefaultLayout() {
	m.clearInputs()
	m.clearForms()

	m.Grid.
		SetBorder(true).
		SetBorderColor(tcell.ColorBlack).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(constant.TitleGrid)

	m.View.Menu.
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(constant.LightBitcoinYellow).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(constant.TitleMenu)

	m.View.Logs.
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetBorder(true).
		SetBorderColor(constant.BitcoinYellow).
		SetBorderPadding(0, 0, 0, 0).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(constant.TitleLogs)

	m.View.Status.SetText(" ").
		SetBorder(true).
		SetBorderColor(constant.LightBlue).
		SetBorderPadding(0, 0, 0, 0).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(constant.TitleStatus)

	m.View.Info.
		SetText(constant.InfoViewMenu).
		SetDynamicColors(true).
		SetRegions(true).
		SetBorder(true).
		SetBorderColor(constant.LimeGreen).
		SetBorderPadding(0, 0, 0, 0).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(constant.TitleInfo)

	m.App.
		SetRoot(m.Grid, true).
		EnableMouse(true)

	m.SyncFocus(m.View.Menu, true)
}

// clearInputs sets focus of all inputs to false
func (m *Monitor) clearInputs() {
	if len(m.Inputs) > 0 {
		for _, input := range m.Inputs {
			input.Focus = false
		}
	}
}

// clearForms sets focus of all forms to false
func (m *Monitor) clearForms() {
	if len(m.Inputs) > 0 {
		for _, input := range m.Inputs {
			input.Focus = false
		}
	}
}
