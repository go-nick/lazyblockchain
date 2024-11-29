package terminal

import (
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
)

// Shortcuts sets General Shortcuts
func (i *Instance) Shortcuts() {
	i.Monitor.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {

		// input shortcuts
		case tcell.KeyEnter, tcell.KeyEsc:
			if ev := i.enterEscInputForm(event); ev == nil {
				return nil
			}

		// logs view shortcuts
		case tcell.KeyDelete:
			if i.Monitor.View.Logs.HasFocus() {
				i.Monitor.ClearLogs()
				i.Monitor.DefaultLayout()
				return nil
			}

		case tcell.KeyCtrlY:
			if i.Monitor.View.Logs.HasFocus() {
				if i.Monitor.Record != nil {
					if err := clipboard.WriteAll(i.Monitor.Record.Selected); err != nil {
						panic(err)
					}
					return nil
				}
			}

		// general shortcuts
		case tcell.KeyLeft:
			if i.Monitor.View.Logs.HasFocus() {
				i.Monitor.SyncFocus(i.Monitor.View.Menu, false)
				return nil
			}
		case tcell.KeyRight:
			if i.Monitor.View.Menu.HasFocus() {
				i.Monitor.SyncFocus(i.Monitor.View.Logs, false)
				return nil
			}
		} // end

		return event
	})
}

func (i *Instance) enterEscInputForm(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()

	for method, input := range i.Monitor.Inputs {
		if input.Focus {
			if key == tcell.KeyEnter {
				i.workerSendInput(method)
			}
			if key == tcell.KeyEsc {
				i.stopWorker(method)
			}
			i.Monitor.DefaultLayout()
			return nil
		}
	}

	for method, form := range i.Monitor.Forms {
		if form.Focus {
			if key == tcell.KeyEnter {
				i.workerSendForm(method)
			}
			if key == tcell.KeyEsc {
				i.stopWorker(method)
			}
			i.Monitor.DefaultLayout()
			return nil
		}
	}

	return event
}
