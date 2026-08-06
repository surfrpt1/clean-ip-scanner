package utils

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const headerWidth = 49

func CenterText(width int, text string) string {
	pad := width - len(text)
	if pad <= 0 {
		return text
	}
	left := pad / 2
	return strings.Repeat(" ", left) + text
}

func CenterInBox(width int, text string) string {
	pad := width - len(text)
	if pad <= 0 {
		return text
	}
	left := pad / 2
	right := pad - left
	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}

func PrintHeader() {
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println(strings.Repeat("=", headerWidth))
	cyan.Println(CenterText(headerWidth, "CLEAN IP SCANNER"))
	cyan.Println(CenterText(headerWidth, "Find the fastest clean IPs"))
	cyan.Println(strings.Repeat("=", headerWidth))
}

func PrintDesigner() {
	green := color.New(color.FgGreen, color.Bold)
	green.Println("...:..::.:::Engineered by: Anonymous:::.::..:...")
	fmt.Println()
}