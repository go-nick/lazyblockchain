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

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 829e0d4 (Still thinking what is best software architecture here)
// CreateInput will print an input screen for entering a bitcoin address
func (m *Monitor) CreateInput(label string) {
	m.Input.
		SetLabel("Input[" + label + "]: ").
<<<<<<< HEAD
=======
// writeError will pretty print the error to the logs view
func (t *Terminal) writeError(err error) {
	t.app.QueueUpdateDraw(func() {
		t.logsView.SetText(fmt.Sprintf("Error: %v", err)).SetBackgroundColor(tcell.ColorDefault)
	})
	return
}

// writeResponse function to write at logs view, with padded line numbers for consistent indentation
func (t *Terminal) writeResponse(data interface{}, command, timestamp, elapsed string) {
	var prettyData string
	// Try to marshal data as JSON with indentation; fallback to plain text
	if jsonBytes, err := json.MarshalIndent(data, "", "  "); err == nil {
		prettyData = string(jsonBytes)
	} else {
		prettyData = fmt.Sprintf("%v", data)
	}

	// add line numbers
	lines := strings.Split(prettyData, "\n")           // line numbers with 0-padding to each line
	lineNumWidth := len(fmt.Sprintf("%d", len(lines))) // width based on total line count
	for i, line := range lines {                       // line numbers with leading 0s for equal width
		lines[i] = fmt.Sprintf("%0*d: %s", lineNumWidth, i+1, line)
	}
	prettyDataWithLineNums := strings.Join(lines, "\n")

	// Write
	t.app.QueueUpdateDraw(func() {
		title := command + ": " + timestamp + " duration: " + elapsed // Append new log data to existing content
		oldTxt := t.logsView.GetText(true)                            // Retrieve the existing content
		newTxt := oldTxt + NL + title + NL + prettyDataWithLineNums + DIV
		numSelections := 0
		go func() {
			for _, word := range strings.Split(newTxt, " ") {
				if word == "the" {
					word = "[#ff0000]the[white]"
				}
				if word == "to" {
					word = fmt.Sprintf(`["%d"]to[""]`, numSelections)
					numSelections++
				}
				fmt.Fprintf(t.logsView, "%s ", word)
				time.Sleep(200 * time.Millisecond)
			}
		}()
		t.logsView.SetText(newTxt).SetDoneFunc(func(key tcell.Key) {
			currentSelection := t.logsView.GetHighlights()
			if key == tcell.KeyEnter {
				if len(currentSelection) > 0 {
					t.logsView.Highlight()
				} else {
					t.logsView.Highlight("0").ScrollToHighlight()
				}
			} else if len(currentSelection) > 0 {
				index, _ := strconv.Atoi(currentSelection[0])
				if key == tcell.KeyTab {
					index = (index + 1) % numSelections
				} else if key == tcell.KeyBacktab {
					index = (index - 1 + numSelections) % numSelections
				} else {
					return
				}
				t.logsView.Highlight(strconv.Itoa(index)).ScrollToHighlight()
			}
		}).SetBackgroundColor(tcell.ColorDefault)
	})
	t.app.Sync()
}

// promptBTCAddress will print an input screen for entering a bitcoin address
func (t *Terminal) promptBTCAddress(addressCh chan<- string) {
	t.inputAddress.
		SetLabel("Input address: ").
>>>>>>> 760e218 (changing textArea to textView)
=======
>>>>>>> 829e0d4 (Still thinking what is best software architecture here)
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
