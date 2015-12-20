package billprocessor

import (
	"github.com/luketower/roomies/color"
	f "github.com/luketower/roomies/field"
	line "github.com/luketower/roomies/linebreak"
	"log"
	"sort"
	"strconv"
	"strings"
)

var (
	DEFAULT_LINE_BREAK_LENGTH = 25
	MAX_AMOUNT_LENGTH         = 15
	MONTHS                    = map[string]string{
		"01": "January",
		"02": "February",
		"03": "March",
		"04": "April",
		"05": "May",
		"06": "June",
		"07": "July",
		"08": "August",
		"09": "September",
		"10": "October",
		"11": "November",
		"12": "December",
	}
	EXAMPLE_USAGE = "  'date 12/2015 gas 34.55 electric 45.99 rent 933 -- bob 45 susan 55'\n\n" +
		"* You must include the date!\n" +
		"* You must add '--' followed by name/percentage pairs.\n" +
		"  Ex. '(args) -- bob 45 susan 55'\n\n"
)

func ErrorMsg(args []string) string {
	yellowLine := line.Make("*", 70, "yellow")
	yellowLines := yellowLine + "\n" + yellowLine + "\n\n"
	return "\n" +
		yellowLines +
		color.Text("There was a problem with your inputs:\n\n", "red") +
		"  '" + strings.Join(args, " ") + "'\n\n" +
		color.Text("Input should resemble the following:\n\n", "red") +
		EXAMPLE_USAGE +
		yellowLines
}

func BillReport(args []string) string {
	data := parse(args)
	dottedLine := line.Make("-", data.getLineBreakLength(), "green") + "\n"
	return color.Text(data.header, "blue") + "\n" +
		line.Make("*", data.getLineBreakLength(), "green") + "\n" +
		data.bills.ToString(data.longestName()) +
		dottedLine +
		data.total.ToString(data.longestName()) +
		dottedLine +
		data.shares.ToString(data.longestName())
}

type argsData struct {
	bills             f.Fields
	shares            f.Fields
	header            string
	total             f.Field
	longestNameLength int
	lineBreakLength   int
}

func (d *argsData) longestName() (l int) {
	if d.longestNameLength > 0 {
		l = d.longestNameLength
	} else {
		l = append(d.bills, d.shares...).LongestName()
	}
	return
}

func (d *argsData) getLineBreakLength() (l int) {
	if d.lineBreakLength > 0 {
		l = d.lineBreakLength
	} else {
		lengths := []int{len(d.header),
			d.longestName() + MAX_AMOUNT_LENGTH,
			DEFAULT_LINE_BREAK_LENGTH}
		sort.Ints(lengths)
		l = lengths[len(lengths)-1]
	}
	return
}

func parse(args []string) (d argsData) {
	isShare := false
	for i, arg := range args {
		if isPartOfHeader(arg, i, args) {
			d.header = makeHeader(arg, args, i)
			continue
		}
		if arg == "--" {
			isShare = true
			d.total = f.Field{"Total", d.bills.Total(), false}
			continue
		}
		if i%2 != 0 {
			if isShare {
				d.shares = append(d.shares, f.Field{arg + " Total",
					calcShare(args[i+1], d.total.Amount), isShare})
			} else {
				d.bills = append(d.bills, f.Field{args[i-1],
					parseFloat(arg), isShare})
			}
		}
	}
	return
}

func makeHeader(word string, args []string, i int) (header string) {
	if strings.Contains(word, "/") && isDate(args[i-1]) {
		monthAndYear := strings.Split(word, "/")
		header = MONTHS[monthAndYear[0]] + " " + monthAndYear[1]
	} else {
		header = strings.Title(strings.Replace(word, "-", " ", -1))
	}
	return
}

func isPartOfHeader(s string, i int, args []string) (ans bool) {
	if i > 0 {
		ans = hasHeader(args[i-1])
	} else {
		ans = hasHeader(s)
	}
	return
}

func isDate(s string) bool {
	return s == "date" || s == "month"
}

func calcShare(percent string, total float64) float64 {
	return total * (parseFloat(percent) / 100.00)
}

func parseFloat(num string) (f float64) {
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		log.Fatal("OUCH! ", err)
	}
	return
}

func HasValid(args []string) bool {
	stringArgs := strings.Join(args, " ")
	return hasHeader(stringArgs) &&
		strings.Contains(stringArgs, "--") &&
		len(args) >= 7 &&
		len(args)%2 != 0 &&
		hasBills(stringArgs)
}

func hasHeader(args string) bool {
	return strings.Contains(args, "date") ||
		strings.Contains(args, "month") ||
		strings.Contains(args, "header")
}

func hasBills(args string) bool {
	return len(strings.Split(strings.Split(args, " -- ")[0], " ")) > 2
}
