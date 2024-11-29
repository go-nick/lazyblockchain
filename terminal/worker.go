package terminal

import (
	"fmt"
	"lazyblockchain/node"

	"github.com/rivo/tview"
)

// worker manages the state of input, form and result channel for go-routines
type worker struct {
	inputState  bool
	inputCH     chan string
	formState   bool
	formCH      chan []string
	resultState bool
	resultCH    chan response
}

type response struct {
	data map[string]interface{}
	err  error
}

// workCmd ...
func (i *Instance) workCmd(method, kind string, rpc node.RPCCommandFunc) response {
	// safely start/stop a new worker go-routine
	if err := i.startWorker(method); err != nil {
		return response{err: err}
	}
	defer i.stopWorker(method)

	// fire go-routine for user input
	// fire go-routine for command response only after user-input
	// start rpc-method only after receiving user-input
	// send result to channel
	if kind == "input" {
		go i.Monitor.CreateInput(method)
		go func() {
			input, ok := i.workerReceiveInput(method)
			if ok {
				result, err := rpc(input)
				i.workerSendResult(method, result, err)
			}
		}()
	}
	if kind == "form" {
		go i.Monitor.CreateForm(method)
		go func() {
			form, ok := i.workerReceiveForm(method)
			if ok {
				result, err := rpc(form...)
				i.workerSendResult(method, result, err)
			}
		}()
	}

	// receiving resultCH after go routine is completed
	return i.workerReceiveResult(method)
}

// startWorker starts channels for go-routines
func (i *Instance) startWorker(method string) error {
	i.mut.Lock()
	defer i.mut.Unlock()

	if work, exists := i.workerPool[method]; exists {
		if work.inputState || work.formState || work.resultState {
			return fmt.Errorf("command %s already running", method)
		}
	} else {
		i.workerPool[method] = &worker{
			inputState:  true,
			inputCH:     make(chan string),
			formState:   true,
			formCH:      make(chan []string),
			resultState: true,
			resultCH:    make(chan response),
		}
	}

	return nil
}

// stopWorker will close the open channels for a given worker
func (i *Instance) stopWorker(method string) {
	i.mut.Lock()
	defer i.mut.Unlock()

	if work, exists := i.workerPool[method]; exists {
		if work.inputState {
			close(work.inputCH)
		}
		if work.formState {
			close(work.formCH)
		}
		if work.resultState {
			close(work.resultCH)
		}
		work.inputState = false
		work.formState = false
		work.resultState = false
	}
}

func (i *Instance) workerSendInput(method string) {
	var inputCH chan string
	var monitorInputText string

	// Access workerPool and Inputs under the mutex lock
	i.mut.Lock()
	if work, exist := i.workerPool[method]; exist {
		if monitorInput, alsoExist := i.Monitor.Inputs[method]; alsoExist {
			if work.inputState {
				inputCH = work.inputCH
				monitorInputText = monitorInput.Primitive.GetText()
			}
		}
	}

	// Release the lock before sending to the channel
	i.mut.Unlock()

	// Send to the input channel
	if inputCH != nil {
		inputCH <- monitorInputText
	}
}

func (i *Instance) workerSendForm(method string) {
	var formCH chan []string
	var itens []string

	// Access workerPool and Inputs under the mutex lock
	i.mut.Lock()
	if work, exist := i.workerPool[method]; exist {
		if form, alsoExist := i.Monitor.Forms[method]; alsoExist {
			if work.formState {
				formCH = work.formCH
				fp := form.Primitive

				for i := 0; i < form.ItensCount; i++ {
					item := fp.GetFormItem(i).(*tview.InputField).GetText()
					itens = append(itens, item)
				}
			}
		}
	}

	// Release the lock before sending to the channel
	i.mut.Unlock()

	// Send to the input channel
	if formCH != nil {
		formCH <- itens
	}
}

func (i *Instance) workerReceiveInput(method string) (string, bool) {
	var inputCH chan string

	// Access workerPool and Inputs under the mutex lock
	i.mut.Lock()
	if work, exist := i.workerPool[method]; exist {
		if work.inputState {
			inputCH = work.inputCH
		}
	}

	// Release the lock before receiving from the channel
	i.mut.Unlock()
	if inputCH != nil {
		input, ok := <-inputCH
		return input, ok
	}

	return "", false
}

func (i *Instance) workerReceiveForm(method string) ([]string, bool) {
	var formCH chan []string

	// Access workerPool and Inputs under the mutex lock
	i.mut.Lock()
	if work, exist := i.workerPool[method]; exist {
		if work.formState {
			formCH = work.formCH
		}
	}

	// Release the lock before receiving from the channel
	i.mut.Unlock()
	if formCH != nil {
		form, ok := <-formCH
		return form, ok
	}

	return []string{}, false
}

func (i *Instance) workerSendResult(method string, result map[string]interface{}, err error) {
	var resultCH chan response
	// Access workerPool under mutex lock
	i.mut.Lock()
	if work, exist := i.workerPool[method]; exist {
		if work.resultState {
			resultCH = work.resultCH
		}
	}

	// Relese the lock before sending in the channel
	i.mut.Unlock()
	resultCH <- response{result, err}
}

func (i *Instance) workerReceiveResult(method string) response {
	var resultCH chan response

	// Access workerPool under mutex lock
	i.mut.Lock()
	if work, exist := i.workerPool[method]; exist {
		if work.resultState {
			resultCH = work.resultCH
		}
	}

	// Release the lock before receiving values in the channel
	i.mut.Unlock()
	r := <-resultCH
	return r
}
