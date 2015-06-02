package billsreader

import (
	"bufio"
	"fmt"
	processor "github.com/luketower/roomies/billprocessor"
	a "github.com/mgutz/ansi"
	"os"
	"strings"
)

func errMsg(args []string) string {
	return "\n" + a.Color(lineBreak("*", 50), "yellow") +
		a.Color(lineBreak("*", 50), "yellow") + "\n" +
		a.Color("There was an error processing your input:\n", "red") +
		"  '" + strings.Join(args, " ") + "'\n\n" +
		a.Color("  Ex. 'readbills <filename>'\n\n", "red") +
		a.Color(lineBreak("*", 50), "yellow") +
		a.Color(lineBreak("*", 50), "yellow")
}

func Read(args []string) {
	if len(args) != 1 {
		fmt.Printf("%s", errMsg(args))
		os.Exit(1)
	}
	filename := args[0]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("OUCH! %s\n", err)
		os.Exit(1)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	readFile(scanner)
}

func lineBreak(char string, num int) string {
	return strings.Repeat(char, num) + "\n"
}

func readFile(scanner *bufio.Scanner) {
	for scanner.Scan() {
		args := strings.Split(scanner.Text(), " ")
		if isCommentOrBlankLine(args) {
			continue
		}
		fmt.Println(processor.BillReport(args))
	}
}

func isCommentOrBlankLine(args []string) bool {
	return args[0] == "//" || args[0] == ""
}
