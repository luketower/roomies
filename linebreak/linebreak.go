package linebreak

import (
	"github.com/luketower/roomies/color"
	"strings"
)

func Make(char string, num int, textColor string) (text string) {
	return color.Text(strings.Repeat(char, num), textColor)
}
