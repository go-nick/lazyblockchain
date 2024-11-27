package logs

import (
	"fmt"
	"lazyblockchain/constant"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// Record holds all logs and it's selection
type Record struct {
	Selected  string
	Text      string
	TextLines []string
	Line      int
}

// Clear all recorded logs
func (r *Record) Clear() {
	r.TextLines = []string{""}
	r.join()
}

// Process will receive a rpc function response and prepare it for logging
func (r *Record) Process(data interface{}, title string) {
	// parse marshalled data
	txtLines := parse(marshal(data), title)
	if len(r.TextLines) == 0 {
		r.TextLines = txtLines
	} else {
		r.TextLines = append(r.TextLines, txtLines...)
	}

	r.TextLines = append(r.TextLines, constant.DIV)
	r.join()
}

// Highlight will iterate words inside the record lines and highlight the selected one
func (r *Record) Highlight(pressed tcell.Key) {
	// clear previous word wrap
	for idx, line := range r.TextLines {
		if strings.Contains(line, constant.WrapStyle) {
			line = strings.Replace(line, constant.WrapStyle, "", 1)
			line = strings.Replace(line, "[-:-]", "", 1)
			r.TextLines[idx] = line
		}
	}

	selected := r.wrapValue(pressed)
	if selected == "" {
		return // break execution
	}
	r.selectValue(selected)

	word := fmt.Sprintf(constant.WordWrap, selected)
	line := strings.Replace(r.TextLines[r.Line], selected, word, 1)
	r.TextLines[r.Line] = line
	r.join()
}
