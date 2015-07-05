package billsreader

import (
	"bufio"
	"fmt"
	processor "github.com/luketower/roomies/billprocessor"
	"github.com/luketower/roomies/color"
	"github.com/luketower/roomies/linebreak"
	"os"
	"strings"
)

func errMsg(args []string) string {
	yellowLines := linebreak.Make("*", 50, "yellow") +
		linebreak.Make("*", 50, "yellow")
	return "\n" +
		yellowLines +
		color.Text("There was an error processing your input:\n", "red") +
		"  '" + strings.Join(args, " ") + "'\n\n" +
		color.Text("  Ex. 'readbills <filename>'\n\n", "red") +
		yellowLines
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
