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
