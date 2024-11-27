package ui

import (
	"lazyblockchain/constant"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateForm ...
func (m *Monitor) CreateForm(label string) {

	_, exist := m.Form[label]
	if !exist {
		m.Form[label] = tview.NewForm()
		switch label {
		case constant.FormBlockFilter:
			blockFilter(m.Form[label])
		}
	}

	m.App.SetRoot(m.Form[label], true)
	m.SyncFocus(m.Form[label], true)

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
