package billprocessor

import (
	"github.com/luketower/roomies/color"
	f "github.com/luketower/roomies/field"
	"github.com/luketower/roomies/linebreak"
	"log"
	"sort"
	"strconv"
	"strings"
)

var (
	DEFAULT_LINE_BREAK_LENGTH = 25
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
)

func ErrorMsg(args []string) string {
	yellowLineBreak := linebreak.Make("*", 70, "yellow")
	return "\n" +
		yellowLineBreak + "\n" + yellowLineBreak + "\n\n" +
		color.Text("There was a problem with your inputs:\n\n", "red") +
		"  '" + strings.Join(args, " ") + "'\n\n" +
		color.Text("Input should resemble the following:\n\n", "red") +
		"  'date 12/2015 gas 34.55 electric 45.99 rent 933 -- bob 45 susan 55'\n\n" +
		"* You must include the date!\n" +
		"* You must add '--' followed by name/percentage pairs.\n" +
		"  Ex. '(args) -- bob 45 susan 55'\n\n" +
		yellowLineBreak + "\n" + yellowLineBreak + "\n\n"
}

func BillReport(args []string) string {
	bills, shares, longestTitle, header, total := parse(args)
	billsStr, longestBill := bills.ToString(longestTitle)
	sharesStr, longestShare := shares.ToString(longestTitle)
	length := lineBreakLength([]int{len(header), longestBill, longestShare})
	dottedLine := linebreak.Make("-", length, "green") + "\n"
	return color.Text(header, "blue") + "\n" +
		linebreak.Make("*", length, "green") + "\n" +
		billsStr +
		dottedLine +
		total.ToString(longestTitle) +
		dottedLine +
		sharesStr
}

func parse(args []string) (bills, shares f.Fields, longestTitle int, header string, total f.Field) {
	isShare := false
	for i, arg := range args {
		if isPartOfHeader(arg, i, args) {
			header = makeHeader(arg, args, i)
			continue
		}
		if arg == "--" {
			isShare, total = true, f.Field{"Total", bills.Total(), false}
			continue
		}
		if i%2 != 0 {
			if isShare {
				shares = append(shares, f.Field{arg + " Total",
					calcShare(args[i+1], total.Amount), isShare})
			} else {
				bills = append(bills, f.Field{args[i-1],
					calcShare("100", parseFloat(arg)), isShare})
			}
		}
	}
	longestTitle = append(bills, shares...).LongestTitle()
	return
}

func lineBreakLength(lengths []int) (l int) {
	sort.Ints(lengths)
	l, highest := DEFAULT_LINE_BREAK_LENGTH, lengths[len(lengths)-1]
	if highest > l {
		l = highest
	}
	return
}

func makeHeader(word string, args []string, i int) (header string) {
	if strings.Contains(word, "/") && isDate(args[i-1]) {
		dateArr := strings.Split(word, "/")
		header = MONTHS[dateArr[0]] + " " + dateArr[1]
	} else {
		header = strings.Title(strings.Replace(word, "-", " ", -1))
	}

	return
}

func isPartOfHeader(s string, i int, args []string) (ans bool) {
	ans = hasHeader(s)
	if i > 0 {
		ans = hasHeader(args[i-1])
	}
	return ans
}

func isDate(s string) bool {
	return s == "date" || s == "month"
}

func calcShare(percent string, total float64) float64 {
	if percent != "100" {
		total = total * (parseFloat(percent) / 100.00)
	}
	return total
}

func parseFloat(num string) (float float64) {
	float, err := strconv.ParseFloat(num, 64)
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
