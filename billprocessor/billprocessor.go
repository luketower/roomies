package billprocessor

import (
	"github.com/luketower/roomies/color"
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

type Field struct {
	name    string
	amount  string
	isShare bool
}

func (f *Field) toString(length int) (s string) {
	name := f.formatName()
	if nameLength := len(name); nameLength < length {
		s = name + ":" + strings.Repeat(" ", length-nameLength)
	} else {
		s = name + ":"
	}
	return s + " $" + f.amount + "\n"
}

func (f *Field) formatName() (s string) {
	s = strings.Title(f.name)
	if f.isShare {
		arr := strings.Split(s, " ")
		first, rest := arr[0]+"'s", arr[1:]
		s = strings.Join(append([]string{first}, rest...), " ")
	}
	return strings.Replace(s, "-", " ", -1)
}

type Fields []Field

func (fields Fields) longestTitle() (l int) {
	for _, b := range fields {
		if length := len(b.formatName()); length > l {
			l = length
		}
	}
	return
}

func (fields Fields) toString(l int) (s string, longest int) {
	sort.Sort(fields)
	for _, f := range fields {
		str := f.toString(l)
		s += str
		if length := len(str); length > longest {
			longest = length
		}
	}
	return s, longest
}

func (fields Fields) total() string {
	var total float64
	for _, f := range fields {
		val := parseFloat(f.amount)
		total += val
	}
	return strconv.FormatFloat(total, 'f', 2, 64)
}

func (fields Fields) Len() int {
	return len(fields)
}

func (fields Fields) Less(i, j int) bool {
	return fields[i].name < fields[j].name
}

func (fields Fields) Swap(i, j int) {
	fields[i], fields[j] = fields[j], fields[i]
}

func BillReport(args []string) string {
	bills, shares, longestTitle, header, total := parse(args)
	billsStr, longestBill := bills.toString(longestTitle)
	sharesStr, longestShare := shares.toString(longestTitle)
	length := lineBreakLength([]int{len(header), longestBill, longestShare})
	dottedLine := linebreak.Make("-", length, "green") + "\n"
	return color.Text(header, "blue") + "\n" +
		linebreak.Make("*", length, "green") + "\n" +
		billsStr +
		dottedLine +
		total.toString(longestTitle) +
		dottedLine +
		sharesStr
}

func parse(args []string) (bills, shares Fields, longestTitle int, header string, total Field) {
	isShare := false
	for i, arg := range args {
		if isPartOfHeader(arg, i, args) {
			header = makeHeader(arg, args, i)
			continue
		}
		if arg == "--" {
			isShare = true
			total = Field{"Total", bills.total(), false}
			continue
		}
		if i%2 != 0 {
			if isShare {
				shares = append(shares, Field{arg + " Total",
					calcShare(args[i+1], total.amount), isShare})
			} else {
				bills = append(bills, Field{args[i-1],
					calcShare("100", arg), isShare})
			}
		}
	}
	longestTitle = append(bills, shares...).longestTitle()
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

func calcShare(percent string, total string) (calc string) {
	s := parseFloat(total)
	if percent != "100" {
		s = s * (parseFloat(percent) / 100.00)
	}
	return strconv.FormatFloat(s, 'f', 2, 64)
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
