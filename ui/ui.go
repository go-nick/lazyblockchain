package ui

import (
	"fmt"
	"lazyblockchain/constant"
	"lazyblockchain/logs"
	"time"

	// "github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ClearLogs from the record and screen
func (m *Monitor) ClearLogs() {
	m.Record.Clear()
	m.View.Logs.SetText(m.Record.Text)
}

// LogResponse function to write at logs view, with padded line numbers for consistent indentation
func (m *Monitor) LogResponse(data interface{}, command, timestamp, elapsed string) {
	// Append new log data to existing content
	title := command + ": " + timestamp + " duration: " + elapsed

	// init logs if it is the first time
	if m.Record == nil {
		m.Record = &logs.Record{
			Selected:  "",
			Text:      "",
			TextLines: []string{},
			Line:      0,
		}
	}
	// Parse log into standard format and record
	m.Record.Process(data, title)
	// Write
	m.App.QueueUpdateDraw(func() {
		m.logNavigation(m.Record.Text)
	})
}

func (m *Monitor) logNavigation(txt string) {

	m.View.Logs.SetText(txt).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

			if !m.View.Logs.HasFocus() || len(m.Record.TextLines) <= 1 {
				return event
			}

			pressed := event.Key()

			m.Record.Highlight(pressed)
			m.View.Logs.SetText(m.Record.Text)
			m.View.Logs.Highlight(constant.WordWrap)

			// Deactivate default shortcut behavior for these specific cases
			// by returning nil and scrolling to target line
			if pressed == tcell.KeyPgUp ||
				pressed == tcell.KeyPgDn ||
				pressed == tcell.KeyHome ||
				pressed == tcell.KeyEnd {

				m.View.Logs.ScrollTo(m.Record.Line-2, 0)
				return nil
			}

			return event
		})
}

// LogError will pretty print the error to the logs view
func (m *Monitor) LogError(err error) {
	m.App.QueueUpdateDraw(func() {
		m.View.Logs.SetText(fmt.Sprintf("Error: %v", err))
	})
	return
}

// Loading just prints a simple character animation
func (m *Monitor) Loading(loadCH <-chan bool) {
	loadingChars := []string{"|", "/", "-", "\\",
		"|", "/", "-", "\\", "|", "/", "-", "\\",
		"_", "|", "/", "-", "#", "%", "$", "\\",
		"|", "\\", "-", "/"}

	idx := 0
	for {
		select {
		case <-loadCH:
			m.App.QueueUpdateDraw(func() {
				m.View.Status.SetText("loading[done]").SetBackgroundColor(tcell.ColorDefault)
				m.View.Info.SetBackgroundColor(tcell.ColorDefault)
			})
			return
		default:
			m.App.QueueUpdateDraw(func() {
				m.View.Status.SetText("loading[" + loadingChars[idx] + "]").SetBackgroundColor(tcell.ColorDefault)
				m.View.Info.SetBackgroundColor(tcell.ColorDefault)
			})
			idx = (idx + 1) % len(loadingChars)
			time.Sleep(50 * time.Millisecond) // Adjust speed as needed
		}
	}
}

// CreateInput will print an input screen for entering a bitcoin address
func (m *Monitor) CreateInput(label string) {
	m.Input.
		SetLabel("Input[" + label + "]: ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorRebeccaPurple).
		SetLabelColor(constant.LightBitcoinYellow).
		SetBorderColor(constant.BitcoinYellow).
		SetTitleColor(constant.LimeGreen).
		SetTitle(" input mode ").
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

	m.App.SetRoot(m.Input, true)
	m.SyncFocus(m.Input, true)

	return
}

// SyncFocus will change the focus, update the info texts and sync
func (m *Monitor) SyncFocus(p tview.Primitive, sync bool) {
	switch p {
	case m.View.Menu:
		m.View.Info.SetText(constant.InfoViewMenu).SetBackgroundColor(tcell.ColorDefault)
	case m.View.Logs:
		m.View.Info.SetText(constant.InfoViewLogs).SetBackgroundColor(tcell.ColorDefault)
	}

	m.App.SetFocus(p)
	if sync {
		m.App.Sync()
	}
}
