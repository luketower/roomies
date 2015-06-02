package main

import (
	"fmt"
	processor "github.com/luketower/roomies/billprocessor"
	reader "github.com/luketower/roomies/billsreader"
	"os"
)

func main() {
	args := os.Args[1:]
	if processingOneMonth(args) {
		if processor.HasValid(args) {
			fmt.Println(processor.BillReport(args))
		} else {
			fmt.Println(processor.ErrorMsg(args))
		}
	} else {
		reader.Read(args)
	}
}

func processingOneMonth(args []string) bool {
	return len(args) > 1
}
