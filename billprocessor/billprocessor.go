package billprocessor

import (
	"github.com/luketower/roomies/color"
	"github.com/luketower/roomies/linebreak"
	"log"
	"sort"
	"strconv"
	"strings"
)

var MONTHS = map[string]string{
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

func ErrorMsg(args []string) string {
	yellowLineBreak := linebreak.Make("*", 70, "yellow")
	return "\n" +
		yellowLineBreak + "\n" + yellowLineBreak + "\n\n" +
		color.Text("There was a problem with your inputs:\n\n", "red") +
		"  '" + strings.Join(args, " ") + "'\n\n" +
		color.Text("Input should resemble the following:\n\n", "red") +
		"  'date 12/2015 Gas 34.55 Electric 45.99 Rent 933 -- bob 45 susan 55'\n\n" +
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
	name := f.formattedName()
	if len(name) < length {
		s = name + ":" + strings.Repeat(" ", length-len(name))
	} else {
		s = name + ":"
	}
	floatAmount, err := strconv.ParseFloat(f.amount, 64)
	if err != nil {
		log.Fatal("OUCH! ", err)
	}
	return s + " $" + strconv.FormatFloat(floatAmount, 'f', 2, 64) + "\n"
}

func (f *Field) formattedName() (s string) {
	s = strings.Title(f.name)
	if f.isShare {
		arr := strings.Split(s, " ")
		first, rest := arr[0]+"'s", arr[1:]
		s = strings.Join(append([]string{first}, rest...), " ")
	}
	return strings.Replace(s, "-", " ", -1)
}

func (f *Field) nameLength() (l int) {
	if f.isShare {
		l = len(f.name) + 2
	} else {
		l = len(f.name)
	}
	return
}

type Fields []Field

func (fields Fields) longestTitle() (l int) {
	for _, b := range fields {
		if length := b.nameLength(); length > l {
			l = length
		}
	}
	return
}

func (fields Fields) toArr(l int) (arr []string, longest int) {
	sort.Sort(fields)
	for _, f := range fields {
		str := f.toString(l)
		arr = append(arr, str)
		if length := len(str); length > longest {
			longest = length
		}
	}
	return arr, longest
}

func (fields Fields) total() string {
	var total float64
	for _, f := range fields {
		val, _ := strconv.ParseFloat(f.amount, 64)
		total += val
	}
	return strconv.FormatFloat(total, 'f', 2, 64)
}

func (slice Fields) Len() int {
	return len(slice)
}

func (slice Fields) Less(i, j int) bool {
	return slice[i].name < slice[j].name
}

func (slice Fields) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func BillReport(args []string) string {
	bills, shares, header, total := parse(args)
	longestTitle := longestTitleIn(bills, shares)
	billsArr, longestBill := bills.toArr(longestTitle)
	sharesArr, longestShare := shares.toArr(longestTitle)
	length := lineBreakLength([]int{len(header), longestBill, longestShare})
	dottedLine := linebreak.Make("-", length, "green") + "\n"
	return color.Text(header, "blue") + "\n" +
		linebreak.Make("*", length, "green") + "\n" +
		strings.Join(billsArr, "") +
		dottedLine +
		total.toString(longestTitle) +
		dottedLine +
		strings.Join(sharesArr, "")
}

func parse(args []string) (bills Fields, shares Fields, header string, total Field) {
	isShare := false
	for i, arg := range args {
		if isHeader(arg, i, args) {
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
				shares = append(shares, Field{arg + " Total", calcShare(args[i+1], total.amount), isShare})
			} else {
				bills = append(bills, Field{args[i-1], arg, isShare})
			}
		}
	}
	return
}

func longestTitleIn(bills Fields, shares Fields) (l int) {
	if shares.longestTitle() > bills.longestTitle() {
		l = shares.longestTitle()
	} else {
		l = bills.longestTitle()
	}
	return
}

func lineBreakLength(allLengths []int) (l int) {
	l = 25
	for _, num := range allLengths {
		if num > l {
			l = num
		}
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

	return header
}

func isHeader(s string, i int, args []string) (ans bool) {
	ans = isHeaderName(s)
	if i > 0 {
		ans = isHeaderName(args[i-1])
	}
	return ans
}

func isHeaderName(s string) bool {
	return isDate(s) || s == "header"
}

func isDate(s string) bool {
	return s == "date" || s == "month"
}

func calcShare(percent string, total string) string {
	percentFloat, err := strconv.ParseFloat(percent, 64)
	if err != nil {
		log.Fatal("OUCH! Can't parse percent! ", err)
	}
	totalFloat, err := strconv.ParseFloat(total, 64)
	if err != nil {
		log.Fatal("OUCH! Can't parse total! ", err)
	}
	share := totalFloat * (percentFloat / 100.00)
	return strconv.FormatFloat(share, 'f', 2, 64)
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
