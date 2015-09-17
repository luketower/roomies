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
	line := linebreak.Make("*", 70, "yellow")
	return line + "\n" + line + "\n\n"
}

type stringify func(name string) string

func titleize(name string) string       { return strings.Title(name) }
func titleizeAndOwn(name string) string { return owner(strings.Title(name)) + " Total" }

func BillReport(args []string) string {
	bills := map[string]string{}
	billsToMap(args, bills)
	total := total(bills)
	shares := sharesToMap(args, total)
	longestTitle := longestTitleIn(bills, shares)
	billsArr := eachToArr(bills, longestTitle, titleize)
	sharesArr := eachToArr(shares, longestTitle, titleizeAndOwn)
	l := lineBreakLength([]int{len(header(args)), longestIn(sharesArr), longestIn(billsArr)})
	dottedLine := linebreak.Make("-", l, "green") + "\n"
	return color.Text(header(args), "blue") + "\n" +
		linebreak.Make("*", l, "green") + "\n" +
		strings.Join(billsArr, "") +
		dottedLine +
		billToString("Total", total, longestTitle) +
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
			header = titleize(strings.Replace(args[i+1], "-", " ", -1))
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

func isHeaderName(s string) bool {
	return s == "date" || s == "month" || s == "header"
}

func eachToArr(m map[string]string, l int, fn stringify) (arr []string) {
	for _, name := range sortedKeys(m) {
		arr = append(arr, billToString(fn(name), m[name], l))
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
	return strings.Replace(adjustedName, "-", " ", -1) +
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

func owner(shares string) string {
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

func sortedKeys(m map[string]string) (names []string) {
	for k := range m {
		names = append(names, k)
	}
	sort.Strings(names)
	return
}
