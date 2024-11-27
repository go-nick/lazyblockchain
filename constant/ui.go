package constant

import (
	"regexp"

	"github.com/gdamore/tcell/v2"
)

const (

	// WrapStyle defines word wrap color style
	WrapStyle string = "[black:blue]"
	// WordWrap implements tview color style tag standards
	WordWrap string = WrapStyle + "%s[-:-]"
	// DIV is a simple divider; TODO: find something better
	DIV                    = " "
	TitleGrid       string = " LazyBlockchain "
	TitleMenu       string = " menu "
	TitleLogs       string = " logs "
	TitleStatus     string = " status "
	TitleInfo       string = " info "
	FormBlockFilter string = "getblockfilter"
)

var (
	LimeGreen          tcell.Color = tcell.NewRGBColor(15, 251, 207)
	LightBlue          tcell.Color = tcell.NewRGBColor(15, 200, 255)
	BitcoinYellow      tcell.Color = tcell.NewRGBColor(247, 148, 19)
	LightBitcoinYellow tcell.Color = tcell.NewRGBColor(255, 184, 19)

	HexLimeGreen         string = "#0ffbcf"
	HexLightBlue         string = "#0fc8ff"
	HexBitcoinYellow     string = "#f79413"
	HexLightBitcoinYello string = "#ffb813"
)

var (
	// Regular expression to match strings like \"<any_string>\":
	JsonFieldRegexp = regexp.MustCompile(`"([^"]+)":`)
)
