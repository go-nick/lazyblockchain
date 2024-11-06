package ui

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

// "github.com/gdamore/tcell/v2"

// ExecuteRPC runs an RPC command, handles errors, and displays the response
func (t *Terminal) ExecuteRPC(command string, args ...string) error {
	commandFunc, _ := t.commandsCache[command]          // Get the selected command function from the list
	startTime := time.Now()                             // Capture start time
	timestamp := time.Now().Format("15:04:05-02/01/06") // Capture the timestamp
	doneCH := make(chan bool)                           // Channel to control the loading indicator

	go t.writeLoadingAnimation(doneCH) // start loading go routine
	go func() {                        // Execute the command in a goroutine
		result, err := commandFunc()     // Run the command function
		elapsed := time.Since(startTime) // calculate elapsed time
		if err != nil {
			t.writeError(err)
			return
		}
		// print to logs
		t.writeResponse(result, command, timestamp, elapsed.String())
		defer close(doneCH) // Ensure doneCH is closed when done
	}()
	return nil
}

// writeLoadingAnimation just prints a simple character animation
func (t *Terminal) writeLoadingAnimation(done <-chan bool) {
	loadingChars := []string{"|", "/", "-", "\\",
		"|", "/", "-", "\\", "|", "/", "-", "\\",
		"_", "|", "/", "-", "#", "%", "$", "\\",
		"|", "\\", "-", "/"}

	i := 0
	for {
		select {
		case <-done:
			t.app.QueueUpdateDraw(func() {
				t.statusView.SetText("loading[done]").SetBackgroundColor(tcell.ColorDefault)
				t.helpInfoView.
					SetText("Help(F1) | exit(Ctrl C) | Menu(<- arrow ->)Logs").SetBackgroundColor(tcell.ColorDefault)
			})
			return
		default:
			t.app.QueueUpdateDraw(func() {
				t.statusView.SetText("loading[" + loadingChars[i] + "]").SetBackgroundColor(tcell.ColorDefault)
				t.helpInfoView.
					SetText("Help(F1) | exit(Ctrl C) | Menu(<-) | Logs(->)").SetBackgroundColor(tcell.ColorDefault)
			})
			i = (i + 1) % len(loadingChars)
			time.Sleep(50 * time.Millisecond) // Adjust speed as needed
		}
	}
}

// writeError will pretty print the error to the logs view
func (t *Terminal) writeError(err error) {
	t.app.QueueUpdateDraw(func() {
		t.logsView.SetText(fmt.Sprintf("Error: %v", err), false).SetBackgroundColor(tcell.ColorDefault)
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
		oldTxt := t.logsView.GetText()                                // Retrieve the existing content
		newTxt := oldTxt + NL + title + NL + prettyDataWithLineNums + DIV
		t.logsView.SetText(newTxt, false).SetBackgroundColor(tcell.ColorDefault)
	})
	t.app.Sync()
}

// promptBTCAddress will print an input screen for entering a bitcoin address
func (t *Terminal) promptBTCAddress(addressCh chan<- string) {
	t.inputAddress.
		SetLabel("Input address: ").
		SetFieldWidth(0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				time.Sleep(50 * time.Millisecond) // fixes typing bug
				address := t.inputAddress.GetText()
				addressCh <- address
				close(addressCh)
				t.resetRootFocus()
			}
		}).
		SetFieldBackgroundColor(tcell.ColorRebeccaPurple).
		SetLabelColor(lightBitcoinYellow).
		SetBorderColor(bitcoinYellow).
		SetTitleColor(limeGreen).
		SetTitle(" input mode ")

	t.app.SetRoot(t.inputAddress, true).SetFocus(t.inputAddress).Sync()
	return
}
