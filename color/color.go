package color

var (
	COLORS = map[string]string{
		"red":    "31",
		"green":  "32",
		"yellow": "33",
		"blue":   "34",
	}
)

func Text(text string, color string) string {
	return "\033[0" + COLORS[color] + "m" + text + "\033[0m"
}
