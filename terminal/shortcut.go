package terminal

import (
	"lazyblockchain/constant"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Shortcuts sets General Shortcuts
func (i *Instance) Shortcuts() {
	i.Monitor.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {

		// input shortcuts
		case tcell.KeyEnter:
			if i.Monitor.Input.HasFocus() {
				i.inputCH <- i.Monitor.Input.GetText()
				close(i.inputCH)
				i.Monitor.DefaultLayout()
				return nil
			}

			for label, form := range i.Monitor.Form {
				if form.HasFocus() {
					switch label {
					case constant.FormBlockFilter:
						f1 := form.GetFormItem(0).(*tview.InputField).GetText()
						f2 := form.GetFormItem(1).(*tview.InputField).GetText()
						i.formCH <- map[string]string{
							"blockhash":  f1,
							"filtertype": f2,
						}
						close(i.formCH)
						i.Monitor.DefaultLayout()
						return nil
					}
				}
			}

		case tcell.KeyEsc:
			if i.Monitor.Input.HasFocus() {
				close(i.inputCH)
				close(i.resultCH)
				i.Monitor.DefaultLayout()
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
