package logs

import (
	"encoding/json"
	"fmt"
	"lazyblockchain/constant"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func (r *Record) join() {
	r.Text = strings.Join(r.TextLines, "\n")
}

// marshal will try to marshal data as JSON with indentation; fallback to plain text
func marshal(data interface{}) string {
	prettyData := ""
	if jsonBytes, err := json.MarshalIndent(data, "", "  "); err == nil {
		prettyData = string(jsonBytes)
	} else {
		prettyData = fmt.Sprintf("%v", data)
	}
	return prettyData
}

func parse(newData, title string) []string {
	txtLines := strings.Split(newData, "\n")
	lineNumWidth := len(fmt.Sprintf("%d", len(txtLines))) // width based on total line count

	for i, line := range txtLines { // line numbers with leading 0s for equal width
		coloredLine := constant.JsonFieldRegexp.ReplaceAllString(line, fmt.Sprintf(`"[%s]${1}[white]":`, constant.HexLightBitcoinYello))
		txtLines[i] = fmt.Sprintf("[%s]%0*d: [white]%s", constant.HexBitcoinYellow, lineNumWidth, i+1, coloredLine)
	}

	newTxtLines := []string{title}
	newTxtLines = append(newTxtLines, txtLines...)
	return newTxtLines
}

func (r *Record) wrapValue(pressed tcell.Key) string {

	totalLines := len(r.TextLines) - 1

	switch pressed {
	case tcell.KeyUp:
		if r.Line > 0 {
			r.Line--
		} else if r.Line <= 0 {
			r.Line = 0
			return ""
		}
	case tcell.KeyDown:
		if r.Line < totalLines {
			r.Line++
		} else if r.Line >= totalLines {
			r.Line = totalLines
			return ""
		}

	case tcell.KeyPgUp:
		if r.Line > 0 {
			r.Line = r.Line - 10
		}
		if r.Line <= 0 {
			r.Line = 2
		}
	case tcell.KeyPgDn:
		if r.Line < totalLines {
			r.Line = r.Line + 10
		}
		if r.Line >= totalLines {
			r.Line = totalLines - 2
		}
	case tcell.KeyHome:
		r.Line = 2
	case tcell.KeyEnd:
		r.Line = totalLines - 2
	}

	if !constant.JsonFieldRegexp.MatchString(r.TextLines[r.Line]) {
		return r.wrapValue(pressed)
	}

	splittedLine := strings.Split(r.TextLines[r.Line], " ")
	word := splittedLine[len(splittedLine)-1]
	return word
}

func (r *Record) selectValue(selected string) {
	if strings.HasPrefix(selected, `"`) {
		selected, _ = strings.CutPrefix(selected, `"`)
	}
	if strings.HasSuffix(selected, `,`) {
		selected, _ = strings.CutSuffix(selected, `,`)
	}
	if strings.HasSuffix(selected, `"`) {
		selected, _ = strings.CutSuffix(selected, `"`)
	}
	r.Selected = selected
}
