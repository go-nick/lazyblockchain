package ui

import (
	"lazyblockchain/constant"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateInput will print an input screen for entering a bitcoin address
func (m *Monitor) CreateInput(method string) {

	_, exist := m.Inputs[method]
	if !exist {
		var variable, label string

		m.Inputs[method] = &Input{
			Method:    method,
			Primitive: tview.NewInputField(),
		}

		switch method {
		case "getblock", "getblockheader":
			variable = "BlockHash"
		case "getblockhash":
			variable = "Height"
		case "getblockstats":
			variable = "Height or Hash"
		case "getchaintxstats", "scantxoutset":
			variable = "Wallet Address"
		}
		label = "Input(" + variable + "):"

		m.Inputs[method].Primitive.
			SetLabel(label).
			SetLabelColor(constant.LightBitcoinYellow).
			SetFieldWidth(128).
			SetFieldBackgroundColor(constant.BitcoinYellow).
			SetBorder(true).
			SetBorderColor(constant.LimeGreen).
			SetTitle(" " + method + " ").
			SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
				// Draw a horizontal line across the middle of the box.
				centerY := y + height/5
				for cx := x + 1; cx < x+width-1; cx++ {
					screen.SetContent(cx, centerY, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
				}
				// Write some text along the horizontal line.
				tview.Print(screen, "Finish Entry (Enter); Cancel (Esc) ", x+1, centerY-1, width-2, tview.AlignCenter, tcell.ColorYellow)
				// Space for other content.
				return x + 1, centerY + 1, width - 2, height - (centerY + 1 - y)
			})
	}

	m.Inputs[method].Focus = true
	m.App.SetRoot(m.Inputs[method].Primitive, true)
	m.SyncFocus(m.Inputs[method].Primitive, true)

	return
}

// CreateForm ...
func (m *Monitor) CreateForm(method string) {

	_, exist := m.Forms[method]
	if !exist {
		m.Forms[method] = &Form{
			Method:    method,
			Focus:     true,
			Primitive: tview.NewForm(),
		}
		switch method {
		case "getblockfilter":
			blockFilter(m.Forms[method].Primitive)
			m.Forms[method].ItensCount = 2
		}
	}

	m.App.SetRoot(m.Forms[method].Primitive, true)
	m.SyncFocus(m.Forms[method].Primitive, true)

	return
}

func blockFilter(form *tview.Form) {
	form.
		AddInputField("block hash", "", 20, nil, nil).
		AddInputField("filter type", "", 20, nil, nil).
		SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			// Draw a horizontal line across the middle of the box.
			centerY := y + height/5
			for cx := x + 1; cx < x+width-1; cx++ {
				screen.SetContent(cx, centerY, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}

			// Write some text along the horizontal line.
			tview.Print(screen, "[green]finish[white](Enter) [green]switch field[white](TAB) [green]cancel[white](Esc) ", x+1, centerY-1, width-2, tview.AlignCenter, tcell.ColorYellow)

			// Space for other content.
			return x + 1, centerY + 1, width - 2, height - (centerY + 1 - y)
		})
}
