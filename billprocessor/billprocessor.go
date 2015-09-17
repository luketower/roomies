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
	headerAndFooter := yellowLines()
	return "\n" +
		headerAndFooter +
		color.Text("There was a problem with your inputs:\n\n", "red") +
		"  '" + strings.Join(args, " ") + "'\n\n" +
		color.Text("Input should resemble the following:\n\n", "red") +
		"  'date 12/2015 Gas 34.55 Electric 45.99 Rent 933 -- bob 45 susan 55'\n\n" +
		"* You must include the date!\n" +
		"* You must add '--' followed by name/percentage pairs.\n" +
		"  Ex. '(args) -- bob 45 susan 55'\n\n" +
		headerAndFooter
}

func yellowLines() (lines string) {
	return linebreak.Make("*", 70, "yellow") + "\n" +
		linebreak.Make("*", 70, "yellow") + "\n\n"
}

type stringify func(name, amount string, length int) string

func billStringifier(name string, amount string, length int) string {
	return billToString(strings.Title(name), amount, length)
}
func shareStringifier(name string, amount string, length int) string {
	return possessive(billToString(name+" Total", amount, length))
}

func BillReport(args []string) string {
	bills := map[string]string{}
	billsToMap(args, bills)
	total := total(bills)
	shares := sharesToMap(args, total)
	longestBillTitle := longestTitleIn(bills, shares)
	billsArr := eachToArr(bills, longestBillTitle, billStringifier)
	sharesArr := eachToArr(shares, longestBillTitle-2, shareStringifier)
	l := lineBreakLength([]int{len(header(args)), longestIn(sharesArr), longestIn(billsArr)})
	dottedLine := linebreak.Make("-", l, "green") + "\n"
	return color.Text(header(args), "blue") + "\n" +
		linebreak.Make("*", l, "green") + "\n" +
		strings.Join(billsArr, "") +
		dottedLine +
		billToString("Total", total, longestBillTitle) +
		dottedLine +
		strings.Join(sharesArr, "")
}

func longestTitleIn(bills map[string]string, shares map[string]string) (l int) {
	for k := range bills {
		if length := len(k); length > l {
			l = length
		}
	}
	for k := range shares {
		if length := len(k) + 8; length > l {
			l = length
		}
	}
	return
}

func longestIn(arr []string) (l int) {
	for _, str := range arr {
		if strLength := len(str); strLength > l {
			l = strLength
		}
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

func header(args []string) (header string) {
	for i, w := range args {
		if w == "month" || w == "date" {
			dateArr := strings.Split(args[i+1], "/")
			header = MONTHS[dateArr[0]] + " " + dateArr[1]
		}
		if w == "header" {
			header = strings.Title(strings.Replace(args[i+1], "-", " ", -1))
		}
	}
	return
}

func billsToMap(args []string, bills map[string]string) {
	for i, s := range args {
		if s == "--" {
			break
		}
		if isHeader(s, i, args) {
			continue
		}
		if i%2 != 0 {
			key := args[i-1]
			bills[key] = s
		}
	}
}

func isHeader(s string, i int, args []string) (ans bool) {
	if isHeaderName(s) {
		ans = true
	}
	if i > 0 {
		ans = isHeaderName(args[i-1])
	}
	return ans
}

func hasHeader(args []string) bool {
	return includeIn(args, "date") ||
		includeIn(args, "month") ||
		includeIn(args, "header")
}

func isHeaderName(s string) bool {
	return s == "date" || s == "month" || s == "header"
}

func eachToArr(m map[string]string, l int, fn stringify) (arr []string) {
	for _, name := range sortedKeys(m) {
		arr = append(arr, fn(name, m[name], l))
	}
	return
}

func billToString(name, amount string, i int) string {
	var adjustedName string
	if len(name) < i {
		adjustedName = name + ":" + strings.Repeat(" ", i-len(name))
	} else {
		adjustedName = name + ":"
	}
	floatAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Fatal("OUCH! ", err)
	}
	return strings.Title(strings.Replace(adjustedName, "-", " ", -1)) +
		" $" +
		strconv.FormatFloat(floatAmount, 'f', 2, 64) + "\n"
}

func total(bills map[string]string) string {
	var total float64
	for _, v := range bills {
		val, _ := strconv.ParseFloat(v, 64)
		total += val
	}
	return strconv.FormatFloat(total, 'f', 2, 64)
}

func possessive(shares string) string {
	arr := strings.Split(shares, " ")
	first, rest := arr[0]+"'s", arr[1:]
	return strings.Join(append([]string{first}, rest...), " ")
}

func sharesToMap(args []string, total string) map[string]string {
	var (
		shares       = map[string]string{}
		stringArgs   = strings.Join(args, " ")
		stringShares = strings.Split(stringArgs, " -- ")[1]
		sharesArr    = strings.Split(stringShares, " ")
	)
	for i, s := range sharesArr {
		if i%2 != 0 {
			shares[sharesArr[i-1]] = calcShare(s, total)
		}
	}
	return shares
}

func calcShare(percent string, total string) string {
	percentFloat, err := strconv.ParseFloat(percent, 64)
	if err != nil {
		log.Fatal("OUCH! ", err)
	}
	totalFloat, err := strconv.ParseFloat(total, 64)
	if err != nil {
		log.Fatal("OUCH! ", err)
	}
	share := totalFloat * (percentFloat / 100.00)
	return strconv.FormatFloat(share, 'f', 2, 64)
}

func HasValid(args []string) bool {
	hasDash := includeIn(args, "--")
	hasMinimumCount := len(args) >= 7
	isOdd := len(args)%2 != 0
	return hasHeader(args) && hasDash && hasMinimumCount && isOdd
}

func includeIn(args []string, flag string) bool {
	stringArgs := strings.Join(args, " ")
	return strings.Contains(stringArgs, flag)
}

func sortedKeys(m map[string]string) (names []string) {
	for k := range m {
		names = append(names, k)
	}
	sort.Strings(names)
	return
}
