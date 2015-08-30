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
	MONTHS = map[string]string{
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
	return "\n" +
		linebreak.Make("*", 70, "yellow") + "\n" +
		linebreak.Make("*", 70, "yellow") + "\n\n" +
		color.Text("There was a problem with your inputs:\n", "red") + "\n" +
		"  '" + strings.Join(args, " ") + "'\n\n" +
		color.Text("Input should resemble the following:\n", "red") + "\n" +
		"  'date 12/2015 Gas 34.55 Electric 45.99 Rent 933 -- bob 45 susan 55'\n\n" +
		"* You must include the date!\n" +
		"* You must add '--' followed by name/percentage pairs.\n" +
		"  Ex. '(args) -- bob 45 susan 55'\n\n" +
		linebreak.Make("*", 70, "yellow") + "\n" +
		linebreak.Make("*", 70, "yellow") + "\n"
}

func BillReport(args []string) string {
	bills := map[string]string{}
	billsMap(args, bills)
	total := total(bills)
	dottedLine := linebreak.Make("-", 25, "green") + "\n"
	return color.Text(header(args), "blue") + "\n" +
		linebreak.Make("*", 25, "green") + "\n" +
		eachBill(bills) +
		dottedLine +
		billToString("Total", total) +
		dottedLine +
		individualShares(total, args)
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

func billsMap(args []string, bills map[string]string) {
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

func eachBill(bills map[string]string) (allBills string) {
	for _, name := range sortedKeys(bills) {
		allBills += billToString(strings.Title(name), bills[name])
	}
	return
}

func billToString(name, amount string) string {
	floatAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Fatal("OUCH! ", err)
	}
	return name + ":" + tabsFor(name) +
		strconv.FormatFloat(floatAmount, 'f', 2, 64) + "\n"
}

func tabsFor(text string) (tabs string) {
	tabs = "\t$"
	if len(text) <= 5 {
		tabs = "\t" + tabs
	}
	return
}

func total(bills map[string]string) string {
	var total float64
	for _, v := range bills {
		val, _ := strconv.ParseFloat(v, 64)
		total += val
	}
	return strconv.FormatFloat(total, 'f', 2, 64)
}

func individualShares(total string, args []string) (totalShares string) {
	shares := shares(args, total)
	for _, name := range sortedKeys(shares) {
		billName := strings.Title(name) + "'s Total"
		totalShares += billToString(billName, shares[name])
	}
	return
}

func shares(args []string, total string) map[string]string {
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
