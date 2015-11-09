package billprocessor

import (
	"github.com/luketower/roomies/color"
	f "github.com/luketower/roomies/field"
	"github.com/luketower/roomies/linebreak"
	"log"
	"math"
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
	yellowLineBreak := linebreak.Make("*", 70, "yellow")
	return "\n" +
		yellowLineBreak + "\n" + yellowLineBreak + "\n\n" +
		color.Text("There was a problem with your inputs:\n\n", "red") +
		"  '" + strings.Join(args, " ") + "'\n\n" +
		color.Text("Input should resemble the following:\n\n", "red") +
		EXAMPLE_USAGE +
		yellowLineBreak + "\n" + yellowLineBreak + "\n\n"
}

func BillReport(args []string) string {
	bills, shares, header, total := parse(args)
	longestName := append(bills, shares...).LongestName()
	l := lineBreakLength([]int{len(header), longestName + MAX_AMOUNT_LENGTH})
	dottedLine := linebreak.Make("-", l, "green") + "\n"
	return color.Text(header, "blue") + "\n" +
		linebreak.Make("*", l, "green") + "\n" +
		bills.ToString(longestName) +
		dottedLine +
		total.ToString(longestName) +
		dottedLine +
		shares.ToString(longestName)
}

func parse(args []string) (bills, shares f.Fields, header string, total f.Field) {
	isShare := false
	for i, arg := range args {
		if isPartOfHeader(arg, i, args) {
			header = makeHeader(arg, args, i)
			continue
		}
		if arg == "--" {
			isShare = true
			total = f.Field{"Total", bills.Total(), false}
			continue
		}
		if i%2 != 0 {
			if isShare {
				shares = append(shares, f.Field{arg + " Total",
					calcShare(args[i+1], total.Amount), isShare})
			} else {
				bills = append(bills, f.Field{args[i-1],
					parseInt(arg), isShare})
			}
		}
	}
	return
}

func lineBreakLength(lengths []int) int {
	lengths = append(lengths, DEFAULT_LINE_BREAK_LENGTH)
	sort.Ints(lengths)
	return lengths[len(lengths)-1]
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

func calcShare(percent string, total int) int {
	floatTotal := float64(total) / 100.00
	floatPercent, err := strconv.ParseFloat(percent, 64)
	if err != nil {
		log.Fatal("OUCH! ", err)
	}
	floatShare := round((floatTotal * (floatPercent / 100.00)), .5, 2)
	return int(floatShare * 100)
}

// round func taken from http://play.golang.org/p/yjfShH_uEy
func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func parseInt(num string) (i int) {
	if strings.Contains(num, ".") {
		num = strings.Join(strings.Split(num, "."), "")
	} else {
		num += "00"
	}
	i, err := strconv.Atoi(num)
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
