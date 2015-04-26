package main

import (
	"fmt"
	a "github.com/mgutz/ansi"
	"log"
	"os"
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

func main() {
	args := os.Args[1:]
	if valid(args) {
		fmt.Println(billReport(args))
	} else {
		fmt.Println(errorMsg(args))
	}
}

func errorMsg(args []string) string {
	return "\n" +
		a.Color(lineBreak("*", 70), "yellow") + "\n" +
		a.Color(lineBreak("*", 70), "yellow") + "\n\n" +
		a.Color("There was a problem with your inputs:\n", "red") + "\n" +
		"  '" + strings.Join(args, " ") + "'\n\n" +
		a.Color("Input should resemble the following:\n", "red") + "\n" +
		"  'date 12/2015 Gas 34.55 Electric 45.99 Rent 933 -- bob 45 susan 55'\n\n" +
		"* You must include the date!\n" +
		"* You must add '--' followed by name/percentage pairs.\n" +
		"  Ex. '(args) -- bob 45 susan 55'\n\n" +
		a.Color(lineBreak("*", 70), "yellow") + "\n" +
		a.Color(lineBreak("*", 70), "yellow") + "\n"
}

func billReport(args []string) string {
	bills := map[string]string{}
	billsMap(args, bills)
	return a.Color(monthHeader(args), "blue") + "\n" +
		a.Color(lineBreak("*", 25), "green") + "\n" +
		eachBill(bills) +
		a.Color(lineBreak("-", 25), "green") + "\n" +
		billToString("Total", calcTotal(bills)) +
		a.Color(lineBreak("-", 25), "green") + "\n" +
		individualShares(calcTotal(bills), args)
}

func lineBreak(char string, num int) string {
	return strings.Repeat(char, num)
}

func monthHeader(args []string) string {
	dateArr := strings.Split(args[1], "/")
	return MONTHS[dateArr[0]] + " " + dateArr[1]
}

func billsMap(args []string, billsMap map[string]string) {
	for i, s := range args {
		if s == "--" {
			break
		}
		if isMonth(s) {
			continue
		}
		if i%2 != 0 {
			key := args[i-1]
			billsMap[key] = s
		}
	}
}

func isMonth(s string) bool {
	return s == "date" ||
		s == "month" ||
		strings.Contains(s, "/") && s == os.Args[2]
}

func eachBill(bills map[string]string) (allBills string) {
	names := orderKeys(bills)
	for _, n := range names {
		allBills += billToString(strings.Title(n), bills[n])
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

func calcTotal(bills map[string]string) string {
	var total float64
	for _, v := range bills {
		val, _ := strconv.ParseFloat(v, 64)
		total += val
	}
	return strconv.FormatFloat(total, 'f', 2, 64)
}

func individualShares(total string, args []string) (totalShares string) {
	shares := map[string]string{}
	calcShares(args, total, shares)
	names := orderKeys(shares)
	for _, name := range names {
		billName := strings.Title(name) + "'s Total"
		totalShares += billToString(billName, shares[name])
	}
	return
}

func calcShares(args []string, total string, shares map[string]string) {
	stringArgs := strings.Join(args, " ")
	stringShares := strings.Split(stringArgs, " -- ")[1]
	sharesArr := strings.Split(stringShares, " ")
	for i, s := range sharesArr {
		if i%2 != 0 {
			key := sharesArr[i-1]
			shares[key] = calcShare(s, total)
		}
	}
}

func calcShare(percent string, total string) string {
	percentFloat, err := strconv.ParseFloat(percent, 64)
	totalFloat, err := strconv.ParseFloat(total, 64)
	if err != nil {
		log.Fatal("OUCH! ", err)
	}
	share := totalFloat * (percentFloat / 100.00)
	return strconv.FormatFloat(share, 'f', 2, 64)
}

func valid(args []string) bool {
	hasDate := includeIn(args, "date") || includeIn(args, "month")
	hasDash := includeIn(args, "--")
	hasMinimumCount := len(args) >= 9
	isOdd := len(args)%2 != 0
	return hasDate && hasDash && hasMinimumCount && isOdd
}

func includeIn(args []string, flag string) bool {
	stringArgs := strings.Join(args, " ")
	return strings.Contains(stringArgs, flag)
}

func orderKeys(m map[string]string) (names []string) {
	for k := range m {
		names = append(names, k)
	}
	sort.Strings(names)
	return
}
