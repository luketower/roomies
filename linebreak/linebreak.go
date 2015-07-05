package linebreak

import (
	"github.com/luketower/roomies/color"
	"strings"
)

func Make(char string, num int, textColor string) (text string) {
	if textColor == "none" {
		text = strings.Repeat(char, num)
	} else {
		text = color.Text(strings.Repeat(char, num), textColor)
	}
	return
}
