package asciichart

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var symbolsPalette = []string{
	"#",
	"@",
	"$",
	"O",
}

var colorsPalette = []func(format string, a ...interface{}) string{
	color.RedString,
	color.YellowString,
	color.BlueString,
	color.GreenString,
	color.MagentaString,
	color.CyanString,
}

func getSymbol(index int, useColors bool) string {
	if useColors {
		return "#"
	}
	return symbolsPalette[index%len(symbolsPalette)]
}

func addColor(format string, index int, useColors bool, a ...interface{}) string {
	if useColors {
		return colorsPalette[index%len(colorsPalette)](format, a...)
	}

	return fmt.Sprintf(format, a...)
}

func strPad(str, pad string, width int, left bool) string {
	padStr := strings.Repeat(pad, (width-len(str))/len(pad))

	if left {
		return padStr + str
	}
	return str + padStr
}
